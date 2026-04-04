package req

type GetVideosQuery struct {
	Limit   int    `form:"limit,default=10" example:"10"`
	LastKey string `form:"last_key" example:"eyJ2b3RlSWQiOjF9"`
}

type GetVideoByIDParam struct {
	ID string `uri:"id" binding:"required" example:"vid_123456"`
}

type VideoItemReq struct {
	AssetID    string `json:"assetId" binding:"required" example:"asset_789"`
	Duration   int    `json:"duration" binding:"required,gt=0" example:"30"`
	Transition string `json:"transition" binding:"required,oneof=FADE SLIDE" example:"FADE"`
}

type CreateVideoReq struct {
	Title string         `json:"title" binding:"required" example:"Summer Promotion"`
	Style string         `json:"style" binding:"required,oneof=KEN_BURNS" example:"KEN_BURNS"`
	Items []VideoItemReq `json:"items" binding:"required,min=1,dive"`
}
