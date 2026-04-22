package services

import (
	"context"
	"encoding/json"
	"net/http"

	clerkUser "github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/gianghp123/Vidmerce/backend/internal/configs"
	"github.com/gianghp123/Vidmerce/backend/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	svix "github.com/svix/svix-webhooks/go"
	"go.uber.org/zap"
)

type WebhookService interface {
	HandleClerkWebhook(ctx context.Context, payload []byte, headers http.Header) *response.AppError
}

type ClerkEvent struct {
	Data json.RawMessage `json:"data"`
	Type string          `json:"type"`
}

type webhookService struct {
	clerkWebhookSecret string
}

func NewWebhookService(webhookSecret string) WebhookService {
	return &webhookService{clerkWebhookSecret: webhookSecret}
}

func (s *webhookService) HandleClerkWebhook(ctx context.Context, payload []byte, headers http.Header) *response.AppError {
	logger := configs.GetLogger()

	// 1. Initialize Svix with your 'whsec_...' secret
	wh, err := svix.NewWebhook(s.clerkWebhookSecret)
	if err != nil {
		logger.Error("Failed to initialize Svix", zap.Error(err))
		return response.Internal("failed to initialize webhook")
	}
	// 2. Verify the headers (svix-id, svix-timestamp, svix-signature)
	// This prevents replay attacks and ensures the request is from Clerk
	if err := wh.Verify(payload, headers); err != nil {
		logger.Warn("Webhook verification failed", zap.Error(err))
		return response.BadRequest("verification failed")
	}

	// 3. Unmarshal the outer event
	var evt ClerkEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		return response.BadRequest("invalid json")
	}

	// 4. Handle specific event types
	switch evt.Type {
	case "user.created":
		return s.handleUserCreated(ctx, evt.Data)
	default:
		logger.Info("Unsupported webhook event", zap.String("type", evt.Type))
		return response.BadRequest("Unsupported webhook event")
	}
}

func (s *webhookService) handleUserCreated(ctx context.Context, data json.RawMessage) *response.AppError {
	logger := configs.GetLogger()

	// Define just the fields we need from the 'data' object
	var user struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(data, &user); err != nil {
		return response.BadRequest("invalid user data")
	}

	logger.Info("New user created in Clerk", zap.String("clerkId", user.ID))

	roleMeta := map[string]interface{}{
		"role": string(enums.UserRoleUser),
	}
	metaJSON, _ := json.Marshal(roleMeta)
	meta := json.RawMessage(metaJSON)

	_, err := clerkUser.Update(ctx, user.ID, &clerkUser.UpdateParams{
		PublicMetadata: &meta,
	})

	if err != nil {
		logger.Error("Failed to sync role to Clerk", zap.Error(err))
		return response.Internal("failed to sync role to Clerk")
	}

	return nil
}
