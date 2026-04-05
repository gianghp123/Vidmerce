import type { VideoBuilder } from "../models/video-sequence.model";
import type { VideoCreateResponse } from "../models/video.model";
import type { BaseResponse } from "@/lib/base.model";

const API_BASE_URL = import.meta.env.VITE_API_URL || "/api/v1";

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
  const response = await fetch(`${API_BASE_URL}/videos`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    throw new Error(`Failed to create video: ${response.statusText}`);
  }

  return response.json();
}