import type { AssetUploadResponse, AssetConfirmResponse } from "../models/asset.model";
import type { BaseResponse } from "@/lib/base.model";

const API_BASE_URL = import.meta.env.VITE_API_URL || "/api/v1";

export interface CreateAssetPayload {
  name: string;
  price: number;
  productUrl: string;
  description?: string;
}

export async function createAssetGetUploadUrl(
  payload: CreateAssetPayload
): Promise<BaseResponse<AssetUploadResponse>> {
  const response = await fetch(`${API_BASE_URL}/assets/upload-url`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    throw new Error(`Failed to get upload URL: ${response.statusText}`);
  }

  return response.json();
}

export async function uploadImageToS3(
  uploadUrl: string,
  file: File
): Promise<void> {
  const response = await fetch(uploadUrl, {
    method: "PUT",
    body: file,
    headers: {
      "Content-Type": file.type,
    },
  });

  if (!response.ok) {
    throw new Error(`Failed to upload image: ${response.statusText}`);
  }
}

export async function confirmAssetUpload(
  assetId: string
): Promise<BaseResponse<AssetConfirmResponse>> {
  const response = await fetch(`${API_BASE_URL}/assets/${assetId}/confirm`, {
    method: "POST",
  });

  if (!response.ok) {
    throw new Error(`Failed to confirm asset upload: ${response.statusText}`);
  }

  return response.json();
}