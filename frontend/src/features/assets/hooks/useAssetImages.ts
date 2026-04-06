import { useState } from "react";
import { toast } from "sonner";
import { getImageUploadUrl } from "../api/asset.api";
import { uploadImageToS3, confirmAssetUpload } from "../api/create-asset.api";
import { MAX_IMAGES } from "@/lib/constants";
import type { Asset } from "../models/asset.model";

export interface UseAssetImagesOptions {
  asset: Asset | null;
  assetId: string;
  onRefresh: () => Promise<void>;
}

export function useAssetImages({ asset, assetId, onRefresh }: UseAssetImagesOptions) {
  const [isUploading, setIsUploading] = useState(false);

  const imageCount = asset?.images?.length ?? 0;
  const canUpload = imageCount < MAX_IMAGES;

  const uploadImages = async (files: FileList | null) => {
    if (!files || files.length === 0 || !asset) return;

    const remainingSlots = MAX_IMAGES - asset.images.length;
    if (remainingSlots <= 0) {
      toast.error(`Maximum ${MAX_IMAGES} images allowed`);
      return;
    }

    const filesToUpload = Array.from(files).slice(0, remainingSlots);
    setIsUploading(true);

    try {
      for (const file of filesToUpload) {
        const urlResponse = await getImageUploadUrl(assetId, file.name);
        if (urlResponse.error) {
          throw new Error(urlResponse.error.message);
        }
        const uploadData = urlResponse.data;
        if (!uploadData) {
          throw new Error("Failed to get upload URL");
        }
        const { uploadUrl } = uploadData;

        await uploadImageToS3(uploadUrl, file);

        await confirmAssetUpload(assetId);
      }

      toast.success("Image uploaded successfully");
      await onRefresh();
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to upload image";
      toast.error(message);
    } finally {
      setIsUploading(false);
    }
  };

  return {
    isUploading,
    imageCount,
    canUpload,
    uploadImages,
  };
}