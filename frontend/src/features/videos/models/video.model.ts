export type VideoStatus =
  | "PENDING"
  | "PROCESSING"
  | "COMPLETED"
  | "FAILED"

  
export interface HotspotAsset {
  asset_id: string;
  name: string;
  price: number;
  image_url: string;
  product_url: string;
}

export interface Hotspot {
  x: number;
  y: number;
  asset: HotspotAsset;
}

export interface Video {
  id: string;
  title: string;
  description?: string;
  status: VideoStatus;
  video_url: string;
  thumbnail_url: string;
  created_at: string;
  updated_at?: string;
  views?: number;
  hotspots: Hotspot[];
}

export interface VideoCreateResponse {
  videoId: string;
  status: VideoStatus;
  message: string;
}