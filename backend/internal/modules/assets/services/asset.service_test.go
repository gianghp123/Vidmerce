package services

import (
	"context"
	"errors"
	"testing"

	"github.com/gianghp123/Vidmerce/backend/services/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/services/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/services/internal/database/models"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/assets/dtos/req"
	repoMocks "github.com/gianghp123/Vidmerce/backend/services/internal/modules/assets/repositories"
	"github.com/gianghp123/Vidmerce/backend/services/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAssetService(t *testing.T) (*assetService, *repoMocks.MockAssetRepository, *repoMocks.MockImageRepository, *storage.MockStorage) {
	t.Helper()
	assetRepo := repoMocks.NewMockAssetRepository()
	imageRepo := repoMocks.NewMockImageRepository()
	storageMock := storage.NewMockStorage()
	svc := NewAssetService(assetRepo, imageRepo, storageMock).(*assetService)
	return svc, assetRepo, imageRepo, storageMock
}

// ==================== CreateAsset Tests ====================

func TestCreateAsset(t *testing.T) {
	tests := []struct {
		name          string
		req           req.CreateAssetReq
		setupMock     func(*repoMocks.MockAssetRepository, *repoMocks.MockImageRepository, *storage.MockStorage)
		wantErr       bool
		wantCode      int
		wantStatus    string
		wantUploadLen int
	}{
		{
			name: "success",
			req: req.CreateAssetReq{
				Name:       "Test Asset",
				Price:      99.99,
				ProductURL: "https://example.com/product",
				ImageCount: 1,
			},
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				assetRepo.On("Create", mock.Anything, mock.AnythingOfType("models.AssetEntity")).Return(nil)
				imageRepo.On("Create", mock.Anything, mock.AnythingOfType("[]models.ImageEntity")).Return(nil)
				storageMock.On("GeneratePresignedUploadURL", mock.Anything, mock.Anything, "image/jpeg", mock.Anything).Return("https://upload.url/1", nil)
			},
			wantErr:       false,
			wantStatus:    "UPLOADING",
			wantUploadLen: 1,
		},
		{
			name: "success multiple images",
			req: req.CreateAssetReq{
				Name:       "Test Asset",
				Price:      99.99,
				ProductURL: "https://example.com/product",
				ImageCount: 3,
			},
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				assetRepo.On("Create", mock.Anything, mock.AnythingOfType("models.AssetEntity")).Return(nil)
				imageRepo.On("Create", mock.Anything, mock.AnythingOfType("[]models.ImageEntity")).Return(nil)
				storageMock.On("GeneratePresignedUploadURL", mock.Anything, mock.Anything, "image/jpeg", mock.Anything).Return("https://upload.url", nil)
			},
			wantErr:       false,
			wantUploadLen: 3,
		},
		{
			name: "default image count",
			req: req.CreateAssetReq{
				Name:       "Test Asset",
				Price:      99.99,
				ProductURL: "https://example.com/product",
				ImageCount: 0,
			},
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				assetRepo.On("Create", mock.Anything, mock.AnythingOfType("models.AssetEntity")).Return(nil)
				imageRepo.On("Create", mock.Anything, mock.AnythingOfType("[]models.ImageEntity")).Return(nil)
				storageMock.On("GeneratePresignedUploadURL", mock.Anything, mock.Anything, "image/jpeg", mock.Anything).Return("https://upload.url", nil)
			},
			wantErr:       false,
			wantUploadLen: 1,
		},
		{
			name: "exceeds max images",
			req: req.CreateAssetReq{
				Name:       "Test Asset",
				Price:      99.99,
				ProductURL: "https://example.com/product",
				ImageCount: 6,
			},
			setupMock: nil,
			wantErr:   true,
		},
		{
			name: "storage error",
			req: req.CreateAssetReq{
				Name:       "Test Asset",
				Price:      99.99,
				ProductURL: "https://example.com/product",
				ImageCount: 1,
			},
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				storageMock.On("GeneratePresignedUploadURL", mock.Anything, mock.Anything, "image/jpeg", mock.Anything).Return("", errors.New("storage error"))
			},
			wantErr:  true,
			wantCode: 500,
		},
		{
			name: "image repo error",
			req: req.CreateAssetReq{
				Name:       "Test Asset",
				Price:      99.99,
				ProductURL: "https://example.com/product",
				ImageCount: 1,
			},
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				storageMock.On("GeneratePresignedUploadURL", mock.Anything, mock.Anything, "image/jpeg", mock.Anything).Return("https://upload.url", nil)
				imageRepo.On("Create", mock.Anything, mock.AnythingOfType("[]models.ImageEntity")).Return(errors.New("db error"))
			},
			wantErr:  true,
			wantCode: 500,
		},
		{
			name: "asset repo error",
			req: req.CreateAssetReq{
				Name:       "Test Asset",
				Price:      99.99,
				ProductURL: "https://example.com/product",
				ImageCount: 1,
			},
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				storageMock.On("GeneratePresignedUploadURL", mock.Anything, mock.Anything, "image/jpeg", mock.Anything).Return("https://upload.url", nil)
				imageRepo.On("Create", mock.Anything, mock.AnythingOfType("[]models.ImageEntity")).Return(nil)
				assetRepo.On("Create", mock.Anything, mock.AnythingOfType("models.AssetEntity")).Return(errors.New("db error"))
			},
			wantErr:  true,
			wantCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, assetRepo, imageRepo, storageMock := setupAssetService(t)
			if tt.setupMock != nil {
				tt.setupMock(assetRepo, imageRepo, storageMock)
			}

			result, appErr := svc.CreateAsset(context.Background(), tt.req)

			if tt.wantErr {
				assert.Nil(t, result)
				assert.NotNil(t, appErr)
				if tt.wantCode > 0 {
					assert.Equal(t, tt.wantCode, appErr.Code)
				}
			} else {
				assert.Nil(t, appErr)
				assert.NotNil(t, result)
				if tt.wantStatus != "" {
					assert.Equal(t, tt.wantStatus, result.Status)
				}
				if tt.wantUploadLen > 0 {
					assert.Equal(t, tt.wantUploadLen, len(result.Uploads))
				}
			}
		})
	}
}

// ==================== ConfirmUpload Tests ====================

func TestConfirmUpload(t *testing.T) {
	uploadingAsset := &models.AssetEntity{
		BaseItem: models.BaseItem{
			PK: "ASSET#test-asset-id",
			SK: "METADATA",
		},
		Name:       "Test Asset",
		Price:      99.99,
		ProductURL: "https://example.com",
		Status:     enums.StatusAssetUploading,
		ImageCount: 1,
		CreatedAt:  "2024-01-01T00:00:00Z",
	}

	tests := []struct {
		name      string
		assetID   string
		setupMock func(*repoMocks.MockAssetRepository, *repoMocks.MockImageRepository, *storage.MockStorage)
		wantErr   bool
		wantCode  int
	}{
		{
			name:    "success",
			assetID: "test-asset-id",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				assetRepo.On("FindByID", mock.Anything, "test-asset-id").Return(uploadingAsset, nil)
				storageMock.On("ObjectExists", mock.Anything, "assets/test-asset-id/1.jpg").Return(true, nil)
				imageRepo.On("UpdateStatus", mock.Anything, "test-asset-id", 1, "COMPLETED").Return(nil)
				assetRepo.On("UpdateAssetStatus", mock.Anything, "test-asset-id", "COMPLETED").Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "asset not found",
			assetID: "non-existent",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, _ *repoMocks.MockImageRepository, _ *storage.MockStorage) {
				assetRepo.On("FindByID", mock.Anything, "non-existent").Return(nil, nil)
			},
			wantErr:  true,
			wantCode: 404,
		},
		{
			name:    "asset repo error",
			assetID: "test-id",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, _ *repoMocks.MockImageRepository, _ *storage.MockStorage) {
				assetRepo.On("FindByID", mock.Anything, "test-id").Return(nil, errors.New("db error"))
			},
			wantErr:  true,
			wantCode: 500,
		},
		{
			name:    "wrong status",
			assetID: "test-id",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, _ *repoMocks.MockImageRepository, _ *storage.MockStorage) {
				asset := &models.AssetEntity{
					BaseItem: models.BaseItem{
						PK: "ASSET#test-id",
						SK: "METADATA",
					},
					Status:     enums.StatusAssetCompleted,
					ImageCount: 1,
				}
				assetRepo.On("FindByID", mock.Anything, "test-id").Return(asset, nil)
			},
			wantErr:  true,
			wantCode: 400,
		},
		{
			name:    "image not uploaded",
			assetID: "test-id",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, _ *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				assetRepo.On("FindByID", mock.Anything, "test-id").Return(uploadingAsset, nil)
				storageMock.On("ObjectExists", mock.Anything, "assets/test-id/1.jpg").Return(false, nil)
			},
			wantErr:  true,
			wantCode: 400,
		},
		{
			name:    "storage error",
			assetID: "test-id",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, _ *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				assetRepo.On("FindByID", mock.Anything, "test-id").Return(uploadingAsset, nil)
				storageMock.On("ObjectExists", mock.Anything, "assets/test-id/1.jpg").Return(false, errors.New("storage error"))
			},
			wantErr:  true,
			wantCode: 500,
		},
		{
			name:    "update image status error",
			assetID: "test-id",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				assetRepo.On("FindByID", mock.Anything, "test-id").Return(uploadingAsset, nil)
				storageMock.On("ObjectExists", mock.Anything, "assets/test-id/1.jpg").Return(true, nil)
				imageRepo.On("UpdateStatus", mock.Anything, "test-id", 1, "COMPLETED").Return(errors.New("db error"))
			},
			wantErr:  true,
			wantCode: 500,
		},
		{
			name:    "update asset status error",
			assetID: "test-id",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				assetRepo.On("FindByID", mock.Anything, "test-id").Return(uploadingAsset, nil)
				storageMock.On("ObjectExists", mock.Anything, "assets/test-id/1.jpg").Return(true, nil)
				imageRepo.On("UpdateStatus", mock.Anything, "test-id", 1, "COMPLETED").Return(nil)
				assetRepo.On("UpdateAssetStatus", mock.Anything, "test-id", "COMPLETED").Return(errors.New("db error"))
			},
			wantErr:  true,
			wantCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, assetRepo, imageRepo, storageMock := setupAssetService(t)
			tt.setupMock(assetRepo, imageRepo, storageMock)

			result, appErr := svc.ConfirmUpload(context.Background(), tt.assetID)

			if tt.wantErr {
				assert.Nil(t, result)
				assert.NotNil(t, appErr)
				if tt.wantCode > 0 {
					assert.Equal(t, tt.wantCode, appErr.Code)
				}
			} else {
				assert.Nil(t, appErr)
				assert.NotNil(t, result)
				assert.Equal(t, tt.assetID, result.AssetID)
				assert.Equal(t, "COMPLETED", result.Status)
			}
		})
	}
}

// ==================== GetAsset Tests ====================

func TestGetAsset(t *testing.T) {
	asset := &models.AssetEntity{
		BaseItem: models.BaseItem{
			PK: "ASSET#test-id",
			SK: "METADATA",
		},
		Name:       "Test Asset",
		Price:      99.99,
		ProductURL: "https://example.com",
		Status:     enums.StatusAssetCompleted,
		ImageCount: 1,
		CreatedAt:  "2024-01-01T00:00:00Z",
	}

	tests := []struct {
		name      string
		assetID   string
		setupMock func(*repoMocks.MockAssetRepository, *repoMocks.MockImageRepository)
		wantErr   bool
		wantID    string
		wantLen   int
		wantCode  int
	}{
		{
			name:    "success",
			assetID: "test-id",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository) {
				images := []models.ImageEntity{
					{
						BaseItem: models.BaseItem{
							PK: "ASSET#test-id",
							SK: "IMAGE#1",
						},
						FileKey: "assets/test-id/1.jpg",
						Status:  enums.StatusImageCompleted,
						Order:   1,
					},
				}
				assetRepo.On("FindByID", mock.Anything, "test-id").Return(asset, nil)
				imageRepo.On("FindByAssetID", mock.Anything, "test-id").Return(images, nil)
			},
			wantErr: false,
			wantID:  "ASSET#test-id",
			wantLen: 1,
		},
		{
			name:    "asset not found",
			assetID: "non-existent",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, _ *repoMocks.MockImageRepository) {
				assetRepo.On("FindByID", mock.Anything, "non-existent").Return(nil, nil)
			},
			wantErr:  true,
			wantCode: 404,
		},
		{
			name:    "asset repo error",
			assetID: "test-id",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, _ *repoMocks.MockImageRepository) {
				assetRepo.On("FindByID", mock.Anything, "test-id").Return(nil, errors.New("db error"))
			},
			wantErr:  true,
			wantCode: 500,
		},
		{
			name:    "image repo error",
			assetID: "test-id",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository) {
				assetRepo.On("FindByID", mock.Anything, "test-id").Return(asset, nil)
				imageRepo.On("FindByAssetID", mock.Anything, "test-id").Return(nil, errors.New("db error"))
			},
			wantErr:  true,
			wantCode: 500,
		},
		{
			name:    "only completed images",
			assetID: "test-id",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository) {
				asset := &models.AssetEntity{
					BaseItem: models.BaseItem{
						PK: "ASSET#test-id",
						SK: "METADATA",
					},
					Name:       "Test Asset",
					Price:      99.99,
					ProductURL: "https://example.com",
					Status:     enums.StatusAssetCompleted,
					ImageCount: 2,
					CreatedAt:  "2024-01-01T00:00:00Z",
				}
				images := []models.ImageEntity{
					{
						BaseItem: models.BaseItem{
							PK: "ASSET#test-id",
							SK: "IMAGE#1",
						},
						FileKey: "assets/test-id/1.jpg",
						Status:  enums.StatusImageCompleted,
						Order:   1,
					},
					{
						BaseItem: models.BaseItem{
							PK: "ASSET#test-id",
							SK: "IMAGE#2",
						},
						FileKey: "assets/test-id/2.jpg",
						Status:  enums.StatusImageUploading,
						Order:   2,
					},
				}
				assetRepo.On("FindByID", mock.Anything, "test-id").Return(asset, nil)
				imageRepo.On("FindByAssetID", mock.Anything, "test-id").Return(images, nil)
			},
			wantErr: false,
			wantLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, assetRepo, imageRepo, _ := setupAssetService(t)
			tt.setupMock(assetRepo, imageRepo)

			result, appErr := svc.GetAsset(context.Background(), tt.assetID)

			if tt.wantErr {
				assert.Nil(t, result)
				assert.NotNil(t, appErr)
				if tt.wantCode > 0 {
					assert.Equal(t, tt.wantCode, appErr.Code)
				}
			} else {
				assert.Nil(t, appErr)
				assert.NotNil(t, result)
				if tt.wantID != "" {
					assert.Equal(t, tt.wantID, result.AssetID)
				}
				if tt.wantLen > 0 {
					assert.Equal(t, tt.wantLen, len(result.Images))
				}
			}
		})
	}
}

// ==================== ListAssets Tests ====================

func TestListAssets(t *testing.T) {
	tests := []struct {
		name      string
		limit     int
		cursor    string
		setupMock func(*repoMocks.MockAssetRepository, *repoMocks.MockImageRepository)
		wantErr   bool
		wantLen   int
		wantCode  int
	}{
		{
			name:   "success",
			limit:  10,
			cursor: "",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository) {
				assets := []models.AssetEntity{
					{
						BaseItem: models.BaseItem{
							PK:     "ASSET#asset-1",
							SK:     "METADATA",
							GSI1PK: "ENTITY#ASSET",
							GSI1SK: "asset-1",
						},
						Name:       "Asset 1",
						Price:      99.99,
						ProductURL: "https://example.com/1",
						Status:     enums.StatusAssetCompleted,
						ImageCount: 1,
						CreatedAt:  "2024-01-01T00:00:00Z",
					},
				}
				images := []models.ImageEntity{
					{
						BaseItem: models.BaseItem{
							PK: "ASSET#asset-1",
							SK: "IMAGE#1",
						},
						FileKey: "assets/asset-1/1.jpg",
						Status:  enums.StatusImageCompleted,
						Order:   1,
					},
				}
				paginatedResult := &response.PaginatedResult[models.AssetEntity]{
					Data: assets,
					Meta: response.Meta{Limit: 10, HasMore: false},
				}
				assetRepo.On("FindAll", mock.Anything, 10, "").Return(paginatedResult, nil)
				imageRepo.On("FindByAssetID", mock.Anything, "ASSET#asset-1").Return(images, nil)
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name:   "default limit",
			limit:  0,
			cursor: "",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, _ *repoMocks.MockImageRepository) {
				paginatedResult := &response.PaginatedResult[models.AssetEntity]{
					Data: []models.AssetEntity{},
					Meta: response.Meta{Limit: 20, HasMore: false},
				}
				assetRepo.On("FindAll", mock.Anything, 20, "").Return(paginatedResult, nil)
			},
			wantErr: false,
		},
		{
			name:   "repo error",
			limit:  10,
			cursor: "",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, _ *repoMocks.MockImageRepository) {
				assetRepo.On("FindAll", mock.Anything, 10, "").Return(nil, errors.New("db error"))
			},
			wantErr:  true,
			wantCode: 500,
		},
		{
			name:   "skips failed image fetch",
			limit:  10,
			cursor: "",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository) {
				assets := []models.AssetEntity{
					{
						BaseItem: models.BaseItem{
							PK:     "ASSET#asset-1",
							SK:     "METADATA",
							GSI1PK: "ENTITY#ASSET",
							GSI1SK: "asset-1",
						},
						Name:       "Asset 1",
						Price:      99.99,
						ProductURL: "https://example.com/1",
						Status:     enums.StatusAssetCompleted,
						ImageCount: 1,
						CreatedAt:  "2024-01-01T00:00:00Z",
					},
				}
				paginatedResult := &response.PaginatedResult[models.AssetEntity]{
					Data: assets,
					Meta: response.Meta{Limit: 10, HasMore: false},
				}
				assetRepo.On("FindAll", mock.Anything, 10, "").Return(paginatedResult, nil)
				imageRepo.On("FindByAssetID", mock.Anything, "ASSET#asset-1").Return(nil, errors.New("db error"))
			},
			wantErr: false,
			wantLen: 0,
		},
		{
			name:   "only first completed image",
			limit:  10,
			cursor: "",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository) {
				assets := []models.AssetEntity{
					{
						BaseItem: models.BaseItem{
							PK:     "ASSET#asset-1",
							SK:     "METADATA",
							GSI1PK: "ENTITY#ASSET",
							GSI1SK: "asset-1",
						},
						Name:       "Asset 1",
						Price:      99.99,
						ProductURL: "https://example.com/1",
						Status:     enums.StatusAssetCompleted,
						ImageCount: 2,
						CreatedAt:  "2024-01-01T00:00:00Z",
					},
				}
				images := []models.ImageEntity{
					{
						BaseItem: models.BaseItem{
							PK: "ASSET#asset-1",
							SK: "IMAGE#1",
						},
						FileKey: "assets/asset-1/1.jpg",
						Status:  enums.StatusImageCompleted,
						Order:   1,
					},
					{
						BaseItem: models.BaseItem{
							PK: "ASSET#asset-1",
							SK: "IMAGE#2",
						},
						FileKey: "assets/asset-1/2.jpg",
						Status:  enums.StatusImageCompleted,
						Order:   2,
					},
				}
				paginatedResult := &response.PaginatedResult[models.AssetEntity]{
					Data: assets,
					Meta: response.Meta{Limit: 10, HasMore: false},
				}
				assetRepo.On("FindAll", mock.Anything, 10, "").Return(paginatedResult, nil)
				imageRepo.On("FindByAssetID", mock.Anything, "ASSET#asset-1").Return(images, nil)
			},
			wantErr: false,
			wantLen: 1,
		},
		{
			name:   "with pagination",
			limit:  10,
			cursor: "",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, _ *repoMocks.MockImageRepository) {
				paginatedResult := &response.PaginatedResult[models.AssetEntity]{
					Data: []models.AssetEntity{},
					Meta: response.Meta{Limit: 10, LastKey: "next-cursor", HasMore: true},
				}
				assetRepo.On("FindAll", mock.Anything, 10, "").Return(paginatedResult, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, assetRepo, imageRepo, _ := setupAssetService(t)
			tt.setupMock(assetRepo, imageRepo)

			result, appErr := svc.ListAssets(context.Background(), tt.limit, tt.cursor)

			if tt.wantErr {
				assert.Nil(t, result)
				assert.NotNil(t, appErr)
				if tt.wantCode > 0 {
					assert.Equal(t, tt.wantCode, appErr.Code)
				}
			} else {
				assert.Nil(t, appErr)
				assert.NotNil(t, result)
				if tt.wantLen > 0 {
					assert.Equal(t, tt.wantLen, len(result.Data))
				}
				if tt.limit == 0 {
					assetRepo.AssertCalled(t, "FindAll", mock.Anything, 20, "")
				}
			}
		})
	}
}
