package services

import (
	"context"
	"fmt"
	"time"

	"github.com/gianghp123/Vidmerce/backend/services/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/services/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/services/internal/database/models"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/assets/dtos/req"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/assets/dtos/res"
	imageRepo "github.com/gianghp123/Vidmerce/backend/services/internal/modules/assets/repositories"
	"github.com/gianghp123/Vidmerce/backend/services/internal/storage"
	"github.com/gianghp123/Vidmerce/backend/services/internal/utils"
	"github.com/google/uuid"
)

const MaxImagesPerAsset = 5

type AssetService interface {
	CreateAsset(ctx context.Context, req req.CreateAssetReq) (*res.CreateAssetRes, *response.AppError)
	ConfirmUpload(ctx context.Context, assetID string) (*res.ConfirmAssetRes, *response.AppError)
	GetAsset(ctx context.Context, assetID string) (*res.AssetRes, *response.AppError)
	ListAssets(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.AssetRes], *response.AppError)
}

type assetService struct {
	assetRepo imageRepo.AssetRepository
	imageRepo imageRepo.ImageRepository
	storage   storage.Storage
}

func NewAssetService(assetRepo imageRepo.AssetRepository, imageRepo imageRepo.ImageRepository, storage storage.Storage) AssetService {
	return &assetService{assetRepo: assetRepo, imageRepo: imageRepo, storage: storage}
}

func (s *assetService) CreateAsset(ctx context.Context, req req.CreateAssetReq) (*res.CreateAssetRes, *response.AppError) {
	imageCount := req.ImageCount
	if imageCount <= 0 {
		imageCount = 1
	}
	if imageCount > MaxImagesPerAsset {
		return nil, response.BadRequest(fmt.Sprintf("maximum %d images per asset", MaxImagesPerAsset))
	}

	assetID := uuid.New().String()
	now := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	asset := models.AssetEntity{
		BaseItem: models.BaseItem{
			PK:     "ASSET#" + assetID,
			SK:     "METADATA",
			GSI1PK: "ENTITY#ASSET",
			GSI1SK: assetID,
		},
		Name:       req.Name,
		Price:      req.Price,
		ProductURL: req.ProductURL,
		Status:     enums.StatusAssetUploading,
		ImageCount: imageCount,
		CreatedAt:  now,
	}

	uploads := make([]res.UploadInfo, 0, imageCount)
	images := make([]models.ImageEntity, 0, imageCount)

	for i := 1; i <= imageCount; i++ {
		fileKey := fmt.Sprintf("assets/%s/%d.jpg", assetID, i)

		uploadURL, err := s.storage.GeneratePresignedUploadURL(ctx, fileKey, "image/jpeg", 5*time.Minute)
		if err != nil {
			return nil, response.Internal("failed to generate upload URL")
		}

		uploads = append(uploads, res.UploadInfo{
			FileKey:   fileKey,
			UploadURL: uploadURL,
			Order:     i,
			ExpiresIn: 300,
		})

		images = append(images, models.ImageEntity{
			BaseItem: models.BaseItem{
				PK: "ASSET#" + assetID,
				SK: fmt.Sprintf("IMAGE#%d", i),
			},
			FileKey: fileKey,
			Status:  enums.StatusImageUploading,
			Order:   i,
		})
	}

	err := s.imageRepo.Create(ctx, images)
	if err != nil {
		return nil, response.Internal("failed to create image records")
	}

	err = s.assetRepo.Create(ctx, asset)
	if err != nil {
		return nil, response.Internal("failed to create asset")
	}

	return &res.CreateAssetRes{
		AssetID: assetID,
		Status:  string(enums.StatusAssetUploading),
		Uploads: uploads,
	}, nil
}

func (s *assetService) ConfirmUpload(ctx context.Context, assetID string) (*res.ConfirmAssetRes, *response.AppError) {
	asset, err := s.assetRepo.FindByID(ctx, assetID)
	if err != nil {
		return nil, response.Internal("failed to find asset")
	}
	if asset == nil {
		return nil, response.NotFound("asset not found")
	}

	if asset.Status != enums.StatusAssetUploading {
		return nil, response.BadRequest("asset is not in UPLOADING state")
	}

	imageCount := asset.ImageCount
	images := make([]res.ImageInfo, 0, imageCount)
	uploadedCount := 0

	for i := 1; i <= imageCount; i++ {
		fileKey := fmt.Sprintf("assets/%s/%d.jpg", assetID, i)

		exists, err := s.storage.ObjectExists(ctx, fileKey)
		if err != nil {
			return nil, response.Internal("failed to verify upload")
		}

		if !exists {
			return nil, response.BadRequest(fmt.Sprintf("image %d/%d not uploaded", i, imageCount))
		}

		imageUrl := utils.GetCDNURL(fileKey)

		err = s.imageRepo.UpdateStatus(ctx, assetID, i, string(enums.StatusImageCompleted))
		if err != nil {
			return nil, response.Internal("failed to update image status")
		}

		images = append(images, res.ImageInfo{
			ImageURL: imageUrl,
			Order:    i,
		})
		uploadedCount++
	}

	err = s.assetRepo.UpdateAssetStatus(ctx, assetID, string(enums.StatusAssetCompleted))
	if err != nil {
		return nil, response.Internal("failed to update asset status")
	}

	return &res.ConfirmAssetRes{
		AssetID: assetID,
		Status:  string(enums.StatusAssetCompleted),
		Images:  images,
	}, nil
}

func (s *assetService) GetAsset(ctx context.Context, assetID string) (*res.AssetRes, *response.AppError) {
	asset, err := s.assetRepo.FindByID(ctx, assetID)
	if err != nil {
		return nil, response.Internal("failed to find asset")
	}
	if asset == nil {
		return nil, response.NotFound("asset not found")
	}

	images, err := s.imageRepo.FindByAssetID(ctx, assetID)
	if err != nil {
		return nil, response.Internal("failed to find images")
	}

	imageInfos := make([]res.ImageInfo, 0, len(images))
	for _, img := range images {
		if img.Status == enums.StatusImageCompleted {
			imageURL := utils.GetCDNURL(img.FileKey)
			imageInfos = append(imageInfos, res.ImageInfo{
				ImageURL: imageURL,
				Order:    img.Order,
			})
		}
	}

	assetIDOnly := asset.PK
	if len(assetIDOnly) > 7 && assetIDOnly[:7] == "ASSET#" {
		assetIDOnly = assetIDOnly[7:]
	}

	return &res.AssetRes{
		AssetID:    assetIDOnly,
		Name:       asset.Name,
		Price:      asset.Price,
		Images:     imageInfos,
		ProductURL: asset.ProductURL,
		Status:     string(asset.Status),
		CreatedAt:  asset.CreatedAt,
	}, nil
}

func (s *assetService) ListAssets(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.AssetRes], *response.AppError) {
	if limit <= 0 {
		limit = 20
	}

	result, err := s.assetRepo.FindAll(ctx, limit, cursor)
	if err != nil {
		return nil, response.Internal("failed to fetch assets")
	}

	assets := make([]res.AssetRes, 0, len(result.Data))
	for _, item := range result.Data {
		assetID := item.PK
		if len(assetID) > 7 && assetID[:7] == "ASSET#" {
			assetID = assetID[7:]
		}

		images, err := s.imageRepo.FindByAssetID(ctx, assetID)
		if err != nil {
			continue
		}

		var imageInfos []res.ImageInfo
		for _, img := range images {
			if img.Status == enums.StatusImageCompleted && img.Order == 1 {
				imageURL := utils.GetCDNURL(img.FileKey)
				imageInfos = append(imageInfos, res.ImageInfo{
					ImageURL: imageURL,
					Order:    img.Order,
				})
				break
			}
		}

		assets = append(assets, res.AssetRes{
			AssetID:    assetID,
			Name:       item.Name,
			Price:      item.Price,
			Images:     imageInfos,
			ProductURL: item.ProductURL,
			Status:     string(item.Status),
			CreatedAt:  item.CreatedAt,
		})
	}

	return &response.PaginatedResult[res.AssetRes]{
		Data: assets,
		Meta: result.Meta,
	}, nil
}
