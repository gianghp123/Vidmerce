package res

import "github.com/gianghp123/Vidmerce/backend/services/internal/core/enums"

type VideoResponse struct {
	ID           string            `json:"id"`
	Title        string            `json:"title"`
	Status       enums.VideoStatus `json:"status"`
	VideoURL     string            `json:"video_url"`
	ThumbnailURL string            `json:"thumbnail_url"`
	CreatedAt    string            `json:"created_at"`
}

type VideoDetailResponse struct {
	VideoResponse
	Hotspots []HotspotResponse `json:"hotspots"`
}

type HotspotResponse struct {
	X     float64          `json:"x"`
	Y     float64          `json:"y"`
	Asset AssetSnapshotDTO `json:"asset"`
}

type AssetSnapshotDTO struct {
	AssetID    string  `json:"asset_id"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	ImageURL   string  `json:"image_url"`
	ProductURL string  `json:"product_url"`
}

type CreateVideoRes struct {
	VideoID string            `json:"videoId"`
	Status  enums.VideoStatus `json:"status"`
	Message string            `json:"message"`
}
