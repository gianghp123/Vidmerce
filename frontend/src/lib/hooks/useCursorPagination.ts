import { useState, useCallback } from "react";

export interface CursorPaginationResult<T> {
  data: T[];
  isLoading: boolean;
  error: Error | null;
  lastKey: string | null;
  hasMore: boolean;
  fetchInitial: () => Promise<void>;
  fetchNext: () => Promise<void>;
  setData: (data: T[]) => void;
}

export interface CursorPaginationOptions<T> {
  fetchFn: (cursor?: string | null) => Promise<{ items: T[]; nextCursor: string | null }>;
  enabled?: boolean;
}

export function useCursorPagination<T>(
  options: CursorPaginationOptions<T>
): CursorPaginationResult<T> {
  const { fetchFn, enabled = true } = options;

  const [data, setData] = useState<T[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const [lastKey, setLastKey] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(true);

  const fetchInitial = useCallback(
    async () => {
      if (!enabled) return;

      setIsLoading(true);
      setError(null);
      setData([]);
      setLastKey(null);
      setHasMore(true);

      try {
        const result = await fetchFn(null);
        setData(result.items);
        setLastKey(result.nextCursor);
        setHasMore(result.nextCursor !== null);
      } catch (err) {
        if (err instanceof Error && err.name !== "AbortError") {
          setError(err instanceof Error ? err : new Error("Failed to fetch initial data"));
        }
      } finally {
        setIsLoading(false);
      }
    },
    [fetchFn, enabled]
  );

  const fetchNext = useCallback(
    async () => {
      if (!enabled || !hasMore || isLoading || !lastKey) return;

      setIsLoading(true);
      setError(null);

      try {
        const result = await fetchFn(lastKey);
        setData((prev) => [...prev, ...result.items]);
        setLastKey(result.nextCursor);
        setHasMore(result.nextCursor !== null);
      } catch (err) {
        if (err instanceof Error && err.name !== "AbortError") {
          setError(err instanceof Error ? err : new Error("Failed to fetch more data"));
        }
      } finally {
        setIsLoading(false);
      }
    },
    [fetchFn, enabled, hasMore, isLoading, lastKey]
  );

  return {
    data,
    isLoading,
    error,
    lastKey,
    hasMore,
    fetchInitial,
    fetchNext,
    setData,
  };
}