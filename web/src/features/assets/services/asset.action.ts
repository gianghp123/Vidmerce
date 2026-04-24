'use server'

import { apiFetch } from "@/lib/api-fetch";
import type { BaseResponse } from "@/types/base.model";
import type { ConfirmUploadDto } from "../dtos/req/create-asset.req.dto";
import type { IAsset } from "@/types/models/asset.model";
import type { CreateAssetRes } from "../dtos/res/create-asset.res.dto";

export async function createAssetGetUploadUrl(
  payload: { imageCount: number }
): Promise<BaseResponse<CreateAssetRes>> {
  return apiFetch<CreateAssetRes>("/assets", {
    withCredentials: true,
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
  assetId: string,
  body: ConfirmUploadDto
): Promise<BaseResponse<IAsset>> {
  return apiFetch<IAsset>(`/assets/${assetId}/confirm`, {
    withCredentials: true,
    method: "POST",
    body: JSON.stringify(body),
  });
}


export async function deleteAssetImage(
  assetId: string,
  imageId: string
): Promise<BaseResponse<void>> {
  return apiFetch<void>(`/assets/${assetId}/images/${imageId}`, {
    withCredentials: true,
    method: "DELETE",
  });
}