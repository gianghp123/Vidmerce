package services

import (
	"context"

	"github.com/gianghp123/Vidmerce/backend/services/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/videos/dtos/req"
	res "github.com/gianghp123/Vidmerce/backend/services/internal/modules/videos/dtos/res"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/videos/repositories"
)

type VideoService interface {
	CreateVideo(ctx context.Context, req req.CreateVideoReq) (*res.CreateVideoRes, *response.AppError)
	ListVideos(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.VideoResponse], *response.AppError)
	GetVideo(ctx context.Context, id string) (*res.VideoDetailResponse, *response.AppError)
}

type videoService struct {
	repo repositories.VideoRepository
}

func NewVideoService(repo repositories.VideoRepository) VideoService {
	return &videoService{repo: repo}
}

func (s *videoService) CreateVideo(ctx context.Context, req req.CreateVideoReq) (*res.CreateVideoRes, *response.AppError) {
	return nil, response.Internal("not implemented")
}

func (s *videoService) ListVideos(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.VideoResponse], *response.AppError) {
	if limit <= 0 {
		limit = 10
	}

	query := req.GetVideosQuery{Limit: limit, LastKey: cursor}
	result, err := s.repo.FindAll(ctx, query.Limit, query.LastKey)
	if err != nil {
		return nil, response.Internal("failed to fetch videos")
	}

	videos := make([]res.VideoResponse, 0, len(result.Data))
	for _, item := range result.Data {
		videos = append(videos, res.VideoResponse{
			ID:           item.PK,
			Title:        item.Title,
			Status:       item.Status,
			VideoURL:     item.VideoURL,
			ThumbnailURL: item.ThumbnailURL,
			CreatedAt:    item.CreatedAt,
		})
	}

	return &response.PaginatedResult[res.VideoResponse]{
		Data: videos,
		Meta: result.Meta,
	}, nil
}

func (s *videoService) GetVideo(ctx context.Context, id string) (*res.VideoDetailResponse, *response.AppError) {
	param := req.GetVideoByIDParam{ID: id}
	video, err := s.repo.FindByID(ctx, param.ID)
	if err != nil {
		return nil, response.Internal("failed to fetch video")
	}
	if video == nil {
		return nil, response.NotFound("video not found")
	}

	interactive, err := s.repo.FindInteractiveByVideoID(ctx, param.ID)
	if err != nil {
		return nil, response.Internal("failed to fetch video interactivity")
	}

	hotspots := make([]res.HotspotResponse, 0)
	if interactive != nil {
		for _, h := range interactive.Hotspots {
			hotspots = append(hotspots, res.HotspotResponse{
				X: h.X,
				Y: h.Y,
				Asset: res.AssetSnapshotDTO{
					AssetID:    h.Asset.AssetID,
					Name:       h.Asset.Name,
					Price:      h.Asset.Price,
					ImageURL:   h.Asset.ImageURL,
					ProductURL: h.Asset.ProductURL,
				},
			})
		}
	}

	return &res.VideoDetailResponse{
		VideoResponse: res.VideoResponse{
			ID:           video.PK,
			Title:        video.Title,
			Status:       video.Status,
			VideoURL:     video.VideoURL,
			ThumbnailURL: video.ThumbnailURL,
			CreatedAt:    video.CreatedAt,
		},
		Hotspots: hotspots,
	}, nil
}
