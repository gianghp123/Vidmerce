package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gianghp123/Vidmerce/backend/internal/configs"
	"github.com/gianghp123/Vidmerce/backend/internal/core"
	"github.com/gianghp123/Vidmerce/backend/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/assets/dtos/req"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/assets/dtos/res"
	imageRepo "github.com/gianghp123/Vidmerce/backend/internal/modules/assets/repositories"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/auth/guards"
	"github.com/gianghp123/Vidmerce/backend/internal/storage"
	"github.com/gianghp123/Vidmerce/backend/internal/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AssetService interface {
	CreateAsset(ctx context.Context, req req.CreateAssetReq) (*res.PresignedUrlsRes, *response.AppError)
	ConfirmUpload(ctx context.Context, assetID string, req req.ConfirmUploadReq) (*res.AssetRes, *response.AppError)
	GetAsset(ctx context.Context, assetID string) (*res.AssetRes, *response.AppError)
	ListAssets(ctx context.Context, limit int, cursor string) (*response.PaginatedResult[res.AssetRes], *response.AppError)
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

	if appErr := guards.GuardOwn(ctx, utils.GetOwnerIDFromPK(asset.Pk)); appErr != nil {
		return nil, appErr
	}

	images, err := s.imageRepo.FindByAssetID(ctx, assetID)
	if err != nil {
		log.Error("Failed to find images", zap.String("assetId", assetID), zap.Error(err))
		return nil, response.Internal("failed to find images")
	}

	imageInfos := make([]res.ImageInfo, 0)
	for _, img := range images {
		imageInfos = append(imageInfos, res.ImageInfo{
			ID:       img.Sk,
			ImageURL: utils.GetCDNURL(img.FileKey),
		})
	}

	return &res.AssetRes{
		ID:         assetID,
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

	auth, err := guards.FromAuthContext(ctx)
	if err != nil {
		return nil, response.Unauthorized("authentication required")
	}

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
		assetID := item.Sk

		images, err := s.imageRepo.FindByAssetID(ctx, assetID)
		if err != nil {
			log.Warn("Failed to fetch images", zap.String("assetId", assetID), zap.Error(err))
			continue
		}

		var thumbnail res.ImageInfo
		for _, img := range images {
			thumbnail = res.ImageInfo{ImageURL: utils.GetCDNURL(img.FileKey)}
			break
		}

	assets = append(assets, res.AssetRes{
		ID:         assetID,
		Name:       item.Name,
		Price:      item.Price,
		Images:     []res.ImageInfo{thumbnail},
		ProductURL: item.ProductURL,
		Status:     string(item.Status),
		CreatedAt:  item.CreatedAt,
	})
	}

	log.Debug("Assets listed", zap.String("userId", auth.UserID), zap.Int("count", len(assets)), zap.Bool("hasMore", result.Meta.HasMore))
	return &response.PaginatedResult[res.AssetRes]{
		Data: assets,
		Meta: result.Meta,
	}, nil
}

func (s *assetService) CreateAsset(ctx context.Context, req req.CreateAssetReq) (*res.PresignedUrlsRes, *response.AppError) {
	log := configs.GetLogger()

	auth, err := guards.FromAuthContext(ctx)
	if err != nil {
		log.Error("Unauthorized", zap.Error(err))
		return nil, response.Unauthorized("authentication required")
	}

	imageCount := req.ImageCount
	if imageCount <= 0 {
		imageCount = 1
	}
	if imageCount > core.MaxImagesPerAsset {
		return nil, response.BadRequest(fmt.Sprintf("maximum %d images per asset", core.MaxImagesPerAsset))
	}

	assetID := uuid.New().String()
	uploads := make([]res.UploadInfo, 0, imageCount)

	for i := 1; i <= imageCount; i++ {
		fileKey := fmt.Sprintf("user/%s/asset/%s/%d.jpg", auth.UserID, assetID, i)
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
	}

	log.Debug("Asset created", zap.String("assetId", assetID), zap.Int("imageCount", imageCount))
	return &res.PresignedUrlsRes{
		AssetID: assetID,
		Uploads: uploads,
	}, nil
}

func (s *assetService) ConfirmUpload(ctx context.Context, assetID string, req req.ConfirmUploadReq) (*res.AssetRes, *response.AppError) {
	log := configs.GetLogger()
	auth, err := guards.FromAuthContext(ctx)
	if err != nil {
		log.Error("Unauthorized", zap.Error(err))
		return nil, response.Unauthorized("authentication required")
	}

	for _, key := range req.FileKeys {
		exists, err := s.storage.ObjectExists(ctx, key)
		if err != nil {
			log.Error("S3 verification failed", zap.String("key", key), zap.Error(err))
			return nil, response.Internal("storage verification failed")
		}
		if !exists {
			return nil, response.BadRequest(fmt.Sprintf("file not uploaded: %s", key))
		}
	}

	asset := models.AssetEntity{
		BaseItem:   utils.BuildUserAssetBaseItem(auth.UserID, assetID),
		Name:       req.Name,
		Price:      req.Price,
		ProductURL: req.ProductURL,
		Status:     enums.AssetStatusCompleted,
		ImageCount: len(req.FileKeys),
		CreatedAt:  utils.Now(),
	}
	
	// Normalize Sk to remove ASSET# prefix for DTO mapping
	asset.Sk = strings.TrimPrefix(asset.Sk, string(core.SkPrefixAsset)+core.KeySeparator)

	images := make([]models.ImageEntity, len(req.FileKeys))
	for i, fileKey := range req.FileKeys {
		images[i] = models.ImageEntity{
			BaseItem: models.BaseItem{
				Pk: utils.BuildPk(core.EntityTypeAsset, assetID),
				Sk: fmt.Sprintf("%s#%d", core.EntityTypeImage, i+1),
			},
			FileKey: fileKey,
		}
	}

	if err := s.assetRepo.CreateWithImages(ctx, asset, images); err != nil {
		log.Error("Failed to create asset/images", zap.Error(err))
		return nil, response.Internal("failed to create asset")
	}

	log.Info("Asset confirmed", zap.String("assetId", assetID), zap.Int("imageCount", len(images)))

	imageInfos := make([]res.ImageInfo, len(images))
	for i, img := range images {
		imageInfos[i] = res.ImageInfo{
			ID:       fmt.Sprintf("%d", i+1),
			ImageURL: utils.GetCDNURL(img.FileKey),
		}
	}

	var result res.AssetRes

	err = utils.MapToDTO(asset, &result)
	if err != nil {
		log.Error("Failed to map asset to DTO", zap.Error(err))
		return nil, response.Internal("failed to map asset to DTO")
	}

	result.Images = imageInfos
	return &result, nil
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

	if appErr := guards.GuardOwn(ctx, utils.GetOwnerIDFromPK(asset.Pk)); appErr != nil {
		return appErr
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

	auth, err := guards.FromAuthContext(ctx)
	if err != nil {
		return nil, response.Unauthorized("authentication required")
	}

	assetID := uuid.New().String()
	jobID := uuid.New().String()

	log.Debug("Importing asset", zap.String("assetId", assetID), zap.String("jobId", jobID), zap.String("productUrl", req.ProductURL))

	asset := models.AssetEntity{
		BaseItem:   utils.BuildUserAssetBaseItem(auth.UserID, assetID),
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

	if err := s.assetRepo.CreateWithJob(ctx, asset, job); err != nil {
		log.Error("Failed to create import job", zap.String("assetId", assetID), zap.String("jobId", jobID), zap.Error(err))
		return nil, response.Internal("failed to create import job: " + err.Error())
	}

	log.Info("Asset import started", zap.String("assetId", assetID), zap.String("jobId", jobID), zap.String("userId", auth.UserID))
	return &res.ImportAssetRes{
		AssetID: assetID,
		JobID:   jobID,
	}, nil
}
