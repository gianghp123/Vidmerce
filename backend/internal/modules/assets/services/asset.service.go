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
	"github.com/gianghp123/Vidmerce/backend/internal/database/repositories"
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
	ListAssets(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.AssetPreviewRes], *response.AppError)
	GetImageUploadUrl(ctx context.Context, assetID string, fileName string) (*res.UploadInfo, *response.AppError)
	DeleteAssetImage(ctx context.Context, assetID string, imageID string) *response.AppError
	ImportAsset(ctx context.Context, req req.ImportAssetReq) (*res.ImportAssetRes, *response.AppError)
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

	asset := models.AssetEntity{
		BaseItem:   utils.BuildAssetBaseItem(assetID),
		Name:       req.Name,
		Price:      req.Price,
		ProductURL: req.ProductURL,
		Status:     enums.StatusAssetUploading,
		ImageCount: imageCount,
		CreatedAt:  utils.Now(),
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
				PK: utils.BuildPK(enums.EntityTypeAsset, assetID),
				SK: fmt.Sprintf("%s#%d", enums.EntityTypeImage, i),
			},
			FileKey: fileKey,
			Status:  enums.StatusImageUploading,
			Order:   i,
		})
	}

	baseRepo := repositories.NewBaseRepository(s.assetRepo.DBClient())

	items := make([]interface{}, 0, len(images)+1)
	items = append(items, asset)
	for _, img := range images {
		items = append(items, img)
	}

	if err := baseRepo.TransactWriteItems(ctx, items...); err != nil {
		return nil, response.Internal("failed to create asset and images: " + err.Error())
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

	imageInfos := make([]res.ImageInfo, 0)
	for _, img := range images {
		if img.Status == enums.StatusImageCompleted {
			imageURL := utils.GetCDNURL(img.FileKey)
			imageInfos = append(imageInfos, res.ImageInfo{
				ImageID:  img.SK,
				ImageURL: imageURL,
				Order:    img.Order,
				Status:   string(img.Status),
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

func (s *assetService) ListAssets(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.AssetPreviewRes], *response.AppError) {
	if limit <= 0 {
		limit = 20
	}

	result, err := s.assetRepo.FindAll(ctx, limit, cursor)
	if err != nil {
		return nil, response.Internal("failed to fetch assets: " + err.Error())
	}

	assets := make([]res.AssetPreviewRes, 0, len(result.Data))
	for _, item := range result.Data {
		assetID := item.PK

		img, err := s.imageRepo.FindOneCompletedImageByAssetId(ctx, assetID)
		if err != nil {
			continue
		}

		var image res.ImageInfo
		if img != nil {
			image = res.ImageInfo{
				ImageURL: utils.GetCDNURL(img.FileKey),
				Order:    img.Order,
			}
		}

		assets = append(assets, res.AssetPreviewRes{
			AssetID:    assetID,
			Name:       item.Name,
			Price:      item.Price,
			Image:      image,
			ProductURL: item.ProductURL,
			Status:     string(item.Status),
			CreatedAt:  item.CreatedAt,
		})
	}

	return &response.PaginatedResult[res.AssetPreviewRes]{
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
			PK: utils.BuildPK(enums.EntityTypeAsset, assetID),
			SK: fmt.Sprintf("%s#%d", enums.EntityTypeImage, newOrder),
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

	var targetImg *models.ImageEntity // adjust type if needed
	for _, img := range images {
		if img.SK == imageID {
			targetImg = &img
			break
		}
	}

	if targetImg == nil {
		return response.NotFound("image not found")
	}

	// delete file from storage
	if err := s.storage.DeleteObject(ctx, targetImg.FileKey); err != nil {
		return response.Internal("failed to delete file from storage")
	}

	// delete DB record
	if err := s.imageRepo.Delete(ctx, assetID, targetImg.SK); err != nil {
		return response.Internal("failed to delete image record")
	}

	return nil
}

func (s *assetService) ImportAsset(ctx context.Context, req req.ImportAssetReq) (*res.ImportAssetRes, *response.AppError) {
	assetID := uuid.New().String()
	jobID := uuid.New().String()

	asset := models.AssetEntity{
		BaseItem:   utils.BuildAssetBaseItem(assetID),
		ProductURL: req.ProductURL,
		Status:     enums.StatusAssetImporting,
		CreatedAt:  utils.Now(),
	}

	job := models.JobEntity{
		BaseItem:  utils.BuildJobBaseItem(jobID),
		TargetID:  assetID,
		Type:      enums.TypeJobScrapeProduct,
		Status:    enums.StatusJobPending,
		Payload:   map[string]any{"url": req.ProductURL},
		CreatedAt: utils.Now(),
	}

	baseRepo := repositories.NewBaseRepository(s.assetRepo.DBClient())
	if err := baseRepo.TransactWriteItems(ctx, asset, job); err != nil {
		return nil, response.Internal("failed to create import job: " + err.Error())
	}

	return &res.ImportAssetRes{
		AssetID: assetID,
		JobID:   jobID,
	}, nil
}
