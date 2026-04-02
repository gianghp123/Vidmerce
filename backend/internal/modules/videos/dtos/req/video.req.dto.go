package req

type GetVideosQuery struct {
	Limit   int    `form:"limit,default=10"`
	LastKey string `form:"last_key"`
}

type GetVideoByIDParam struct {
	ID string `uri:"id" binding:"required"`
}

type VideoItemReq struct {
	AssetID    string `json:"assetId" binding:"required"`
	Duration   int    `json:"duration" binding:"required,gt=0"`
	Transition string `json:"transition" binding:"required,oneof=FADE SLIDE"`
}

type CreateVideoReq struct {
	Title string         `json:"title" binding:"required"`
	Style string         `json:"style" binding:"required,oneof=KEN_BURNS"`
	Items []VideoItemReq `json:"items" binding:"required,min=1,dive"`
}
