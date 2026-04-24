import type { IAsset } from "@/types/models/asset.model";
import type { BaseResponse } from "@/types/base.model";
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
  return apiFetch<IAsset[]>("/assets", { query: params, withCredentials: true });
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
