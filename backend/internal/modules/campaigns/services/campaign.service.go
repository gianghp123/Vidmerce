package services

import (
	"context"

	"github.com/gianghp123/Vidmerce/backend/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
	"github.com/gianghp123/Vidmerce/backend/internal/database/repositories"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/campaigns/dtos/req"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/campaigns/dtos/res"
	campaignRepo "github.com/gianghp123/Vidmerce/backend/internal/modules/campaigns/repositories"
	"github.com/gianghp123/Vidmerce/backend/internal/utils"
	"github.com/google/uuid"
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
	campaignID := uuid.New().String()
	jobID := uuid.New().String()

	campaign := models.CampaignEntity{
		BaseItem:  utils.BuildCampaignBaseItem(campaignID),
		AssetID:   req.AssetID,
		Status:    enums.StatusCampaignDraft,
		CreatedAt: utils.Now(),
	}

	job := models.JobEntity{
		BaseItem:  utils.BuildJobBaseItem(jobID),
		TargetID:  campaignID,
		Type:      enums.TypeJobGenerateCampaign,
		Status:    enums.StatusJobPending,
		CreatedAt: utils.Now(),
	}

	baseRepo := repositories.NewBaseRepository(s.campaignRepo.DBClient())
	if err := baseRepo.TransactWriteItems(ctx, campaign, job); err != nil {
		return nil, response.Internal("failed to create campaign: " + err.Error())
	}

	return toCampaignRes(&campaign), nil
}

func (s *campaignService) GetCampaign(ctx context.Context, id string) (*res.CampaignRes, *response.AppError) {
	campaign, err := s.campaignRepo.FindByID(ctx, id)
	if err != nil {
		return nil, response.Internal("failed to find campaign")
	}
	if campaign == nil {
		return nil, response.NotFound("campaign not found")
	}

	return toCampaignRes(campaign), nil
}

func (s *campaignService) ListCampaigns(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.CampaignRes], *response.AppError) {
	if limit <= 0 {
		limit = 20
	}

	result, err := s.campaignRepo.FindAll(ctx, limit, cursor)
	if err != nil {
		return nil, response.Internal("failed to fetch campaigns: " + err.Error())
	}

	campaigns := make([]res.CampaignRes, 0, len(result.Data))
	for _, item := range result.Data {
		campaigns = append(campaigns, *toCampaignRes(&item))
	}

	return &response.PaginatedResult[res.CampaignRes]{
		Data: campaigns,
		Meta: result.Meta,
	}, nil
}

func (s *campaignService) UpdateCampaign(ctx context.Context, id string, req req.UpdateCampaignReq) (*res.CampaignRes, *response.AppError) {
	campaign, err := s.campaignRepo.FindByID(ctx, id)
	if err != nil {
		return nil, response.Internal("failed to find campaign")
	}
	if campaign == nil {
		return nil, response.NotFound("campaign not found")
	}

	if req.Storyboard != nil {
		if campaign.Storyboard.Headline == "" {
			campaign.Storyboard.Headline = req.Storyboard.Headline
		}
		if req.Storyboard.Headline != "" {
			campaign.Storyboard.Headline = req.Storyboard.Headline
		}
		if req.Storyboard.BodyCopy != "" {
			campaign.Storyboard.BodyCopy = req.Storyboard.BodyCopy
		}
		if req.Storyboard.CTA != "" {
			campaign.Storyboard.CTA = req.Storyboard.CTA
		}
		if req.Storyboard.Slides != nil {
			campaign.Storyboard.Slides = toModelSlides(req.Storyboard.Slides)
		}
	}

	if err := s.campaignRepo.Update(ctx, *campaign); err != nil {
		return nil, response.Internal("failed to update campaign")
	}

	return toCampaignRes(campaign), nil
}

func (s *campaignService) DeleteCampaign(ctx context.Context, id string) *response.AppError {
	campaign, err := s.campaignRepo.FindByID(ctx, id)
	if err != nil {
		return response.Internal("failed to find campaign")
	}
	if campaign == nil {
		return response.NotFound("campaign not found")
	}

	if err := s.campaignRepo.Delete(ctx, id); err != nil {
		return response.Internal("failed to delete campaign")
	}

	return nil
}

func toCampaignRes(campaign *models.CampaignEntity) *res.CampaignRes {
	var result res.CampaignRes
	_ = utils.MapToDTO(campaign, &result)
	result.Status = string(campaign.Status)
	result.Storyboard = toResStoryboard(campaign.Storyboard)
	return &result
}

func toResStoryboard(sb models.Storyboard) res.Storyboard {
	var result res.Storyboard
	_ = utils.MapToDTO(sb, &result)
	if len(sb.Slides) > 0 {
		result.Slides = toResSlides(sb.Slides)
	}
	return result
}

func toResSlides(slides []models.Slide) []res.Slide {
	var result []res.Slide
	_ = utils.MapToDTOs(slides, &result)
	for i := range result {
		result[i].ImageRole = string(slides[i].ImageRole)
	}
	return result
}

func toModelSlides(slides []req.Slide) []models.Slide {
	var result []models.Slide
	_ = utils.MapToDTOs(slides, &result)
	return result
}
