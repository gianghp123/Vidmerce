import { apiFetch } from "@/lib/api-fetch";
import type { BaseResponse } from "@/lib/base.model";
import type { CreateAssetDto } from "../dtos/req/create-asset.req.dto";

export async function createAssetGetUploadUrl(
  payload: CreateAssetDto
): Promise<BaseResponse<{ assetId: string; uploads: Array<{ order: number; uploadUrl: string }> }>> {
  return apiFetch<{ assetId: string; uploads: Array<{ order: number; uploadUrl: string }> }>("/assets", {
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
): Promise<BaseResponse<{ assetId: string; status: string }>> {
  return apiFetch<{ assetId: string; status: string }>(`/assets/${assetId}/confirm`, {
    method: "POST",
  });
}