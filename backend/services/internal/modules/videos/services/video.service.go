package services

import (
	"context"

	"github.com/gianghp123/Vidmerce/backend/services/internal/core"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/videos/dtos/req"
	res "github.com/gianghp123/Vidmerce/backend/services/internal/modules/videos/dtos/res"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/videos/repositories"
)

type VideoService interface {
	CreateVideo(ctx context.Context, req interface{}) (*res.CreateVideoRes, *core.AppError)
	ListVideos(ctx context.Context, limit int, cursor string) (*core.PaginatedResult[res.VideoResponse], *core.AppError)
	GetVideo(ctx context.Context, id string) (*res.VideoDetailResponse, *core.AppError)
}

type videoService struct {
	repo repositories.VideoRepository
}

func NewVideoService(repo repositories.VideoRepository) VideoService {
	return &videoService{repo: repo}
}

func (s *videoService) CreateVideo(ctx context.Context, req interface{}) (*res.CreateVideoRes, *core.AppError) {
	return nil, core.Internal("not implemented")
}

func (s *videoService) ListVideos(ctx context.Context, limit int, cursor string) (*core.PaginatedResult[res.VideoResponse], *core.AppError) {
	if limit <= 0 {
		limit = 10
	}

	query := req.GetVideosQuery{Limit: limit, LastKey: cursor}
	items, lastKey, hasMore, err := s.repo.FindAll(ctx, query.Limit, query.LastKey)
	if err != nil {
		return nil, core.Internal("failed to fetch videos")
	}

	videos := make([]res.VideoResponse, 0, len(items))
	for _, item := range items {
		videos = append(videos, res.VideoResponse{
			ID:           item.PK,
			Title:        item.Title,
			Status:       item.Status,
			VideoURL:     item.VideoURL,
			ThumbnailURL: item.ThumbnailURL,
			CreatedAt:    item.CreatedAt,
		})
	}

	return &core.PaginatedResult[res.VideoResponse]{
		Data: videos,
		Meta: core.NewCursorMeta(limit, lastKey, hasMore),
	}, nil
}

func (s *videoService) GetVideo(ctx context.Context, id string) (*res.VideoDetailResponse, *core.AppError) {
	param := req.GetVideoByIDParam{ID: id}
	video, err := s.repo.FindByID(ctx, param.ID)
	if err != nil {
		return nil, core.Internal("failed to fetch video")
	}
	if video == nil {
		return nil, core.NotFound("video not found")
	}

	interactive, err := s.repo.FindInteractiveByVideoID(ctx, param.ID)
	if err != nil {
		return nil, core.Internal("failed to fetch video interactivity")
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
