import type { Video } from "../models/video.model";
import type { BaseResponse } from "@/lib/base.model";
import { apiFetch } from "@/lib/api-fetch";

export interface FetchVideosParams {
  cursor?: string | null;
  limit?: number;
  status?: string;
  search?: string;
}

export async function fetchVideos(
  params: FetchVideosParams = {}
): Promise<BaseResponse<Video[]>> {
  return apiFetch<Video[]>("/videos", { query: params });
}

export async function fetchVideo(id: string): Promise<BaseResponse<Video>> {
  return apiFetch<Video>(`/videos/${id}`);
}

export async function retryVideo(id: string): Promise<BaseResponse<Video>> {
  return apiFetch<Video>(`/videos/${id}/retry`, { method: "POST" });
}
