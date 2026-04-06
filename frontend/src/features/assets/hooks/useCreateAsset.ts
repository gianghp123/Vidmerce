import { useCallback, useState } from "react";
import {
  confirmAssetUpload,
  createAssetGetUploadUrl,
  uploadImageToS3,
} from "../api/create-asset.api";
import type { CreateAssetDto } from "../dtos/create-asset.dto";
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
    async (payload: CreateAssetDto, images: ImageFile[]) => {
      setIsLoading(true);
      setError(null);

      try {
        const uploadUrlResponse = await createAssetGetUploadUrl(payload);

        if (uploadUrlResponse.error || !uploadUrlResponse.data) {
          throw new Error(uploadUrlResponse.error?.message || "Failed to get upload URL");
        }

        const { assetId, uploads } = uploadUrlResponse.data;

        const imagesWithOrder = images.map((img, index) => ({
          ...img,
          order: index,
        }));

        const results = await Promise.allSettled(
          imagesWithOrder.map((img) => {
            const upload = uploads.find(u => u.order === img.order);
            if (!upload) return Promise.reject("Missing upload URL");

            return uploadImageToS3(upload.uploadUrl, img.file);
          })
        );

        const failed = results.filter(r => r.status === "rejected");

        // still call confirm to let backend decide status
        if (failed.length > 0) {
          setError(new Error(`${failed.length} images failed`));
        }

        // all success
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