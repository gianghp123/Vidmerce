package services

import (
	"context"

	"github.com/gianghp123/Vidmerce/backend/services/internal/core/response"
	res "github.com/gianghp123/Vidmerce/backend/services/internal/modules/assets/dtos/res"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/assets/repositories"
)

type AssetService interface {
	GenerateUploadUrls(ctx context.Context, count int) (*res.GenerateUploadUrlsRes, *response.AppError)
	CreateAsset(ctx context.Context, req interface{}) (*res.CreateAssetRes, *response.AppError)
	ListAssets(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.AssetRes], *response.AppError)
}

type assetService struct {
	repo repositories.AssetRepository
}

func NewAssetService(repo repositories.AssetRepository) AssetService {
	return &assetService{repo: repo}
}

func (s *assetService) GenerateUploadUrls(ctx context.Context, count int) (*res.GenerateUploadUrlsRes, *response.AppError) {
	return nil, response.Internal("not implemented")
}

func (s *assetService) CreateAsset(ctx context.Context, req interface{}) (*res.CreateAssetRes, *response.AppError) {
	return nil, response.Internal("not implemented")
}

func (s *assetService) ListAssets(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.AssetRes], *response.AppError) {
	if limit <= 0 {
		limit = 20
	}

	result, err := s.repo.FindAll(ctx, limit, cursor)
	if err != nil {
		return nil, response.Internal("failed to fetch assets")
	}

	assets := make([]res.AssetRes, 0, len(result.Data))
	for _, item := range result.Data {
		assets = append(assets, res.AssetRes{
			AssetID:    item.PK,
			Name:       item.Name,
			Price:      item.Price,
			ImageURL:   item.ImageURL,
			ProductURL: item.ProductURL,
			CreatedAt:  item.CreatedAt,
		})
	}

	return &response.PaginatedResult[res.AssetRes]{
		Data: assets,
		Meta: result.Meta,
	}, nil
}
