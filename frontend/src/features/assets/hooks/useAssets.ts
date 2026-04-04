import { useState, useEffect } from "react";
import { fetchAssets, fetchAsset, type FetchAssetsParams } from "../api/asset-api";
import type { Asset } from "../models/asset.model";

export interface UseAssetsOptions extends FetchAssetsParams {
  enabled?: boolean;
}

export function useAssets(options: UseAssetsOptions = {}) {
  const { page = 1, limit = 12, status, search, enabled = true } = options;
  
  const [data, setData] = useState<{ data: Asset[]; total: number; page: number; limit: number; totalPages: number } | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    if (!enabled) return;
    
    const fetchData = async () => {
      setIsLoading(true);
      setError(null);
      
      try {
        const result = await fetchAssets({ page, limit, status, search });
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
        setError(err instanceof Error ? err : new Error("Failed to fetch assets"));
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchData();
  }, [page, limit, status, search, enabled]);

  return {
    assets: data?.data ?? [],
    total: data?.total ?? 0,
    page: data?.page ?? 1,
    totalPages: data?.totalPages ?? 1,
    isLoading,
    error,
  };
}

export function useAsset(id: string, enabled = true) {
  const [asset, setAsset] = useState<Asset | null>(null);
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