package services

import (
	"context"
	"errors"
	"testing"

	"github.com/gianghp123/Vidmerce/backend/internal/core"
	"github.com/gianghp123/Vidmerce/backend/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/internal/database/models"
	"github.com/gianghp123/Vidmerce/backend/internal/modules/assets/dtos/req"
	repoMocks "github.com/gianghp123/Vidmerce/backend/internal/modules/assets/repositories"
	"github.com/gianghp123/Vidmerce/backend/internal/storage"
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

func createAuthContext(userID string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, core.UserIDKey, userID)
	ctx = context.WithValue(ctx, core.RoleKey, enums.UserRoleUser)
	return ctx
}

// ==================== CreateAsset Tests ====================

func TestCreateAsset_PresignedUrls(t *testing.T) {
	tests := []struct {
		name          string
		req           req.CreateAssetReq
		setupMock     func(*repoMocks.MockAssetRepository, *repoMocks.MockImageRepository, *storage.MockStorage)
		wantErr       bool
		wantCode      int
		wantUploadLen int
	}{
		{
			name: "success",
			req: req.CreateAssetReq{
				ImageCount: 1,
			},
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				storageMock.On("GeneratePresignedUploadURL", mock.Anything, mock.Anything, "image/jpeg", mock.Anything).Return("https://upload.url/1", nil)
			},
			wantErr:       false,
			wantUploadLen: 1,
		},
		{
			name: "success multiple images",
			req: req.CreateAssetReq{
				ImageCount: 3,
			},
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				storageMock.On("GeneratePresignedUploadURL", mock.Anything, mock.Anything, "image/jpeg", mock.Anything).Return("https://upload.url", nil)
			},
			wantErr:       false,
			wantUploadLen: 3,
		},
		{
			name: "default image count",
			req: req.CreateAssetReq{
				ImageCount: 0,
			},
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				storageMock.On("GeneratePresignedUploadURL", mock.Anything, mock.Anything, "image/jpeg", mock.Anything).Return("https://upload.url", nil)
			},
			wantErr:       false,
			wantUploadLen: 1,
		},
		{
			name: "exceeds max images",
			req: req.CreateAssetReq{
				ImageCount: 6,
			},
			setupMock: nil,
			wantErr:   true,
		},
		{
			name: "storage error",
			req: req.CreateAssetReq{
				ImageCount: 1,
			},
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				storageMock.On("GeneratePresignedUploadURL", mock.Anything, mock.Anything, "image/jpeg", mock.Anything).Return("", errors.New("storage error"))
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

			ctx := createAuthContext("test-user-id")
			result, appErr := svc.CreateAsset(ctx, tt.req)

			if tt.wantErr {
				assert.Nil(t, result)
				assert.NotNil(t, appErr)
				if tt.wantCode > 0 {
					assert.Equal(t, tt.wantCode, appErr.Code)
				}
			} else {
				assert.Nil(t, appErr)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AssetID)
				if tt.wantUploadLen > 0 {
					assert.Equal(t, tt.wantUploadLen, len(result.Uploads))
				}
			}
		})
	}
}

// ==================== ConfirmUpload Tests ====================

func TestConfirmUpload(t *testing.T) {
	tests := []struct {
		name      string
		assetID   string
		req       req.ConfirmUploadReq
		setupMock func(*repoMocks.MockAssetRepository, *repoMocks.MockImageRepository, *storage.MockStorage)
		wantErr   bool
		wantCode  int
	}{
		{
			name:    "success",
			assetID: "test-asset-id",
			req: req.ConfirmUploadReq{
				Name:       "Test Asset",
				Price:      99.99,
				ProductURL: "https://example.com",
				FileKeys:   []string{"user/test-user-id/asset/test-asset-id/1.jpg"},
			},
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				storageMock.On("ObjectExists", mock.Anything, "user/test-user-id/asset/test-asset-id/1.jpg").Return(true, nil)
				assetRepo.On("CreateWithImages", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "image not uploaded",
			assetID: "test-asset-id",
			req: req.ConfirmUploadReq{
				Name:       "Test Asset",
				Price:      99.99,
				ProductURL: "https://example.com",
				FileKeys:   []string{"user/test-user-id/asset/test-asset-id/1.jpg"},
			},
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				storageMock.On("ObjectExists", mock.Anything, "user/test-user-id/asset/test-asset-id/1.jpg").Return(false, nil)
			},
			wantErr:  true,
			wantCode: 400,
		},
		{
			name:    "storage error",
			assetID: "test-asset-id",
			req: req.ConfirmUploadReq{
				Name:       "Test Asset",
				Price:      99.99,
				ProductURL: "https://example.com",
				FileKeys:   []string{"user/test-user-id/asset/test-asset-id/1.jpg"},
			},
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, imageRepo *repoMocks.MockImageRepository, storageMock *storage.MockStorage) {
				storageMock.On("ObjectExists", mock.Anything, "user/test-user-id/asset/test-asset-id/1.jpg").Return(false, errors.New("storage error"))
			},
			wantErr:  true,
			wantCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, assetRepo, imageRepo, storageMock := setupAssetService(t)
			tt.setupMock(assetRepo, imageRepo, storageMock)

			ctx := createAuthContext("test-user-id")
			result, appErr := svc.ConfirmUpload(ctx, tt.assetID, tt.req)

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
			}
		})
	}
}

// ==================== GetAsset Tests ====================

func TestGetAsset(t *testing.T) {
	asset := &models.AssetEntity{
		BaseItem: models.BaseItem{
			Pk: "USER#test-user-id",
			Sk: "ASSET#test-id",
		},
		Name:       "Test Asset",
		Price:      99.99,
		ProductURL: "https://example.com",
		Status:     enums.AssetStatusCompleted,
		ImageCount: 1,
		CreatedAt:  "2024-01-01T00:00:00Z",
	}

	tests := []struct {
		name      string
		userID    string
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
							Pk: "ASSET#test-id",
							Sk: "IMAGE#1",
						},
						FileKey: "assets/test-id/1.jpg",
						Order:   1,
					},
				}
				assetRepo.On("FindByID", mock.Anything, "test-id").Return(asset, nil)
				imageRepo.On("FindByAssetID", mock.Anything, "test-id").Return(images, nil)
			},
			wantErr: false,
			wantID:  "test-id",
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
						Pk: "USER#test-user-id",
						Sk: "ASSET#test-id",
					},
					Name:       "Test Asset",
					Price:      99.99,
					ProductURL: "https://example.com",
					Status:     enums.AssetStatusCompleted,
					ImageCount: 2,
					CreatedAt:  "2024-01-01T00:00:00Z",
				}
				images := []models.ImageEntity{
					{
						BaseItem: models.BaseItem{
							Pk: "USER#test-user-id",
							Sk: "ASSET#test-id",
						},
						FileKey: "assets/test-id/1.jpg",
						Order:   1,
					},
					{
						BaseItem: models.BaseItem{
							Pk: "USER#test-user-id",
							Sk: "ASSET#test-id",
						},
						FileKey: "assets/test-id/2.jpg",
						Order:   2,
					},
				}
				assetRepo.On("FindByID", mock.Anything, "test-id").Return(asset, nil)
				imageRepo.On("FindByAssetID", mock.Anything, "test-id").Return(images, nil)
			},
			wantErr: false,
			wantLen: 2,
		},
		{
			name:    "ownership check fails",
			assetID: "test-id",
			setupMock: func(assetRepo *repoMocks.MockAssetRepository, _ *repoMocks.MockImageRepository) {
				otherUserAsset := &models.AssetEntity{
					BaseItem: models.BaseItem{
						Pk: "USER#other-user-id",
						Sk: "ASSET#test-id",
					},
				}
				assetRepo.On("FindByID", mock.Anything, "test-id").Return(otherUserAsset, nil)
			},
			wantErr:  true,
			wantCode: 403,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, assetRepo, imageRepo, _ := setupAssetService(t)
			tt.setupMock(assetRepo, imageRepo)

			ctx := createAuthContext("test-user-id")
			result, appErr := svc.GetAsset(ctx, tt.assetID)

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
					assert.NotEmpty(t, result.Images[0].ImageURL)
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
							Pk:     "ASSET#asset-1",
							Sk:     "METADATA",
							Gsi1Pk: "ENTITY#ASSET",
							Gsi1Sk: "asset-1",
						},
						Name:       "Asset 1",
						Price:      99.99,
						ProductURL: "https://example.com/1",
						Status:     enums.AssetStatusCompleted,
						ImageCount: 1,
						CreatedAt:  "2024-01-01T00:00:00Z",
					},
				}
				images := []models.ImageEntity{
					{
						BaseItem: models.BaseItem{
							Pk: "ASSET#asset-1",
							Sk: "IMAGE#1",
						},
						FileKey: "assets/asset-1/1.jpg",
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
							Pk:     "ASSET#asset-1",
							Sk:     "METADATA",
							Gsi1Pk: "ENTITY#ASSET",
							Gsi1Sk: "asset-1",
						},
						Name:       "Asset 1",
						Price:      99.99,
						ProductURL: "https://example.com/1",
						Status:     enums.AssetStatusCompleted,
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
							Pk:     "ASSET#asset-1",
							Sk:     "METADATA",
							Gsi1Pk: "ENTITY#ASSET",
							Gsi1Sk: "asset-1",
						},
						Name:       "Asset 1",
						Price:      99.99,
						ProductURL: "https://example.com/1",
						Status:     enums.AssetStatusCompleted,
						ImageCount: 2,
						CreatedAt:  "2024-01-01T00:00:00Z",
					},
				}
				images := []models.ImageEntity{
					{
						BaseItem: models.BaseItem{
							Pk: "ASSET#asset-1",
							Sk: "IMAGE#1",
						},
						FileKey: "assets/asset-1/1.jpg",
						Order:   1,
					},
					{
						BaseItem: models.BaseItem{
							Pk: "ASSET#asset-1",
							Sk: "IMAGE#2",
						},
						FileKey: "assets/asset-1/2.jpg",
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

			ctx := createAuthContext("test-user-id")
			result, appErr := svc.ListAssets(ctx, tt.limit, tt.cursor)

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
