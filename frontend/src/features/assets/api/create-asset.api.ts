import type { AssetUploadResponse, AssetConfirmResponse } from "../models/asset.model";
import type { BaseResponse } from "@/lib/base.model";
import { apiFetch } from "@/lib/api-fetch";

export interface CreateAssetPayload {
  name: string;
  price: number;
  productUrl: string;
  description?: string;
}

export async function createAssetGetUploadUrl(
  payload: CreateAssetPayload
): Promise<BaseResponse<AssetUploadResponse>> {
  return apiFetch<AssetUploadResponse>("/assets/upload-url", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });
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
  return apiFetch<AssetConfirmResponse>(`/assets/${assetId}/confirm`, {
    method: "POST",
  });
}