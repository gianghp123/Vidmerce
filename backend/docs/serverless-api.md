# Vidmerce Backend - Serverless API Documentation

## Overview

Vidmerce is a serverless backend API built with Go, using Gin framework for HTTP and AWS DynamoDB for data storage. The project runs as a single service that can run locally or deploy as an AWS Lambda function.

### Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.x |
| HTTP Framework | Gin |
| Database | AWS DynamoDB |
| Object Storage | AWS S3 |
| Deployment | AWS Lambda (via aws-lambda-go-api-proxy) |
| Documentation | Swagger (swaggo) |

### Services

| Service | Port | Base Path | Description |
|---------|------|-----------|-------------|
| Assets | 3000 | `/api` | Product/asset management, campaigns, image uploads |

---

## Project Structure

```
backend/
├── cmd/
│   └── assets/
│       ├── main.go           # Entry point (Gin + Lambda)
│       └── docs/             # Swagger docs
├── internal/
│   ├── configs/              # AWS SDK configuration
│   ├── core/
│   │   ├── constants.go      # Table name, max limits
│   │   ├── cursor.go         # Cursor pagination utilities
│   │   ├── enums/            # Status enums
│   │   └── response/         # BaseResponse, AppError helpers
│   ├── database/
│   │   ├── models/           # DynamoDB entity models
│   │   └── repositories/     # Base repository with transaction support
│   ├── modules/
│   │   ├── assets/           # Assets module (controller, service, repo)
│   │   └── campaigns/        # Campaigns module (controller, service, repo)
│   ├── storage/              # S3 storage abstraction
│   └── utils/                # DTO mappers, DynamoDB utilities, CDN utilities
└── go.mod
```

---

## Database Schema

### DynamoDB Table: `MediaProjectTable`

The project uses a **single-table design** pattern with composite keys and GSI (Global Secondary Index) for flexible querying.

#### Base Item Structure

All entities embed `BaseItem` which provides:

| Attribute | Type | Description |
|-----------|------|-------------|
| `PK` | String | Partition key (format: `ENTITY_TYPE#ID`) |
| `SK` | String | Sort key (format: `METADATA` or sub-type) |
| `GSI1PK` | String | GSI1 partition key for secondary queries |
| `GSI1SK` | String | GSI1 sort key for secondary queries |

---

### Entities

#### 1. AssetEntity

Represents a product/asset with images.

| Attribute | Type | Description |
|-----------|------|-------------|
| `PK` | String | `ASSET#{id}` |
| `SK` | String | `METADATA` |
| `GSI1PK` | String | `ENTITY#ASSET` |
| `GSI1SK` | String | `id` |
| `name` | String | Product name |
| `price` | float64 | Product price |
| `productUrl` | String | URL to product page |
| `status` | AssetStatus | Current status |
| `imageCount` | int | Number of uploaded images |
| `createdAt` | String | ISO timestamp |

#### 2. ImageEntity

Represents an image uploaded for an asset.

| Attribute | Type | Description |
|-----------|------|-------------|
| `PK` | String | `ASSET#{assetId}` |
| `SK` | String | `IMAGE#{imageId}` |
| `fileKey` | String | S3 object key |
| `status` | ImageStatus | Upload status |
| `order` | int | Display order |

#### 3. CampaignEntity

Represents a marketing campaign linked to an asset.

| Attribute | Type | Description |
|-----------|------|-------------|
| `PK` | String | `CAMPAIGN#{id}` |
| `SK` | String | `METADATA` |
| `GSI1PK` | String | `ENTITY#CAMPAIGN` |
| `GSI1SK` | String | `id` |
| `assetId` | String | Linked asset ID |
| `status` | CampaignStatus | Campaign status |
| `storyboard` | Storyboard | Campaign storyboard |
| `createdAt` | String | ISO timestamp |
| `updatedAt` | String | ISO timestamp |

#### 4. Storyboard

Campaign storyboard with slides.

| Attribute | Type | Description |
|-----------|------|-------------|
| `headline` | String | Campaign headline |
| `bodyCopy` | String | Body text |
| `cta` | String | Call-to-action |
| `slides` | []Slide | Array of slides |

#### 5. Slide

A single slide in the storyboard.

| Attribute | Type | Description |
|-----------|------|-------------|
| `slideNumber` | int | Slide order |
| `imageRole` | ImageRole | Role of image (HERO, DETAIL, etc.) |
| `s3Key` | String | S3 key for slide image |
| `overlayText` | String | Text overlay on slide |

#### 6. JobEntity

Represents a background job for processing.

| Attribute | Type | Description |
|-----------|------|-------------|
| `PK` | String | `JOB#{id}` |
| `SK` | String | `METADATA` |
| `GSI1PK` | String | `ENTITY#JOB` |
| `GSI1SK` | String | `id` |
| `targetId` | String | ID of entity being processed |
| `type` | JobType | Job type |
| `status` | JobStatus | Job status |
| `payload` | map[string]any | Job-specific data |
| `errorLog` | String | Error message if failed |
| `createdAt` | String | ISO timestamp |

---

## Status Enums

### Asset Status

| Status | Description |
|--------|-------------|
| `IMPORTING` | Asset is being imported from product URL (scraping in progress) |
| `UPLOADING` | Asset created, images being uploaded |
| `DRAFT` | Asset saved as draft |
| `COMPLETED` | All images uploaded successfully |
| `PARTIAL` | Some images uploaded |
| `FAILED` | Upload failed |

### Image Status

| Status | Description |
|--------|-------------|
| `UPLOADING` | Image being uploaded to S3 |
| `COMPLETED` | Upload successful |
| `FAILED` | Upload failed |

### Campaign Status

| Status | Description |
|--------|-------------|
| `DRAFT` | Campaign created, not yet published |
| `PUBLISHED` | Campaign is live |

### Job Status

| Status | Description |
|--------|-------------|
| `PENDING` | Job queued, not started |
| `PROCESSING` | Job is running |
| `COMPLETED` | Job finished successfully |
| `FAILED` | Job failed |

### Job Type

| Type | Description |
|------|-------------|
| `SCRAPE_PRODUCT` | Scrapes product info from URL |
| `GENERATE_CAMPAIGN` | Generates campaign storyboard |

### Image Role

| Role | Description |
|------|-------------|
| `HERO` | Main hero image |
| `DETAIL` | Detail shot |
| `LIFESTYLE` | Lifestyle context |
| `PACKAGING` | Product packaging |

---

## API Endpoints

Base Path: `/api`

### Assets Endpoints

| Method | Endpoint | Handler | Description |
|--------|----------|---------|-------------|
| POST | `/assets` | `CreateAsset` | Create new asset with name, price, product URL, image count |
| POST | `/assets/import` | `ImportAsset` | Import asset from product URL (creates SCRAPE_PRODUCT job) |
| GET | `/assets` | `ListAssets` | List assets with pagination (query: `limit`, `cursor`) |
| GET | `/assets/:id` | `GetAsset` | Get asset details by ID |
| POST | `/assets/:id/confirm` | `ConfirmUpload` | Confirm upload completion, verify images exist in S3 |
| GET | `/assets/:id/images/upload-url` | `GetImageUploadUrl` | Get presigned S3 URL for image upload |
| DELETE | `/assets/:id/images/:imageId` | `DeleteAssetImage` | Delete an image from asset |

### Campaigns Endpoints

| Method | Endpoint | Handler | Description |
|--------|----------|---------|-------------|
| POST | `/campaigns` | `CreateCampaign` | Create new campaign (creates GENERATE_CAMPAIGN job) |
| GET | `/campaigns` | `ListCampaigns` | List campaigns with pagination |
| GET | `/campaigns/:id` | `GetCampaign` | Get campaign details with storyboard |
| PATCH | `/campaigns/:id` | `UpdateCampaign` | Update campaign storyboard |
| DELETE | `/campaigns/:id` | `DeleteCampaign` | Delete a campaign |

---

### API Examples

#### POST /assets/import

Import a product from URL - creates asset with IMPORTING status and a background job to scrape product info.

**Request:**
```json
{
  "productUrl": "https://example.com/product/123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "assetId": "uuid",
    "jobId": "uuid"
  }
}
```

#### POST /assets

**Request:**
```json
{
  "name": "iPhone 15",
  "price": 999.99,
  "productUrl": "https://example.com/iphone15",
  "imageCount": 3
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "assetId": "uuid",
    "status": "UPLOADING",
    "uploads": [
      {
        "fileKey": "assets/uuid/1.jpg",
        "uploadUrl": "https://s3...",
        "order": 1,
        "expiresIn": 300
      }
    ]
  }
}
```

#### POST /campaigns

**Request:**
```json
{
  "assetId": "uuid"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "assetId": "uuid",
    "status": "DRAFT",
    "storyboard": {},
    "createdAt": "2024-01-01T00:00:00Z"
  }
}
```

#### PATCH /campaigns/:id

**Request:**
```json
{
  "storyboard": {
    "headline": "New Headline",
    "bodyCopy": "New body text",
    "cta": "Shop Now",
    "slides": [
      {
        "slideNumber": 1,
        "imageRole": "HERO",
        "s3Key": "assets/asset-id/1.jpg",
        "overlayText": "Amazing Product"
      }
    ]
  }
}
```

---

## Configuration

### Environment Variables

| Variable | Description |
|----------|-------------|
| `IS_LOCAL` | Set to run locally (non-Lambda mode) |
| `AWS_REGION` | AWS region for DynamoDB/S3 |
| `AWS_ACCESS_KEY_ID` | AWS access key |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key |
| `DYNAMODB_ENDPOINT` | Local DynamoDB endpoint (for local dev) |
| `S3_BUCKET` | S3 bucket name for uploads |
| `CDN_BASE_URL` | Base URL for CDN |

---

## Local Development

### Running Service

```bash
# Assets service (http://localhost:3000)
go run cmd/assets/main.go
```

### Swagger UI

When running locally, access Swagger at:
- http://localhost:3000/swagger/index.html

### Build & Verify

```bash
go build ./...
go vet ./...
```

---

## Architecture Patterns

### Data Flow

```
HTTP Request → Gin Controller → Service → Repository → DynamoDB
                                     ↓
                               S3 (for uploads)
```

### DynamoDB Transactions

For atomic operations involving multiple entities, use `TransactWriteItems`:

```go
baseRepo := repositories.NewBaseRepository(dbClient)
if err := baseRepo.TransactWriteItems(ctx, entity1, entity2); err != nil {
    return response.Internal("failed to create: " + err.Error())
}
```

### Dependency Injection

Services are wired in modules:

```go
func RegisterRoutes(r *gin.RouterGroup, dbClient *dynamodb.Client, store storage.Storage) {
    assetRepo := repositories.NewAssetRepository(dbClient)
    imageRepo := repositories.NewImageRepository(dbClient)
    svc := services.NewAssetService(assetRepo, imageRepo, store)
    ctrl := NewAssetController(svc)
    // register routes...
}
```

### Utils

- `utils.BuildPK(entityType, id)` - Build composite primary key
- `utils.BuildAssetBaseItem(id)` - Build BaseItem for Asset
- `utils.BuildJobBaseItem(id)` - Build BaseItem for Job
- `utils.BuildCampaignBaseItem(id)` - Build BaseItem for Campaign
- `utils.Now()` - Get current timestamp in RFC3339 format
- `utils.MapToDTO(source, dest)` - Copy struct fields
- `utils.MapToDTOs(sourceSlice, destSlice)` - Copy slice of structs

### Response Format

All responses use `BaseResponse[T]`:

| Function | Usage |
|----------|-------|
| `core.Success(data)` | Success with data |
| `core.SuccessWithMeta(data, meta)` | Success with pagination |
| `core.Fail(err)` | Error response |

### Error Handling

Errors use `*core.AppError` with standard HTTP status codes:
- `core.BadRequest()` - 400
- `core.NotFound()` - 404
- `core.Internal()` - 500

---

## Naming Conventions

| Layer | File Pattern | Example |
|-------|--------------|---------|
| Service | `<name>.service.go` | `campaign.service.go` |
| Repository | `<name>.repository.go` | `campaign.repository.go` |
| Request DTO | `<action>.<name>.req.dto.go` | `import.req.dto.go` |
| Response DTO | `<action>.<name>.res.dto.go` | `import.res.dto.go` |
| Model | `<name>.model.go` | `campaigns.model.go` |

- Use **singular** for entity names
- DTOs are use-case based
- Models include `dynamodbav` tags for DynamoDB mapping