import type { IAsset } from "../../../lib/models/asset.model";
import type { BaseResponse } from "@/lib/base.model";
import { apiFetch } from "@/lib/api-fetch";

export interface FetchAssetsParams {
  cursor?: string | null;
  limit?: number;
  status?: string;
  search?: string;
}

export async function fetchAssets(
  params: FetchAssetsParams = {}
): Promise<BaseResponse<IAsset[]>> {
  return apiFetch<IAsset[]>("/assets", { query: params });
}

export async function fetchAsset(id: string): Promise<BaseResponse<IAsset>> {
  return apiFetch<IAsset>(`/assets/${id}`);
}

export async function getImageUploadUrl(
  assetId: string,
  fileName: string
): Promise<BaseResponse<{ uploadUrl: string; imageId: string }>> {
  return apiFetch<{ uploadUrl: string; imageId: string }>(`/assets/${assetId}/images/upload-url`, {
    query: { fileName },
  });
}

export async function deleteAssetImage(
  assetId: string,
  imageId: string
): Promise<BaseResponse<void>> {
  return apiFetch<void>(`/assets/${assetId}/images/${imageId}`, {
    method: "DELETE",
  });
}