package services

import (
	"context"
	"errors"
	"testing"

	"github.com/gianghp123/Vidmerce/backend/services/internal/core/enums"
	"github.com/gianghp123/Vidmerce/backend/services/internal/core/response"
	"github.com/gianghp123/Vidmerce/backend/services/internal/database/models"
	"github.com/gianghp123/Vidmerce/backend/services/internal/modules/videos/dtos/req"
	repoMocks "github.com/gianghp123/Vidmerce/backend/services/internal/modules/videos/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupVideoService(t *testing.T) (*videoService, *repoMocks.MockVideoRepository) {
	t.Helper()
	videoRepo := repoMocks.NewMockVideoRepository()
	svc := NewVideoService(videoRepo).(*videoService)
	return svc, videoRepo
}

// ==================== CreateVideo Tests ====================

func TestCreateVideo(t *testing.T) {
	tests := []struct {
		name    string
		req     req.CreateVideoReq
		wantErr bool
	}{
		{
			name: "not implemented",
			req: req.CreateVideoReq{
				Title: "Test Video",
				Style: "KEN_BURNS",
				Items: []req.VideoItemReq{
					{
						AssetID:    "asset-1",
						Duration:   5,
						Transition: "FADE",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, _ := setupVideoService(t)

			result, appErr := svc.CreateVideo(context.Background(), tt.req)

			if tt.wantErr {
				assert.Nil(t, result)
				assert.NotNil(t, appErr)
				assert.Equal(t, 500, appErr.Code)
			} else {
				assert.NotNil(t, result)
				assert.Nil(t, appErr)
			}
		})
	}
}

// ==================== ListVideos Tests ====================

func TestListVideos(t *testing.T) {
	tests := []struct {
		name         string
		limit        int
		cursor       string
		setupMock    func(*repoMocks.MockVideoRepository)
		wantErr      bool
		wantDataLen  int
		wantHasMore  bool
		wantLastKey  string
		wantCallArgs []any
	}{
		{
			name:   "success",
			limit:  10,
			cursor: "",
			setupMock: func(repo *repoMocks.MockVideoRepository) {
				videos := []models.VideoMetadataEntity{
					{
						BaseItem: models.BaseItem{
							PK:     "VIDEO#video-1",
							SK:     "METADATA",
							GSI1PK: "ENTITY#VIDEO",
							GSI1SK: "video-1",
						},
						Title:        "Test Video",
						Status:       enums.StatusCompleted,
						VideoURL:     "https://example.com/video.mp4",
						ThumbnailURL: "https://example.com/thumb.jpg",
						CreatedAt:    "2024-01-01T00:00:00Z",
					},
				}
				paginatedResult := &response.PaginatedResult[models.VideoMetadataEntity]{
					Data: videos,
					Meta: response.Meta{Limit: 10, HasMore: false},
				}
				repo.On("FindAll", mock.Anything, 10, "").Return(paginatedResult, nil)
			},
			wantErr:     false,
			wantDataLen: 1,
		},
		{
			name:   "default limit",
			limit:  0,
			cursor: "",
			setupMock: func(repo *repoMocks.MockVideoRepository) {
				paginatedResult := &response.PaginatedResult[models.VideoMetadataEntity]{
					Data: []models.VideoMetadataEntity{},
					Meta: response.Meta{Limit: 10, HasMore: false},
				}
				repo.On("FindAll", mock.Anything, 10, "").Return(paginatedResult, nil)
			},
			wantErr:      false,
			wantCallArgs: []any{context.Background(), 10, ""},
		},
		{
			name:   "repo error",
			limit:  10,
			cursor: "",
			setupMock: func(repo *repoMocks.MockVideoRepository) {
				repo.On("FindAll", mock.Anything, 10, "").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:   "with pagination",
			limit:  10,
			cursor: "",
			setupMock: func(repo *repoMocks.MockVideoRepository) {
				paginatedResult := &response.PaginatedResult[models.VideoMetadataEntity]{
					Data: []models.VideoMetadataEntity{},
					Meta: response.Meta{Limit: 10, LastKey: "next-cursor", HasMore: true},
				}
				repo.On("FindAll", mock.Anything, 10, "").Return(paginatedResult, nil)
			},
			wantErr:     false,
			wantHasMore: true,
			wantLastKey: "next-cursor",
		},
		{
			name:   "empty result",
			limit:  10,
			cursor: "",
			setupMock: func(repo *repoMocks.MockVideoRepository) {
				paginatedResult := &response.PaginatedResult[models.VideoMetadataEntity]{
					Data: []models.VideoMetadataEntity{},
					Meta: response.Meta{Limit: 10, HasMore: false},
				}
				repo.On("FindAll", mock.Anything, 10, "").Return(paginatedResult, nil)
			},
			wantErr:     false,
			wantDataLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, videoRepo := setupVideoService(t)
			tt.setupMock(videoRepo)

			result, appErr := svc.ListVideos(context.Background(), tt.limit, tt.cursor)

			if tt.wantErr {
				assert.Nil(t, result)
				assert.NotNil(t, appErr)
				assert.Equal(t, 500, appErr.Code)
			} else {
				assert.Nil(t, appErr)
				assert.NotNil(t, result)
				if tt.wantDataLen > 0 {
					assert.Equal(t, tt.wantDataLen, len(result.Data))
				}
				if tt.wantHasMore {
					assert.True(t, result.Meta.HasMore)
				}
				if tt.wantLastKey != "" {
					assert.Equal(t, tt.wantLastKey, result.Meta.LastKey)
				}
			}
			if tt.wantCallArgs != nil {
				videoRepo.AssertCalled(t, "FindAll", tt.wantCallArgs...)
			}
		})
	}
}

// ==================== GetVideo Tests ====================

func TestGetVideo(t *testing.T) {
	video := &models.VideoMetadataEntity{
		BaseItem: models.BaseItem{
			PK: "VIDEO#video-1",
			SK: "METADATA",
		},
		Title:        "Test Video",
		Status:       enums.StatusCompleted,
		VideoURL:     "https://example.com/video.mp4",
		ThumbnailURL: "https://example.com/thumb.jpg",
		CreatedAt:    "2024-01-01T00:00:00Z",
	}

	interactive := &models.VideoInteractiveEntity{
		BaseItem: models.BaseItem{
			PK: "VIDEO#video-1",
			SK: "INTERACTIVE",
		},
		Hotspots: []models.Hotspot{
			{
				X: 10.5,
				Y: 20.5,
				Asset: models.AssetSnapshot{
					AssetID:    "asset-1",
					Name:       "Test Asset",
					Price:      99.99,
					ImageURL:   "https://example.com/asset.jpg",
					ProductURL: "https://example.com/product",
				},
			},
		},
	}

	tests := []struct {
		name      string
		videoID   string
		setupMock func(*repoMocks.MockVideoRepository)
		wantErr   bool
		wantID    string
		wantTitle string
	}{
		{
			name:    "success",
			videoID: "video-1",
			setupMock: func(repo *repoMocks.MockVideoRepository) {
				repo.On("FindByID", mock.Anything, "video-1").Return(video, nil)
				repo.On("FindInteractiveByVideoID", mock.Anything, "video-1").Return(interactive, nil)
			},
			wantErr:   false,
			wantID:    "VIDEO#video-1",
			wantTitle: "Test Video",
		},
		{
			name:    "success no hotspots",
			videoID: "video-1",
			setupMock: func(repo *repoMocks.MockVideoRepository) {
				repo.On("FindByID", mock.Anything, "video-1").Return(video, nil)
				repo.On("FindInteractiveByVideoID", mock.Anything, "video-1").Return(nil, nil)
			},
			wantErr: false,
			wantID:  "VIDEO#video-1",
		},
		{
			name:    "video not found",
			videoID: "non-existent",
			setupMock: func(repo *repoMocks.MockVideoRepository) {
				repo.On("FindByID", mock.Anything, "non-existent").Return(nil, nil)
			},
			wantErr: true,
		},
		{
			name:    "video repo error",
			videoID: "video-1",
			setupMock: func(repo *repoMocks.MockVideoRepository) {
				repo.On("FindByID", mock.Anything, "video-1").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:    "interactive repo error",
			videoID: "video-1",
			setupMock: func(repo *repoMocks.MockVideoRepository) {
				repo.On("FindByID", mock.Anything, "video-1").Return(video, nil)
				repo.On("FindInteractiveByVideoID", mock.Anything, "video-1").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, videoRepo := setupVideoService(t)
			tt.setupMock(videoRepo)

			result, appErr := svc.GetVideo(context.Background(), tt.videoID)

			if tt.wantErr {
				assert.Nil(t, result)
				assert.NotNil(t, appErr)
				if tt.videoID == "non-existent" {
					assert.Equal(t, 404, appErr.Code)
				} else {
					assert.Equal(t, 500, appErr.Code)
				}
			} else {
				assert.Nil(t, appErr)
				assert.NotNil(t, result)
				if tt.wantID != "" {
					assert.Equal(t, tt.wantID, result.ID)
				}
				if tt.wantTitle != "" {
					assert.Equal(t, tt.wantTitle, result.Title)
				}
			}
			videoRepo.AssertExpectations(t)
		})
	}
}
