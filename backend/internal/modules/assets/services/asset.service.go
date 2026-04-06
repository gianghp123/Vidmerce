package services

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gianghp123/Vidmerce/backend/internal/core"
	"github.com/gianghp123/Vidmerce/backend/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/assets/dtos/req"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/assets/dtos/res"
	imageRepo "github.com/gianghp123/Vidmerce/backend/internal/modules/assets/repositories"
	"github.com/gianghp123/Vidmerce/backend/internal/storage"
	"github.com/gianghp123/Vidmerce/backend/internal/utils"
	"github.com/google/uuid"
)

type AssetService interface {
	CreateAsset(ctx context.Context, req req.CreateAssetReq) (*res.CreateAssetRes, *response.AppError)
	ConfirmUpload(ctx context.Context, assetID string) (*res.ConfirmAssetRes, *response.AppError)
	GetAsset(ctx context.Context, assetID string) (*res.AssetRes, *response.AppError)
	ListAssets(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.AssetRes], *response.AppError)
	GetImageUploadUrl(ctx context.Context, assetID string, fileName string) (*res.UploadInfo, *response.AppError)
	DeleteAssetImage(ctx context.Context, assetID string, imageID string) *response.AppError
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
	if imageCount > core.MaxImagesPerAsset {
		return nil, response.BadRequest(fmt.Sprintf("maximum %d images per asset", core.MaxImagesPerAsset))
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

	// ✅ get images from DB (source of truth)
	dbImages, err := s.imageRepo.FindByAssetID(ctx, assetID)
	if err != nil {
		return nil, response.Internal("failed to get images")
	}

	images := make([]res.ImageInfo, 0, len(dbImages))
	uploadedCount := 0

	for _, img := range dbImages {
		fileKey := img.FileKey

		exists, err := s.storage.ObjectExists(ctx, fileKey)
		if err != nil {
			return nil, response.Internal("failed to verify upload")
		}

		imageUrl := utils.GetCDNURL(fileKey)

		var status string

		if exists {
			status = string(enums.StatusImageCompleted)
			uploadedCount++
		} else {
			status = string(enums.StatusImageFailed)
		}

		// ✅ update per image
		if err := s.imageRepo.UpdateStatus(ctx, assetID, img.Order, status); err != nil {
			return nil, response.Internal("failed to update image status")
		}

		images = append(images, res.ImageInfo{
			ImageURL: imageUrl,
			Order:    img.Order,
			Status:   status,
		})
	}

	// ✅ determine asset status dynamically
	total := len(dbImages)

	var finalStatus string

	switch {
	case total == 0:
		finalStatus = string(enums.StatusAssetFailed)
	case uploadedCount == total:
		finalStatus = string(enums.StatusAssetCompleted)
	case uploadedCount > 0:
		finalStatus = string(enums.StatusAssetPartial)
	default:
		finalStatus = string(enums.StatusAssetFailed)
	}

	if err := s.assetRepo.UpdateAssetStatus(ctx, assetID, finalStatus); err != nil {
		return nil, response.Internal("failed to update asset status")
	}

	return &res.ConfirmAssetRes{
		AssetID: assetID,
		Status:  finalStatus,
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

	return &res.AssetRes{
		AssetID:    assetID,
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
		return nil, response.Internal("failed to fetch assets: " + err.Error())
	}

	assets := make([]res.AssetRes, 0, len(result.Data))
	for _, item := range result.Data {
		assetID := item.PK

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

func (s *assetService) GetImageUploadUrl(ctx context.Context, assetID string, fileName string) (*res.UploadInfo, *response.AppError) {
	asset, err := s.assetRepo.FindByID(ctx, assetID)
	if err != nil {
		return nil, response.Internal("failed to find asset")
	}
	if asset == nil {
		return nil, response.NotFound("asset not found")
	}

	// 1. Fetch ALL existing images for this asset
	images, err := s.imageRepo.FindByAssetID(ctx, assetID)
	if err != nil {
		return nil, response.Internal("failed to fetch images")
	}

	activeCount := 0
	newOrder := 1

	// 2. Calculate Count and Order simultaneously
	for _, img := range images {
		// Prevent DynamoDB SK collisions by always finding the absolute highest order
		if img.Order >= newOrder {
			newOrder = img.Order + 1
		}

		// Only count images that are successful or currently in progress
		if img.Status == enums.StatusImageCompleted || img.Status == enums.StatusImageUploading {
			activeCount++
		}
	}

	// 3. Enforce the limit based ONLY on active images
	if activeCount >= core.MaxImagesPerAsset {
		return nil, response.BadRequest(fmt.Sprintf("maximum %d images per asset reached", core.MaxImagesPerAsset))
	}

	// 4. Extract extension safely
	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".jpg"
	}

	// fileKey will now safely be assets/123/6.jpg if images 1-5 exist but 1 failed
	fileKey := fmt.Sprintf("assets/%s/%d%s", assetID, newOrder, ext)
	contentType := "image/jpeg"
	switch strings.ToLower(ext) {
	case ".png":
		contentType = "image/png"
	case ".webp":
		contentType = "image/webp"
	}

	uploadURL, err := s.storage.GeneratePresignedUploadURL(ctx, fileKey, contentType, 5*time.Minute)
	if err != nil {
		return nil, response.Internal("failed to generate upload URL")
	}

	imageEntity := models.ImageEntity{
		BaseItem: models.BaseItem{
			PK: "ASSET#" + assetID,
			SK: fmt.Sprintf("IMAGE#%d", newOrder),
		},
		FileKey: fileKey,
		Status:  enums.StatusImageUploading,
		Order:   newOrder,
	}

	if err := s.imageRepo.Create(ctx, []models.ImageEntity{imageEntity}); err != nil {
		return nil, response.Internal("failed to create image record")
	}

	// IMPORTANT: update asset status back to UPLOADING
	_ = s.assetRepo.UpdateAssetStatus(ctx, assetID, string(enums.StatusAssetUploading))

	return &res.UploadInfo{
		UploadURL: uploadURL,
		FileKey:   fileKey,
		Order:     newOrder,
		ExpiresIn: 300,
	}, nil
}

func (s *assetService) DeleteAssetImage(ctx context.Context, assetID string, imageID string) *response.AppError {
	asset, err := s.assetRepo.FindByID(ctx, assetID)
	if err != nil {
		return response.Internal("failed to find asset")
	}
	if asset == nil {
		return response.NotFound("asset not found")
	}

	images, err := s.imageRepo.FindByAssetID(ctx, assetID)
	if err != nil {
		return response.Internal("failed to find images")
	}

	var targetOrder int
	found := false
	for _, img := range images {
		imgID := fmt.Sprintf("img-%s-%d", assetID, img.Order)
		if imgID == imageID {
			targetOrder = img.Order
			found = true
			break
		}
	}

	if !found {
		return response.NotFound("image not found")
	}

	img, err := s.imageRepo.FindByAssetIDAndOrder(ctx, assetID, targetOrder)
	if err != nil {
		return response.Internal("failed to get image")
	}
	if img == nil {
		return response.NotFound("image not found")
	}

	if err := s.storage.DeleteObject(ctx, img.FileKey); err != nil {
		return response.Internal("failed to delete file from storage")
	}

	if err := s.imageRepo.Delete(ctx, assetID, targetOrder); err != nil {
		return response.Internal("failed to delete image record")
	}

	return nil
}
