package services

import (
	"context"
	"encoding/json"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkUser "github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gianghp123/Vidmerce/backend/internal/configs"
	"github.com/gianghp123/Vidmerce/backend/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/auth/dtos/req"
	userRepo "github.com/gianghp123/Vidmerce/backend/internal/modules/users/repositories"
	"github.com/gianghp123/Vidmerce/backend/internal/utils"
	"go.uber.org/zap"
)

type AuthService struct {
	userRepo userRepo.UserRepository
}

func NewAuthService(userRepo userRepo.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) SignUp(ctx context.Context, signUpReq req.SignUpReq) (*clerk.User, *response.AppError) {
	logger := configs.GetLogger()
	logger.Info("Signing up user", zap.String("email", signUpReq.Email))

	publicMetadata := map[string]any{
		"role": string(enums.UserRoleUser),
	}

	b, err := json.Marshal(publicMetadata)
	if err != nil {
		logger.Error("Failed to marshal user metadata", zap.Error(err))
		return nil, response.Internal("failed to marshal user metadata")
	}

	meta := json.RawMessage(b)

	newUser, err := clerkUser.Create(ctx, &clerkUser.CreateParams{
		EmailAddresses: &[]string{signUpReq.Email},
		Password:       &signUpReq.Password,
		PublicMetadata: &meta,
	})
	if err != nil {
		logger.Error("Failed to create user", zap.Error(err))
		return nil, response.Internal("failed to create user")
	}

	userEntity := models.UserEntity{
		BaseItem:  utils.BuildUserBaseItem(newUser.ID),
		Email:     signUpReq.Email,
		Role:      enums.UserRoleUser,
		CreatedAt: utils.Now(),
	}

	if err := s.userRepo.Create(ctx, userEntity); err != nil {
		logger.Error("Failed to create user record", zap.Error(err), zap.String("clerk_user_id", newUser.ID))

		if delErr := deleteClerkUser(ctx, newUser.ID); delErr != nil {
			logger.Error("Failed to delete Clerk user after DB insert failure",
				zap.Error(delErr),
				zap.String("clerk_user_id", newUser.ID),
			)
		}

		return nil, response.Internal("failed to create user record")
	}

	return newUser, nil
}

func deleteClerkUser(ctx context.Context, userID string) error {
	_, err := clerkUser.Delete(ctx, userID)
	return err
}
