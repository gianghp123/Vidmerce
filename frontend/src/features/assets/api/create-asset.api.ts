import { apiFetch } from "@/lib/api-fetch";
import type { BaseResponse } from "@/lib/base.model";
import type { CreateAssetDto } from "../dtos/create-asset.dto";
import type { AssetConfirmResponse, AssetUploadResponse } from "../models/asset.model";


export async function createAssetGetUploadUrl(
  payload: CreateAssetDto
): Promise<BaseResponse<AssetUploadResponse>> {
  return apiFetch<AssetUploadResponse>("/assets", {
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