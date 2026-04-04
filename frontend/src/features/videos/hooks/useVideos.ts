import { useState, useEffect } from "react";
import { fetchVideos, fetchVideo, retryVideo, type FetchVideosParams } from "../api/video-api";
import type { Video } from "../models/video.model";

export interface UseVideosOptions extends FetchVideosParams {
  enabled?: boolean;
}

export function useVideos(options: UseVideosOptions = {}) {
  const { page = 1, limit = 12, status, search, enabled = true } = options;
  
  const [data, setData] = useState<{ data: Video[]; total: number; page: number; limit: number; totalPages: number } | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    if (!enabled) return;
    
    const fetchData = async () => {
      setIsLoading(true);
      setError(null);
      
      try {
        const result = await fetchVideos({ page, limit, status, search });
        if (result.error) {
          throw new Error(result.error.message);
        }
        const items = result.data || [];
        const meta = result.meta || {};
        setData({
          data: items,
          total: (meta as any).total || 0,
          page: (meta as any).page || 1,
          limit: (meta as any).limit || 12,
          totalPages: (meta as any).totalPages || 1,
        });
      } catch (err) {
        setError(err instanceof Error ? err : new Error("Failed to fetch videos"));
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchData();
  }, [page, limit, status, search, enabled]);

  return {
    videos: data?.data ?? [],
    total: data?.total ?? 0,
    page: data?.page ?? 1,
    totalPages: data?.totalPages ?? 1,
    isLoading,
    error,
  };
}

export function useVideo(id: string, enabled = true) {
  const [video, setVideo] = useState<Video | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    if (!enabled || !id) return;
    
    const fetchData = async () => {
      setIsLoading(true);
      setError(null);
      
      try {
        const result = await fetchVideo(id);
        if (result.error) {
          throw new Error(result.error.message);
        }
        setVideo(result.data);
      } catch (err) {
        setError(err instanceof Error ? err : new Error("Failed to fetch video"));
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchData();
  }, [id, enabled]);

  return { video, isLoading, error };
}

export function useRetryVideo() {
  const [isRetrying, setIsRetrying] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const retry = async (id: string) => {
    setIsRetrying(true);
    setError(null);
    
    try {
      const result = await retryVideo(id);
      if (result.error) {
        throw new Error(result.error.message);
      }
      return result.data;
    } catch (err) {
      setError(err instanceof Error ? err : new Error("Failed to retry video"));
      throw err;
    } finally {
      setIsRetrying(false);
    }
  };

  return { retry, isRetrying, error };
}
