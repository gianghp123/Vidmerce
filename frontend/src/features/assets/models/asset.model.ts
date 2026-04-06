export type AssetStatus =
  | "UPLOADING"
  | "DRAFT"
  | "COMPLETED"
  | "PARTIAL"
  | "FAILED"

export type ImageStatus = 
  | "COMPLETED"
  | "UPLOADING"
  | "FAILED"

export interface AssetImage {
  imageId: string
  imageUrl: string;
  order: number;
  status: ImageStatus
}

export interface AssetUpload {
  fileKey: string;
  uploadUrl: string;
  order: number;
  expiresIn: number;
}

export interface Asset {
  assetId: string;
  name: string;
  description?: string;
  price: number;
  images: AssetImage[];
  productUrl: string;
  status: AssetStatus;
  createdAt: string;
}

export interface AssetPreview {
  assetId: string;
  name: string;
  price: number;
  image: AssetImage;
  productUrl: string;
  status: AssetStatus;
  createdAt: string;
}

export interface AssetUploadResponse {
  assetId: string;
  status: AssetStatus;
  uploads: AssetUpload[];
}

export interface AssetConfirmResponse {
  assetId: string;
  status: AssetStatus;
  images: AssetImage[];
}