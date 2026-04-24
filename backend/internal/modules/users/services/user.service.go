package services

import (
	"context"

	"github.com/gianghp123/Vidmerce/backend/internal/configs"
	"github.com/gianghp123/Vidmerce/backend/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/auth/guards"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/users/dtos/req"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/users/dtos/res"
	userRepo "github.com/gianghp123/Vidmerce/backend/internal/modules/users/repositories"
	"github.com/gianghp123/Vidmerce/backend/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService interface {
	CreateUser(ctx context.Context, req req.CreateUserReq) (*res.UserRes, *response.AppError)
	GetUser(ctx context.Context, id string) (*res.UserRes, *response.AppError)
	ListUsers(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.UserRes], *response.AppError)
	UpdateUser(ctx context.Context, id string, req req.UpdateUserReq) (*res.UserRes, *response.AppError)
	DeleteUser(ctx context.Context, id string) *response.AppError
}

type userService struct {
	userRepo userRepo.UserRepository
}

func NewUserService(userRepo userRepo.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(ctx context.Context, req req.CreateUserReq) (*res.UserRes, *response.AppError) {
	log := configs.GetLogger()

	userID := uuid.New().String()
	role := enums.UserRole(req.Role)
	if role == "" {
		role = enums.UserRoleUser
	}

	user := models.UserEntity{
		BaseItem:  utils.BuildUserBaseItem(userID),
		Email:     req.Email,
		Role:      role,
		CreatedAt: utils.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		log.Error("Failed to create user", zap.String("userId", userID), zap.Error(err))
		return nil, response.Internal("failed to create user: " + err.Error())
	}

	log.Info("User created", zap.String("userId", userID), zap.String("email", req.Email))
	var result res.UserRes

	err := utils.MapToDTO(user, result)
	if err != nil {
		log.Error("Failed to map user to DTO", zap.Error(err))
		return nil, response.Internal("failed to map user to DTO")
	}
	return &result, nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*res.UserRes, *response.AppError) {
	log := configs.GetLogger()

	if appErr := guards.GuardOwn(ctx, id); appErr != nil {
		return nil, appErr
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		log.Error("Failed to find user", zap.String("userId", id), zap.Error(err))
		return nil, response.Internal("failed to find user")
	}
	if user == nil {
		return nil, response.NotFound("user not found")
	}

	var result res.UserRes

	err = utils.MapToDTO(user, result)
	if err != nil {
		log.Error("Failed to map user to DTO", zap.Error(err))
		return nil, response.Internal("failed to map user to DTO")
	}
	return &result, nil
}

func (s *userService) ListUsers(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.UserRes], *response.AppError) {
	log := configs.GetLogger()

	auth, err := guards.FromAuthContext(ctx)
	if err != nil {
		return nil, response.Unauthorized("authentication required")
	}

	// Only admins can list users
	if auth.Role != enums.UserRoleAdmin {
		return nil, response.Forbidden("admin access required")
	}

	if limit <= 0 {
		limit = 20
	}

	result, err := s.userRepo.FindAll(ctx, limit, cursor)
	if err != nil {
		log.Error("Failed to fetch users", zap.Error(err))
		return nil, response.Internal("failed to fetch users: " + err.Error())
	}

	var users []res.UserRes

	err = utils.MapToDTO(result.Data, users)
	if err != nil {
		log.Error("Failed to map user to DTO", zap.Error(err))
		return nil, response.Internal("failed to map user to DTO")
	}

	log.Debug("Users listed", zap.Int("count", len(users)), zap.Bool("hasMore", result.Meta.HasMore))
	return &response.PaginatedResult[res.UserRes]{
		Data: users,
		Meta: result.Meta,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, req req.UpdateUserReq) (*res.UserRes, *response.AppError) {
	log := configs.GetLogger()

	if appErr := guards.GuardOwn(ctx, id); appErr != nil {
		return nil, appErr
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		log.Error("Failed to find user", zap.String("userId", id), zap.Error(err))
		return nil, response.Internal("failed to find user")
	}
	if user == nil {
		return nil, response.NotFound("user not found")
	}

	if req.Role != nil && *req.Role != "" {
		user.Role = enums.UserRole(*req.Role)
	}

	if err := s.userRepo.Update(ctx, *user); err != nil {
		log.Error("Failed to update user", zap.String("userId", id), zap.Error(err))
		return nil, response.Internal("failed to update user")
	}

	log.Info("User updated", zap.String("userId", id))

	var result res.UserRes

	err = utils.MapToDTO(user, result)
	if err != nil {
		log.Error("Failed to map user to DTO", zap.Error(err))
		return nil, response.Internal("failed to map user to DTO")
	}
	return &result, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) *response.AppError {
	log := configs.GetLogger()

	if appErr := guards.GuardOwn(ctx, id); appErr != nil {
		return appErr
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		log.Error("Failed to find user", zap.String("userId", id), zap.Error(err))
		return response.Internal("failed to find user")
	}
	if user == nil {
		return response.NotFound("user not found")
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		log.Error("Failed to delete user", zap.String("userId", id), zap.Error(err))
		return response.Internal("failed to delete user")
	}

	log.Info("User deleted", zap.String("userId", id))
	return nil
}
