package services

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gianghp123/Vidmerce/backend/internal/configs"
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
	"go.uber.org/zap"
)

type AssetService interface {
	CreateAsset(ctx context.Context, req req.CreateAssetReq) (*res.CreateAssetRes, *response.AppError)
	ConfirmUpload(ctx context.Context, assetID string) (*res.ConfirmAssetRes, *response.AppError)
	GetAsset(ctx context.Context, assetID string) (*res.AssetRes, *response.AppError)
	ListAssets(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.AssetRes], *response.AppError)
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
	log := configs.GetLogger()

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
		Status:     enums.AssetStatusUploading,
		ImageCount: imageCount,
		CreatedAt:  utils.Now(),
	}

	uploads := make([]res.UploadInfo, 0, imageCount)
	images := make([]models.ImageEntity, 0, imageCount)

	for i := 1; i <= imageCount; i++ {
		fileKey := fmt.Sprintf("assets/%s/%d.jpg", assetID, i)

		uploadURL, err := s.storage.GeneratePresignedUploadURL(ctx, fileKey, "image/jpeg", 5*time.Minute)
		if err != nil {
			log.Error("Failed to generate upload URL", zap.String("assetId", assetID), zap.Error(err))
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
				Pk: utils.BuildPk(core.EntityTypeAsset, assetID),
				Sk: fmt.Sprintf("%s#%d", core.EntityTypeImage, i),
			},
			FileKey: fileKey,
			Status:  enums.ImageStatusUploading,
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
		log.Error("Failed to create asset", zap.String("assetId", assetID), zap.Error(err))
		return nil, response.Internal("failed to create asset and images: " + err.Error())
	}

	log.Debug("Asset created", zap.String("assetId", assetID), zap.Int("imageCount", imageCount))
	return &res.CreateAssetRes{
		AssetID: assetID,
		Status:  string(enums.AssetStatusUploading),
		Uploads: uploads,
	}, nil
}

func (s *assetService) ConfirmUpload(ctx context.Context, assetID string) (*res.ConfirmAssetRes, *response.AppError) {
	log := configs.GetLogger()

	asset, err := s.assetRepo.FindByID(ctx, assetID)
	if err != nil {
		log.Error("Failed to find asset for confirmation", zap.String("assetId", assetID), zap.Error(err))
		return nil, response.Internal("failed to find asset")
	}
	if asset == nil {
		return nil, response.NotFound("asset not found")
	}

	dbImages, err := s.imageRepo.FindByAssetID(ctx, assetID)
	if err != nil {
		log.Error("Failed to get images", zap.String("assetId", assetID), zap.Error(err))
		return nil, response.Internal("failed to get images")
	}

	images := make([]res.ImageInfo, 0, len(dbImages))
	uploadedCount := 0

	for _, img := range dbImages {
		fileKey := img.FileKey

		exists, err := s.storage.ObjectExists(ctx, fileKey)
		if err != nil {
			log.Error("Failed to verify upload", zap.String("assetId", assetID), zap.String("fileKey", fileKey), zap.Error(err))
			return nil, response.Internal("failed to verify upload")
		}

		imageUrl := utils.GetCDNURL(fileKey)

		var status string

		if exists {
			status = string(enums.ImageStatusCompleted)
			uploadedCount++
		} else {
			status = string(enums.ImageStatusFailed)
		}

		if err := s.imageRepo.UpdateStatus(ctx, assetID, img.Order, status); err != nil {
			log.Error("Failed to update image status", zap.String("assetId", assetID), zap.Int("order", img.Order), zap.Error(err))
			return nil, response.Internal("failed to update image status")
		}

		images = append(images, res.ImageInfo{
			ImageURL: imageUrl,
			Order:    img.Order,
			Status:   status,
		})
	}

	total := len(dbImages)

	var finalStatus string

	switch {
	case total == 0:
		finalStatus = string(enums.AssetStatusFailed)
	case uploadedCount == total:
		finalStatus = string(enums.AssetStatusCompleted)
	case uploadedCount > 0:
		finalStatus = string(enums.AssetStatusPartial)
	default:
		finalStatus = string(enums.AssetStatusFailed)
	}

	if err := s.assetRepo.UpdateAssetStatus(ctx, assetID, finalStatus); err != nil {
		log.Error("Failed to update asset status", zap.String("assetId", assetID), zap.String("status", finalStatus), zap.Error(err))
		return nil, response.Internal("failed to update asset status")
	}

	log.Debug("Asset upload confirmed", zap.String("assetId", assetID), zap.String("status", finalStatus), zap.Int("uploadedCount", uploadedCount), zap.Int("total", total))
	return &res.ConfirmAssetRes{
		AssetID: assetID,
		Status:  finalStatus,
		Images:  images,
	}, nil
}

func (s *assetService) GetAsset(ctx context.Context, assetID string) (*res.AssetRes, *response.AppError) {
	log := configs.GetLogger()

	asset, err := s.assetRepo.FindByID(ctx, assetID)
	if err != nil {
		log.Error("Failed to find asset", zap.String("assetId", assetID), zap.Error(err))
		return nil, response.Internal("failed to find asset")
	}
	if asset == nil {
		return nil, response.NotFound("asset not found")
	}

	images, err := s.imageRepo.FindByAssetID(ctx, assetID)
	if err != nil {
		log.Error("Failed to find images", zap.String("assetId", assetID), zap.Error(err))
		return nil, response.Internal("failed to find images")
	}

	imageInfos := make([]res.ImageInfo, 0)
	for _, img := range images {
		if img.Status == enums.ImageStatusCompleted {
			imageURL := utils.GetCDNURL(img.FileKey)
			imageInfos = append(imageInfos, res.ImageInfo{
				ImageID:  img.Sk,
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

func (s *assetService) ListAssets(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.AssetRes], *response.AppError) {
	log := configs.GetLogger()

	if limit <= 0 {
		limit = 20
	}

	result, err := s.assetRepo.FindAll(ctx, limit, cursor)
	if err != nil {
		log.Error("Failed to fetch assets", zap.Error(err))
		return nil, response.Internal("failed to fetch assets: " + err.Error())
	}

	assets := make([]res.AssetRes, 0, len(result.Data))
	for _, item := range result.Data {
		assetID := item.Pk

		img, err := s.imageRepo.FindOneCompletedImageByAssetId(ctx, assetID)
		if err != nil {
			log.Warn("Failed to fetch thumbnail", zap.String("assetId", assetID), zap.Error(err))
			continue
		}

		var image res.ImageInfo
		if img != nil {
			image = res.ImageInfo{
				ImageURL: utils.GetCDNURL(img.FileKey),
				Order:    img.Order,
			}
		}

		assets = append(assets, res.AssetRes{
			AssetID:    assetID,
			Name:       item.Name,
			Price:      item.Price,
			Images:     []res.ImageInfo{image},
			ProductURL: item.ProductURL,
			Status:     string(item.Status),
			CreatedAt:  item.CreatedAt,
		})
	}

	log.Debug("Assets listed", zap.Int("count", len(assets)), zap.Bool("hasMore", result.Meta.HasMore))
	return &response.PaginatedResult[res.AssetRes]{
		Data: assets,
		Meta: result.Meta,
	}, nil
}

func (s *assetService) GetImageUploadUrl(ctx context.Context, assetID string, fileName string) (*res.UploadInfo, *response.AppError) {
	log := configs.GetLogger()

	asset, err := s.assetRepo.FindByID(ctx, assetID)
	if err != nil {
		log.Error("Failed to find asset", zap.String("assetId", assetID), zap.Error(err))
		return nil, response.Internal("failed to find asset")
	}
	if asset == nil {
		return nil, response.NotFound("asset not found")
	}

	images, err := s.imageRepo.FindByAssetID(ctx, assetID)
	if err != nil {
		log.Error("Failed to fetch images", zap.String("assetId", assetID), zap.Error(err))
		return nil, response.Internal("failed to fetch images")
	}

	activeCount := 0
	newOrder := 1

	for _, img := range images {
		if img.Order >= newOrder {
			newOrder = img.Order + 1
		}

		if img.Status == enums.ImageStatusCompleted || img.Status == enums.ImageStatusUploading {
			activeCount++
		}
	}

	if activeCount >= core.MaxImagesPerAsset {
		return nil, response.BadRequest(fmt.Sprintf("maximum %d images per asset reached", core.MaxImagesPerAsset))
	}
	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".jpg"
	}

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
		log.Error("Failed to generate upload URL", zap.String("assetId", assetID), zap.Error(err))
		return nil, response.Internal("failed to generate upload URL")
	}

	imageEntity := models.ImageEntity{
		BaseItem: models.BaseItem{
			Pk: utils.BuildPk(core.EntityTypeAsset, assetID),
			Sk: fmt.Sprintf("%s#%d", core.EntityTypeImage, newOrder),
		},
		FileKey: fileKey,
		Status:  enums.ImageStatusUploading,
		Order:   newOrder,
	}

	if err := s.imageRepo.Create(ctx, []models.ImageEntity{imageEntity}); err != nil {
		log.Error("Failed to create image record", zap.String("assetId", assetID), zap.Error(err))
		return nil, response.Internal("failed to create image record")
	}

	_ = s.assetRepo.UpdateAssetStatus(ctx, assetID, string(enums.AssetStatusUploading))

	log.Debug("Generated upload URL", zap.String("assetId", assetID), zap.Int("order", newOrder))
	return &res.UploadInfo{
		UploadURL: uploadURL,
		FileKey:   fileKey,
		Order:     newOrder,
		ExpiresIn: 300,
	}, nil
}

func (s *assetService) DeleteAssetImage(ctx context.Context, assetID string, imageID string) *response.AppError {
	log := configs.GetLogger()

	asset, err := s.assetRepo.FindByID(ctx, assetID)
	if err != nil {
		log.Error("Failed to find asset", zap.String("assetId", assetID), zap.Error(err))
		return response.Internal("failed to find asset")
	}
	if asset == nil {
		return response.NotFound("asset not found")
	}

	images, err := s.imageRepo.FindByAssetID(ctx, assetID)
	if err != nil {
		log.Error("Failed to find images", zap.String("assetId", assetID), zap.Error(err))
		return response.Internal("failed to find images")
	}

	var targetImg *models.ImageEntity
	for _, img := range images {
		if img.Sk == imageID {
			targetImg = &img
			break
		}
	}

	if targetImg == nil {
		return response.NotFound("image not found")
	}

	if err := s.storage.DeleteObject(ctx, targetImg.FileKey); err != nil {
		log.Error("Failed to delete file from storage", zap.String("assetId", assetID), zap.String("imageId", imageID), zap.Error(err))
		return response.Internal("failed to delete file from storage")
	}

	if err := s.imageRepo.Delete(ctx, assetID, targetImg.Sk); err != nil {
		log.Error("Failed to delete image record", zap.String("assetId", assetID), zap.String("imageId", imageID), zap.Error(err))
		return response.Internal("failed to delete image record")
	}

	log.Debug("Image deleted", zap.String("assetId", assetID), zap.String("imageId", imageID))
	return nil
}

func (s *assetService) ImportAsset(ctx context.Context, req req.ImportAssetReq) (*res.ImportAssetRes, *response.AppError) {
	log := configs.GetLogger()

	assetID := uuid.New().String()
	jobID := uuid.New().String()

	log.Debug("Importing asset", zap.String("assetId", assetID), zap.String("jobId", jobID), zap.String("productUrl", req.ProductURL))

	asset := models.AssetEntity{
		BaseItem:   utils.BuildAssetBaseItem(assetID),
		ProductURL: req.ProductURL,
		Status:     enums.AssetStatusImporting,
		CreatedAt:  utils.Now(),
	}

	job := models.JobEntity{
		BaseItem:  utils.BuildJobBaseItem(jobID),
		TargetID:  &assetID,
		Type:      enums.JobTypeScrapeProduct,
		Status:    enums.JobStatusPending,
		Payload:   map[string]any{"url": req.ProductURL},
		CreatedAt: utils.Now(),
	}

	baseRepo := repositories.NewBaseRepository(s.assetRepo.DBClient())
	if err := baseRepo.TransactWriteItems(ctx, asset, job); err != nil {
		log.Error("Failed to create import job", zap.String("assetId", assetID), zap.String("jobId", jobID), zap.Error(err))
		return nil, response.Internal("failed to create import job: " + err.Error())
	}

	log.Info("Asset import started", zap.String("assetId", assetID), zap.String("jobId", jobID))
	return &res.ImportAssetRes{
		AssetID: assetID,
		JobID:   jobID,
	}, nil
}
