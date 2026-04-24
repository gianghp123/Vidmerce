import { useCallback, useState } from "react";
import {
	confirmAssetUpload,
	createAssetGetUploadUrl,
	uploadImageToS3,
} from "../services/asset.action";
import type { CreateAssetReqDto } from "../dtos/req/create-asset.req.dto";
import type { IAsset } from "../../../lib/models/asset.model";

interface ImageFile {
	id: string;
	file: File;
	preview: string;
}

interface UseCreateAssetOptions {
	onSuccess?: (asset: IAsset) => void;
	onError?: (error: Error) => void;
}

export function useCreateAsset(options: UseCreateAssetOptions = {}) {
	const [isLoading, setIsLoading] = useState(false);
	const [error, setError] = useState<Error | null>(null);
	const [asset, setAsset] = useState<IAsset | null>(null);

	const createAsset = useCallback(
		async (payload: CreateAssetReqDto, images: ImageFile[]) => {
			setIsLoading(true);
			setError(null);

			try {
				// Step 1: Get upload URLs (only imageCount needed)
				const uploadUrlResponse = await createAssetGetUploadUrl({
					imageCount: payload.imageCount,
				});

				if (uploadUrlResponse.error || !uploadUrlResponse.data) {
					throw new Error(uploadUrlResponse.error?.message || "Failed to get upload URL");
				}

				const { assetId, uploads } = uploadUrlResponse.data;

				// Validate upload URLs count matches images count
				if (uploads.length !== images.length) {
					throw new Error(`Expected ${images.length} upload URLs, got ${uploads.length}`);
				}

				// Step 2: Upload all images in parallel (index-based pairing)
				const results = await Promise.allSettled(
					images.map((img, index) => {
						const upload = uploads[index];
						if (!upload) return Promise.reject(new Error(`Missing upload URL for image ${index}`));
						return uploadImageToS3(upload.uploadUrl, img.file);
					})
				);

				// Collect fileKeys from successful uploads
				const fileKeys: string[] = [];
				results.forEach((result, index) => {
					if (result.status === "fulfilled") {
						fileKeys.push(uploads[index].fileKey);
					}
				});

				// Lenient policy: proceed with partial success
				const failedCount = images.length - fileKeys.length;
				if (failedCount > 0) {
					console.warn(`${failedCount} of ${images.length} images failed; continuing with ${fileKeys.length} successful uploads`);
				}

				if (fileKeys.length === 0) {
					throw new Error("All image uploads failed");
				}

				// Step 3: Confirm asset creation with metadata and file keys
				const confirmResponse = await confirmAssetUpload(assetId, {
					name: payload.name,
					price: payload.price,
					productUrl: payload.productUrl,
					fileKeys,
				});

				if (confirmResponse.error || !confirmResponse.data) {
					throw new Error(confirmResponse.error?.message || "Failed to confirm asset upload");
				}

				const asset = confirmResponse.data;
				setAsset(asset);
				options.onSuccess?.(asset);
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
