package services

import (
	"context"

	"github.com/gianghp123/Vidmerce/backend/internal/configs"
	"github.com/gianghp123/Vidmerce/backend/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
	"github.com/gianghp123/Vidmerce/backend/internal/database/repositories"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/campaigns/dtos/req"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/campaigns/dtos/res"
	campaignRepo "github.com/gianghp123/Vidmerce/backend/internal/modules/campaigns/repositories"
	"github.com/gianghp123/Vidmerce/backend/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CampaignService interface {
	CreateCampaign(ctx context.Context, req req.CreateCampaignReq) (*res.CampaignRes, *response.AppError)
	GetCampaign(ctx context.Context, id string) (*res.CampaignRes, *response.AppError)
	ListCampaigns(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.CampaignRes], *response.AppError)
	UpdateCampaign(ctx context.Context, id string, req req.UpdateCampaignReq) (*res.CampaignRes, *response.AppError)
	DeleteCampaign(ctx context.Context, id string) *response.AppError
}

type campaignService struct {
	campaignRepo campaignRepo.CampaignRepository
}

func NewCampaignService(campaignRepo campaignRepo.CampaignRepository) CampaignService {
	return &campaignService{campaignRepo: campaignRepo}
}

func (s *campaignService) CreateCampaign(ctx context.Context, req req.CreateCampaignReq) (*res.CampaignRes, *response.AppError) {
	log := configs.GetLogger()

	campaignID := uuid.New().String()
	jobID := uuid.New().String()

	campaign := models.CampaignEntity{
		BaseItem:  utils.BuildCampaignBaseItem(campaignID),
		AssetID:   req.AssetID,
		Status:    enums.CampaignStatusDraft,
		CreatedAt: utils.Now(),
	}

	job := models.JobEntity{
		BaseItem:  utils.BuildJobBaseItem(jobID),
		TargetID:  &campaignID,
		Type:      enums.JobTypeGenerateCampaign,
		Status:    enums.JobStatusPending,
		CreatedAt: utils.Now(),
	}

	baseRepo := repositories.NewBaseRepository(s.campaignRepo.DBClient())
	if err := baseRepo.TransactWriteItems(ctx, campaign, job); err != nil {
		log.Error("Failed to create campaign", zap.String("campaignId", campaignID), zap.Error(err))
		return nil, response.Internal("failed to create campaign: " + err.Error())
	}

	log.Info("Campaign created", zap.String("campaignId", campaignID), zap.String("assetId", req.AssetID), zap.String("jobId", jobID))
	return toCampaignRes(&campaign), nil
}

func (s *campaignService) GetCampaign(ctx context.Context, id string) (*res.CampaignRes, *response.AppError) {
	log := configs.GetLogger()

	campaign, err := s.campaignRepo.FindByID(ctx, id)
	if err != nil {
		log.Error("Failed to find campaign", zap.String("campaignId", id), zap.Error(err))
		return nil, response.Internal("failed to find campaign")
	}
	if campaign == nil {
		return nil, response.NotFound("campaign not found")
	}

	return toCampaignRes(campaign), nil
}

func (s *campaignService) ListCampaigns(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.CampaignRes], *response.AppError) {
	log := configs.GetLogger()

	if limit <= 0 {
		limit = 20
	}

	result, err := s.campaignRepo.FindAll(ctx, limit, cursor)
	if err != nil {
		log.Error("Failed to fetch campaigns", zap.Error(err))
		return nil, response.Internal("failed to fetch campaigns: " + err.Error())
	}

	campaigns := make([]res.CampaignRes, 0, len(result.Data))
	for _, item := range result.Data {
		campaigns = append(campaigns, *toCampaignRes(&item))
	}

	log.Debug("Campaigns listed", zap.Int("count", len(campaigns)), zap.Bool("hasMore", result.Meta.HasMore))
	return &response.PaginatedResult[res.CampaignRes]{
		Data: campaigns,
		Meta: result.Meta,
	}, nil
}

func (s *campaignService) UpdateCampaign(ctx context.Context, id string, req req.UpdateCampaignReq) (*res.CampaignRes, *response.AppError) {
	log := configs.GetLogger()

	campaign, err := s.campaignRepo.FindByID(ctx, id)
	if err != nil {
		log.Error("Failed to find campaign", zap.String("campaignId", id), zap.Error(err))
		return nil, response.Internal("failed to find campaign")
	}
	if campaign == nil {
		return nil, response.NotFound("campaign not found")
	}

	if req.Storyboard != nil {
		if campaign.StoryboardEntity.Headline == "" {
			campaign.StoryboardEntity.Headline = req.Storyboard.Headline
		}
		if req.Storyboard.Headline != "" {
			campaign.StoryboardEntity.Headline = req.Storyboard.Headline
		}
		if req.Storyboard.BodyCopy != "" {
			campaign.StoryboardEntity.BodyCopy = req.Storyboard.BodyCopy
		}
		if req.Storyboard.Cta != "" {
			campaign.StoryboardEntity.Cta = req.Storyboard.Cta
		}
		if req.Storyboard.Slides != nil {
			campaign.StoryboardEntity.Slides = toModelSlides(req.Storyboard.Slides)
		}
	}

	if err := s.campaignRepo.Update(ctx, *campaign); err != nil {
		log.Error("Failed to update campaign", zap.String("campaignId", id), zap.Error(err))
		return nil, response.Internal("failed to update campaign")
	}

	log.Info("Campaign updated", zap.String("campaignId", id))
	return toCampaignRes(campaign), nil
}

func (s *campaignService) DeleteCampaign(ctx context.Context, id string) *response.AppError {
	log := configs.GetLogger()

	campaign, err := s.campaignRepo.FindByID(ctx, id)
	if err != nil {
		log.Error("Failed to find campaign", zap.String("campaignId", id), zap.Error(err))
		return response.Internal("failed to find campaign")
	}
	if campaign == nil {
		return response.NotFound("campaign not found")
	}

	if err := s.campaignRepo.Delete(ctx, id); err != nil {
		log.Error("Failed to delete campaign", zap.String("campaignId", id), zap.Error(err))
		return response.Internal("failed to delete campaign")
	}

	log.Info("Campaign deleted", zap.String("campaignId", id))
	return nil
}

func toCampaignRes(campaign *models.CampaignEntity) *res.CampaignRes {
	var result res.CampaignRes
	_ = utils.MapToDTO(campaign, &result)
	return &result
}

func toResStoryboard(sb models.StoryboardEntity) res.Storyboard {
	var result res.Storyboard
	_ = utils.MapToDTO(sb, &result)
	return result
}

func toModelSlides(slides []req.Slide) []models.SlideEntity {
	var result []models.SlideEntity
	_ = utils.MapToDTOs(slides, &result)
	return result
}
