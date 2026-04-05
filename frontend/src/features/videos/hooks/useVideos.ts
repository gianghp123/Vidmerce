import { useCallback, useEffect, useState } from "react";
import { fetchVideos, fetchVideo, retryVideo, type FetchVideosParams } from "../api/video.api";
import type { Video } from "../models/video.model";
import { useCursorPagination } from "@/lib/hooks/useCursorPagination";

export interface UseVideosOptions extends FetchVideosParams {
  enabled?: boolean;
}

export function useVideos(options: UseVideosOptions = {}) {
  const { cursor, limit = 12, status, search, enabled = true } = options;

  const fetchVideosFn = useCallback(
    async (cursorParam?: string | null) => {
      const result = await fetchVideos({
        cursor: cursorParam,
        limit,
        status,
        search,
      });
      if (result.error) {
        throw new Error(result.error.message);
      }
      const items = result.data || [];
      const meta = result.meta || {};
      const nextCursor = (meta as any).lastKey ?? null;
      return { items, nextCursor };
    },
    [limit, status, search]
  );

  const {
    data: videos,
    isLoading,
    error,
    hasMore,
    fetchInitial,
    fetchNext,
  } = useCursorPagination<Video>({
    fetchFn: fetchVideosFn,
    enabled,
  });

  useEffect(() => {
    fetchInitial();
  }, [enabled, cursor, limit, status, search]);

  return {
    videos,
    isLoading,
    error,
    hasMore,
    fetchInitial,
    fetchNext,
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