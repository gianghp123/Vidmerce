import { useState, useCallback } from "react";
import {
  createAssetGetUploadUrl,
  uploadImageToS3,
  confirmAssetUpload,
  type CreateAssetPayload,
} from "../api/create-asset.api";
import type { Asset } from "../models/asset.model";

interface ImageFile {
  id: string;
  file: File;
  preview: string;
}

interface UseCreateAssetOptions {
  onSuccess?: (asset: Asset) => void;
  onError?: (error: Error) => void;
}

export function useCreateAsset(options: UseCreateAssetOptions = {}) {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const [asset, setAsset] = useState<Asset | null>(null);

  const createAsset = useCallback(
    async (payload: CreateAssetPayload, images: ImageFile[]) => {
      setIsLoading(true);
      setError(null);

      try {
        const uploadUrlResponse = await createAssetGetUploadUrl(payload);

        if (uploadUrlResponse.error || !uploadUrlResponse.data) {
          throw new Error(uploadUrlResponse.error?.message || "Failed to get upload URL");
        }

        const { assetId, uploads } = uploadUrlResponse.data;

        for (let i = 0; i < images.length; i++) {
          const upload = uploads.find((u: { order: number }) => u.order === i);
          if (upload) {
            await uploadImageToS3(upload.uploadUrl, images[i].file);
          }
        }

        const confirmResponse = await confirmAssetUpload(assetId);

        if (confirmResponse.error || !confirmResponse.data) {
          throw new Error(confirmResponse.error?.message || "Failed to confirm asset upload");
        }

        setAsset(confirmResponse.data as unknown as Asset);
        options.onSuccess?.(confirmResponse.data as unknown as Asset);
      } catch (err) {
        const error = err instanceof Error ? err : new Error("Failed to create asset");
        setError(error);
        options.onError?.(error);
      } finally {
        setIsLoading(false);
      }
    },
    [options]
  );

  return {
    createAsset,
    isLoading,
    error,
    asset,
    reset: () => {
      setAsset(null);
      setError(null);
    },
  };
}