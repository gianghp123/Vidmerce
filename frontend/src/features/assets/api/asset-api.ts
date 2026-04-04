import type { Asset } from "../models/asset.model";
import type { BaseResponse } from "@/lib/base.model";

const API_BASE_URL = import.meta.env.VITE_API_URL || "/api/v1";

export interface FetchAssetsParams {
  page?: number;
  limit?: number;
  status?: string;
  search?: string;
}

export async function fetchAssets(
  params: FetchAssetsParams = {}
): Promise<BaseResponse<Asset[]>> {
  const searchParams = new URLSearchParams();
  
  if (params.page) searchParams.set("page", params.page.toString());
  if (params.limit) searchParams.set("limit", params.limit.toString());
  if (params.status) searchParams.set("status", params.status);
  if (params.search) searchParams.set("search", params.search);

  const response = await fetch(`${API_BASE_URL}/assets?${searchParams}`);
  
  if (!response.ok) {
    throw new Error(`Failed to fetch assets: ${response.statusText}`);
  }
  
  return response.json();
}

export async function fetchAsset(id: string): Promise<BaseResponse<Asset>> {
  const response = await fetch(`${API_BASE_URL}/assets/${id}`);
  
  if (!response.ok) {
    throw new Error(`Failed to fetch asset: ${response.statusText}`);
  }
  
  return response.json();
}