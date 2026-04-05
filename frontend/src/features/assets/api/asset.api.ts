import type { Asset } from "../models/asset.model";
import type { BaseResponse } from "@/lib/base.model";
import { apiFetch } from "@/lib/api-fetch";

export interface FetchAssetsParams {
  page?: number;
  limit?: number;
  status?: string;
  search?: string;
}

export async function fetchAssets(
  params: FetchAssetsParams = {}
): Promise<BaseResponse<Asset[]>> {
  return apiFetch<Asset[]>("/assets", { query: params });
}

export async function fetchAsset(id: string): Promise<BaseResponse<Asset>> {
  return apiFetch<Asset>(`/assets/${id}`);
}