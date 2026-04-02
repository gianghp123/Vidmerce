package services

import (
	"context"

	"github.com/gianghp123/Vidmerce/backend/services/internal/core"
	res "github.com/gianghp123/Vidmerce/backend/services/internal/modules/assets/dtos/res"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/assets/repositories"
)

type AssetService interface {
	GenerateUploadUrls(ctx context.Context, count int) (*res.GenerateUploadUrlsRes, *core.AppError)
	CreateAsset(ctx context.Context, req interface{}) (*res.CreateAssetRes, *core.AppError)
	ListAssets(ctx context.Context, limit int, cursor string) (*core.PaginatedResult[res.AssetRes], *core.AppError)
}

type assetService struct {
	repo repositories.AssetRepository
}

func NewAssetService(repo repositories.AssetRepository) AssetService {
	return &assetService{repo: repo}
}

func (s *assetService) GenerateUploadUrls(ctx context.Context, count int) (*res.GenerateUploadUrlsRes, *core.AppError) {
	return nil, core.Internal("not implemented")
}

func (s *assetService) CreateAsset(ctx context.Context, req interface{}) (*res.CreateAssetRes, *core.AppError) {
	return nil, core.Internal("not implemented")
}

func (s *assetService) ListAssets(ctx context.Context, limit int, cursor string) (*core.PaginatedResult[res.AssetRes], *core.AppError) {
	if limit <= 0 {
		limit = 20
	}

	items, lastKey, hasMore, err := s.repo.FindAll(ctx, limit, cursor)
	if err != nil {
		return nil, core.Internal("failed to fetch assets")
	}

	assets := make([]res.AssetRes, 0, len(items))
	for _, item := range items {
		assets = append(assets, res.AssetRes{
			AssetID:    item.PK,
			Name:       item.Name,
			Price:      item.Price,
			ImageURL:   item.ImageURL,
			ProductURL: item.ProductURL,
			CreatedAt:  item.CreatedAt,
		})
	}

	return &core.PaginatedResult[res.AssetRes]{
		Data: assets,
		Meta: core.NewCursorMeta(limit, lastKey, hasMore),
	}, nil
}
