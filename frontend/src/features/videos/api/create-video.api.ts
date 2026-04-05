import type { VideoBuilder } from "../models/video-sequence.model";
import type { VideoCreateResponse } from "../models/video.model";
import type { BaseResponse } from "@/lib/base.model";
import { apiFetch } from "@/lib/api-fetch";

export interface CreateVideoPayload {
  title: string;
  style: VideoBuilder["style"];
  sequences: {
    assetId: string;
    duration: number;
    transition: string;
  }[];
}

export async function createVideo(
  payload: CreateVideoPayload
): Promise<BaseResponse<VideoCreateResponse>> {
  return apiFetch<VideoCreateResponse>("/videos", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });
}