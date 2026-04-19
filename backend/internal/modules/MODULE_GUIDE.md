# Module Creation Guide

## Folder Structure

```
backend/internal/modules/<module-name>/
├── <module-name>.module.go        # Route registration + DI wiring
├── <module-name>.controller.go    # HTTP handlers
├── dtos/
│   ├── req/
│   │   └── <module-name>.req.dto.go   # Request DTOs
│   └── res/
│       └── <module-name>.res.dto.go  # Response DTOs
├── services/
│   ├── <module-name>.service.go      # Service interface + implementation
│   └── <module-name>.service_test.go # Service tests
└── repositories/
    ├── <module-name>.repository.go      # Repository interface + implementation
    └── <module-name>.repository.mock.go # Repository mock for testing
```

## Naming Convention

| Layer | File Pattern | Example |
|-------|-------------|---------|
| Module | `<module-name>.module.go` | `videos.module.go` |
| Controller | `<module-name>.controller.go` | `videos.controller.go` |
| Service | `<module-name>.service.go` | `video.service.go` |
| Repository | `<module-name>.repository.go` | `video.repository.go` |
| Request DTO | `<module-name>.req.dto.go` | `video.req.dto.go` |
| Response DTO | `<module-name>.res.dto.go` | `video.res.dto.go` |

**Rules:**
- Use **singular** for code files (`video`, `asset`)
- Use **plural** for folder/module name (`videos`, `assets`)
- DTOs are use-case based, not entity based
- Models are defined in `internal/database/models/`

---

## Code Snippets

### 1. Module (`<module-name>.module.go`)

```go
package <module-plural>

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/<module-plural>/repositories"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/<module-plural>/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, dbClient *dynamodb.Client) {
	repo := repositories.New<ModuleName>Repository(dbClient)
	svc := services.New<ModuleName>Service(repo)
	ctrl := New<ModuleName>Controller(svc)

	group := r.Group("/<module-plural>")
	{
		group.POST("", ctrl.Create<ModuleName>)
		group.GET("", ctrl.List<ModuleName>s)
		group.GET("/:id", ctrl.Get<ModuleName>)
	}
}
```

### 2. Controller (`<module-name>.controller.go`)

```go
package <module-plural>

import (
	"net/http"
	"strconv"

	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/<module-plural>/dtos/req"
	_ "github.com/gianghp123/Vidmerce/backend/internal/modules/<module-plural>/dtos/res"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/<module-plural>/services"
	"github.com/gin-gonic/gin"
)

type <ModuleName>Controller struct {
	svc services.<ModuleName>Service
}

func New<ModuleName>Controller(svc services.<ModuleName>Service) *<ModuleName>Controller {
	return &<ModuleName>Controller{svc: svc}
}

// Create<ModuleName> godoc
// @Summary      Create a new <module-name>
// @Description  Create a <module-name>
// @Tags         <module-plural>
// @Accept       json
// @Produce     json
// @Param        body  body      req.Create<ModuleName>Req  true  "Create <module-name> request"
// @Success      201  {object}  response.BaseResponse[res.Create<ModuleName>Res]
// @Failure      400  {object}  response.BaseResponse[any]
// @Failure      500  {object}  response.BaseResponse[any]
// @Router      /<module-plural> [post]
func (ctrl *<ModuleName>Controller) Create<ModuleName>(c *gin.Context) {
	var body req.Create<ModuleName>Req
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest()))
		return
	}

	result, appErr := ctrl.svc.Create<ModuleName>(c, body)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusCreated, response.Success(result))
}

// List<ModuleName>s godoc
// @Summary      List <module-plural>
// @Description  Get a paginated list of <module-plural>
// @Tags         <module-plural>
// @Accept       json
// @Produce     json
// @Param        limit   query     int     false  "Number of items per page"    Format(int32)
// @Param        last_key query    string  false  "Pagination cursor"
// @Success      200     {object}  response.BaseResponse[[]res.<ModuleName>Response]
// @Failure      400     {object}  response.BaseResponse[any]
// @Failure      500     {object}  response.BaseResponse[any]
// @Router      /<module-plural> [get]
func (ctrl *<ModuleName>Controller) List<ModuleName>s(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	cursor := c.Query("cursor")

	result, appErr := ctrl.svc.List<ModuleName>s(c, limit, cursor)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMeta(result.Data, result.Meta))
}

// Get<ModuleName> godoc
// @Summary      Get <module-name> by ID
// @Description  Retrieve detailed information about a specific <module-name>
// @Tags         <module-plural>
// @Accept       json
// @Produce     json
// @Param        id   path      string  true  "<ModuleName> ID"
// @Success      200  {object}  response.BaseResponse[res.<ModuleName>DetailResponse]
// @Failure      400  {object}  response.BaseResponse[any]
// @Failure      404  {object}  response.BaseResponse[any]
// @Failure      500  {object}  response.BaseResponse[any]
// @Router      /<module-plural>/{id} [get]
func (ctrl *<ModuleName>Controller) Get<ModuleName>(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.Fail(response.BadRequest()))
		return
	}

	result, appErr := ctrl.svc.Get<ModuleName>(c, id)
	if appErr != nil {
		c.JSON(appErr.Code, response.Fail(appErr))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}
```

### 3. Service (`services/<module-name>.service.go`)

```go
package services

import (
	"context"

	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/<module-plural>/dtos/req"
	res "github.com/gianghp123/Vidmerce/backend/internal/modules/<module-plural>/dtos/res"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/<module-plural>/repositories"
)

type <ModuleName>Service interface {
	Create<ModuleName>(ctx context.Context, req req.Create<ModuleName>Req) (*res.Create<ModuleName>Res, *response.AppError)
	List<ModuleName>s(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.<ModuleName>Response], *response.AppError)
	Get<ModuleName>(ctx context.Context, id string) (*res.<ModuleName>DetailResponse, *response.AppError)
}

type <module-name>Service struct {
	repo repositories.<ModuleName>Repository
}

func New<ModuleName>Service(repo repositories.<ModuleName>Repository) <ModuleName>Service {
	return &<module-name>Service{repo: repo}
}

func (s *<module-name>Service) Create<ModuleName>(ctx context.Context, req req.Create<ModuleName>Req) (*res.Create<ModuleName>Res, *response.AppError) {
	// TODO: Implement
	return nil, response.Internal("not implemented")
}

func (s *<module-name>Service) List<ModuleName>s(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.<ModuleName>Response], *response.AppError) {
	if limit <= 0 {
		limit = 10
	}

	query := req.Get<ModuleName>sQuery{Limit: limit, LastKey: cursor}
	result, err := s.repo.FindAll(ctx, query.Limit, query.LastKey)
	if err != nil {
		return nil, response.Internal("failed to fetch <module-plural>")
	}

	items := make([]res.<ModuleName>Response, 0, len(result.Data))
	for _, item := range result.Data {
		items = append(items, res.<ModuleName>Response{
			ID:        item.PK,
			Title:     item.Title,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
		})
	}

	return &response.PaginatedResult[res.<ModuleName>Response]{
		Data: items,
		Meta: result.Meta,
	}, nil
}

func (s *<module-name>Service) Get<ModuleName>(ctx context.Context, id string) (*res.<ModuleName>DetailResponse, *response.AppError) {
	param := req.Get<ModuleName>ByIDParam{ID: id}
	item, err := s.repo.FindByID(ctx, param.ID)
	if err != nil {
		return nil, response.Internal("failed to fetch <module-name>")
	}
	if item == nil {
		return nil, response.NotFound("<module-name> not found")
	}

	return &res.<ModuleName>DetailResponse{
		<ModuleName>Response: res.<ModuleName>Response{
			ID:        item.PK,
			Title:     item.Title,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
		},
	}, nil
}
```

### 4. Repository (`repositories/<module-name>.repository.go`)

```go
package repositories

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/gianghp123/Vidmerce/backend/internal/core"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
)

type <ModuleName>Repository interface {
	FindAll(ctx context.Context, limit int, lastKey string) (*response.PaginatedResult[models.<ModuleName>Entity], error)
	FindByID(ctx context.Context, id string) (*models.<ModuleName>Entity, error)
}

type <module-name>Repository struct {
	dbClient *dynamodb.Client
}

func New<ModuleName>Repository(dbClient *dynamodb.Client) <ModuleName>Repository {
	return &<module-name>Repository{dbClient: dbClient}
}

func (r *<module-name>Repository) FindAll(ctx context.Context, limit int, lastKey string) (*response.PaginatedResult[models.<ModuleName>Entity], error) {
	exclusiveStartKey, err := core.DecodeCursor(lastKey)
	if err != nil {
		return nil, err
	}

	keyCond := expression.Key("GSI1PK").Equal(expression.Value("ENTITY#<MODULE_NAME>"))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 aws.String(core.TableName),
		IndexName:                 aws.String("GSI1"),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(int32(limit)),
		ScanIndexForward:          aws.Bool(false),
	}
	if exclusiveStartKey != nil {
		input.ExclusiveStartKey = exclusiveStartKey
	}

	resp, err := r.dbClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var items []models.<ModuleName>Entity
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &items); err != nil {
		return nil, err
	}

	for i := range items {
		items[i].PK = strings.TrimPrefix(items[i].PK, "<MODULE_NAME>#")
	}

	nextCursor, err := core.EncodeCursor(resp.LastEvaluatedKey)
	if err != nil {
		return nil, err
	}

	hasMore := resp.LastEvaluatedKey != nil
	return &response.PaginatedResult[models.<ModuleName>Entity]{
		Data: items,
		Meta: response.NewCursorMeta(limit, nextCursor, hasMore),
	}, nil
}

func (r *<module-name>Repository) FindByID(ctx context.Context, id string) (*models.<ModuleName>Entity, error) {
	key, err := attributevalue.MarshalMap(map[string]string{
		"PK": "<MODULE_NAME>#" + id,
		"SK": "METADATA",
	})
	if err != nil {
		return nil, err
	}

	resp, err := r.dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(core.TableName),
		Key:       key,
	})
	if err != nil {
		return nil, err
	}

	if resp.Item == nil {
		return nil, nil
	}

	var item models.<ModuleName>Entity
	if err := attributevalue.UnmarshalMap(resp.Item, &item); err != nil {
		return nil, err
	}
	item.PK = strings.TrimPrefix(item.PK, "<MODULE_NAME>#")
	return &item, nil
}
```

### 5. Request DTO (`dtos/req/<module-name>.req.dto.go`)

```go
package req

type Get<ModuleName>sQuery struct {
	Limit   int    `form:"limit,default=10" example:"10"`
	LastKey string `form:"last_key" example:"eyJ2b3RlSWQiOjF9"`
}

type Get<ModuleName>ByIDParam struct {
	ID string `uri:"id" binding:"required" example:"<module-name>_123456"`
}

type Create<ModuleName>Req struct {
	Title string `json:"title" binding:"required" example:"My <ModuleName>"`
}
```

### 6. Response DTO (`dtos/res/<module-name>.res.dto.go`)

```go
package res

import "github.com/gianghp123/Vidmerce/backend/internal/core/enums"

type <ModuleName>Response struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Status    enums.<ModuleName>Status `json:"status"`
	CreatedAt string            `json:"created_at"`
}

type <ModuleName>DetailResponse struct {
	<ModuleName>Response
}

type Create<ModuleName>Res struct {
	<ModuleName>ID string            `json:"<module-name>Id"`
	Status  enums.<ModuleName>Status `json:"status"`
	Message string            `json:"message"`
}
```

---

## Placeholders Reference

| Placeholder | Description | Example |
|-------------|-------------|---------|
| `<module-name>` | Full module name (lowercase) | `video` |
| `<module-plural>` | Module name (plural/lowercase) | `videos` |
| `<ModuleName>` | Module name (PascalCase) | `Video` |
| `<MODULE_NAME>` | Module name (UPPER_SNAKE) | `VIDEO` |
| `<module-plural>` | Service layer struct (lowercase) | `videoService` |
| `res` | Response DTO package alias | `res` |

---

## Key Conventions

1. **Context**: All service/repository methods accept `context.Context` as first param
2. **DI**: Repositories receive `*dynamodb.Client` via constructor
3. **Errors**: Use `*core.AppError` - `response.BadRequest()`, `response.NotFound()`, `response.Internal()`, etc.
4. **Responses**: Use `core.Success()`, `core.SuccessWithMeta()`, `core.Fail()`
5. **Pagination**: Cursor-based via `core.DecodeCursor()` / `core.EncodeCursor()`
6. **DTO Tags**: Include `example` tags for Swagger documentation