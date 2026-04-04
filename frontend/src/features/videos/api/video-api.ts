import type { Video } from "../models/video.model";
import type { BaseResponse } from "@/lib/base.model";

const API_BASE_URL = import.meta.env.VITE_API_URL || "/api/v1";

export interface FetchVideosParams {
  page?: number;
  limit?: number;
  status?: string;
  search?: string;
}

export async function fetchVideos(
  params: FetchVideosParams = {}
): Promise<BaseResponse<Video[]>> {
  const searchParams = new URLSearchParams();
  
  if (params.page) searchParams.set("page", params.page.toString());
  if (params.limit) searchParams.set("limit", params.limit.toString());
  if (params.status) searchParams.set("status", params.status);
  if (params.search) searchParams.set("search", params.search);

  const response = await fetch(`${API_BASE_URL}/videos?${searchParams}`);
  
  if (!response.ok) {
    throw new Error(`Failed to fetch videos: ${response.statusText}`);
  }
  
  return response.json();
}

export async function fetchVideo(id: string): Promise<BaseResponse<Video>> {
  const response = await fetch(`${API_BASE_URL}/videos/${id}`);
  
  if (!response.ok) {
    throw new Error(`Failed to fetch video: ${response.statusText}`);
  }
  
  return response.json();
}

export async function retryVideo(id: string): Promise<BaseResponse<Video>> {
  const response = await fetch(`${API_BASE_URL}/videos/${id}/retry`, {
    method: "POST",
  });
  
  if (!response.ok) {
    throw new Error(`Failed to retry video: ${response.statusText}`);
  }
  
  return response.json();
}
