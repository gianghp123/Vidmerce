import { useCallback, useEffect, useState } from "react";
import { fetchAssets, fetchAsset, type FetchAssetsParams } from "../services/asset.get";
import type { IAsset } from "../../../lib/models/asset.model";
import { useCursorPagination } from "@/lib/hooks/useCursorPagination";

export interface UseAssetsOptions extends FetchAssetsParams {
  enabled?: boolean;
}

export function useAssets(options: UseAssetsOptions = {}) {
  const { cursor, limit = 12, status, search, enabled = true } = options;

  const fetchAssetsFn = useCallback(
    async (cursorParam?: string | null) => {
      const result = await fetchAssets({
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
    data: assets,
    isLoading,
    error,
    hasMore,
    fetchInitial,
    fetchNext,
  } = useCursorPagination<IAsset>({
    fetchFn: fetchAssetsFn,
    enabled,
  });

  useEffect(() => {
    fetchInitial();
  }, [enabled, cursor, limit, status, search]);

  return {
    assets,
    isLoading,
    error,
    hasMore,
    fetchInitial,
    fetchNext,
  };
}

export function useAsset(id: string, enabled = true) {
  const [asset, setAsset] = useState<IAsset | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    if (!enabled || !id) return;
    
    const fetchData = async () => {
      setIsLoading(true);
      setError(null);
      
      try {
        const result = await fetchAsset(id);
        if (result.error) {
          throw new Error(result.error.message);
        }
        setAsset(result.data);
      } catch (err) {
        setError(err instanceof Error ? err : new Error("Failed to fetch asset"));
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchData();
  }, [id, enabled]);

  return { asset, isLoading, error };
}