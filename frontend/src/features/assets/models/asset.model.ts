export type AssetStatus =
  | "UPLOADING"
  | "DRAFT"
  | "COMPLETED"
  | "PARTIAL"
  | "FAILED"

export interface AssetImage {
  imageUrl: string;
  order: number;
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