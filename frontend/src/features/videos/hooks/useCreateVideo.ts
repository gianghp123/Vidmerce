import { useState, useCallback } from "react";
import { createVideo, type CreateVideoPayload } from "../api/create-video.api";
import type { Video } from "../models/video.model";

interface UseCreateVideoOptions {
  onSuccess?: (video: Video) => void;
  onError?: (error: Error) => void;
}

export function useCreateVideo(options: UseCreateVideoOptions = {}) {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const [video, setVideo] = useState<Video | null>(null);

  const create = useCallback(
    async (payload: CreateVideoPayload) => {
      setIsLoading(true);
      setError(null);

      try {
        const response = await createVideo(payload);

        if (response.error || !response.data) {
          throw new Error(response.error?.message || "Failed to create video");
        }

        const newVideo: Video = {
          id: response.data.videoId,
          title: payload.title,
          status: response.data.status,
          video_url: "",
          thumbnail_url: "",
          created_at: new Date().toISOString(),
          hotspots: [],
        };

        setVideo(newVideo);
        options.onSuccess?.(newVideo);
      } catch (err) {
        const error = err instanceof Error ? err : new Error("Failed to create video");
        setError(error);
        options.onError?.(error);
      } finally {
        setIsLoading(false);
      }
    },
    [options]
  );

  return {
    create,
    isLoading,
    error,
    video,
    reset: () => {
      setVideo(null);
      setError(null);
    },
  };
}