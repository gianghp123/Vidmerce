"use client";

import { ImageUploadZone } from "@/components/common/ImageUploadZone";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { MAX_IMAGES } from "@/lib/constants";
import { useServerAction } from "@/lib/hooks/use-server-action";
import type { IAsset } from "@/types/models/asset.model";
import { X } from "lucide-react";
import React, { useState } from "react";
import { toast } from "sonner";
import {
  confirmAssetUpload,
  createAssetGetUploadUrl,
  uploadImageToS3,
} from "../services/asset.action";

interface CreateAssetModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

interface ConfirmPayload {
  assetId: string;
  name: string;
  price: number;
  productUrl: string;
  fileKeys: string[];
}

export function CreateAssetModal({
  open,
  onOpenChange,
  onSuccess,
}: CreateAssetModalProps) {
  const [formData, setFormData] = useState<Partial<IAsset>>({
    name: "",
    price: undefined,
    productUrl: "",
  });
  
  // Updated to use File[] to match the new ImageUploadZone
  const [images, setImages] = useState<File[]>([]);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const [isPreparing, setIsPreparing] = useState(false);

  const { dispatch: confirmAction, isPending: isConfirming } = useServerAction<
    IAsset,
    ConfirmPayload
  >(
    async (prevState, payload) => {
      return confirmAssetUpload(payload.assetId, {
        name: payload.name,
        price: payload.price,
        productUrl: payload.productUrl,
        fileKeys: payload.fileKeys,
      });
    },
    {
      successMessage: "Asset created successfully",
      errorMessage: "Failed to confirm asset upload",
      onSuccess: () => {
        setFormData({ name: "", price: undefined, productUrl: "" });
        setImages([]);
        onOpenChange(false);
        onSuccess?.();
      },
    }
  );

  const isSubmitting = isPreparing || isConfirming;

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};
    if (!formData.name?.trim()) newErrors.name = "Product name is required";
    if (!formData.price || formData.price <= 0)
      newErrors.price = "Valid price is required";
    if (!formData.productUrl?.trim())
      newErrors.productUrl = "Product URL is required";
    if (images.length === 0)
      newErrors.images = "At least one image is required";

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.SubmitEvent) => {
    e.preventDefault();
    if (!validate()) return;

    setIsPreparing(true);

    try {
      // Step 1: Get Pre-signed URLs
      const uploadUrlResponse = await createAssetGetUploadUrl({
        imageCount: images.length,
      });

      if (uploadUrlResponse.error || !uploadUrlResponse.data) {
        throw new Error(uploadUrlResponse.error?.message || "Failed to get upload URLs");
      }

      const { assetId, uploads } = uploadUrlResponse.data;

      // Step 2: Upload images to S3
      const results = await Promise.allSettled(
        images.map((file, index) => {
          const upload = uploads[index];
          if (!upload)
            return Promise.reject(new Error(`Missing upload URL for image ${index}`));
          return uploadImageToS3(upload.uploadUrl, file);
        })
      );

      if (results.some((result) => result.status === "rejected")) {
        throw new Error("Images upload failed, please try again");
      }

      const fileKeys: string[] = [];
      results.forEach((result, index) => {
        if (result.status === "fulfilled") {
          fileKeys.push(uploads[index].fileKey);
        }
      });

      if (fileKeys.length === 0) {
        throw new Error("All image uploads failed");
      }

      if (fileKeys.length < images.length) {
        toast.warning(
          `${images.length - fileKeys.length} of ${images.length} images failed to upload`
        );
      }

      // Step 3: Final confirmation via Server Action Hook
      confirmAction({
        assetId,
        name: formData.name!,
        price: formData.price!,
        productUrl: formData.productUrl!,
        fileKeys,
      });

    } catch (err) {
      const error = err instanceof Error ? err : new Error("Failed to create asset");
      toast.error(error.message);
      setErrors({ submit: error.message });
    } finally {
      setIsPreparing(false);
    }
  };

  const handleInputChange = (field: keyof IAsset, value: string) => {
    setFormData((prev) => ({
      ...prev,
      [field]: field === "price" ? parseFloat(value) || undefined : value,
    }));
    if (errors[field as string]) {
      setErrors((prev) => ({ ...prev, [field]: "" }));
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        className="max-w-none! w-fit rounded-3xl bg-background border-none p-0 shadow-2xl"
        showCloseButton={false}
      >
        <DialogHeader className="px-10 pt-10 pb-6 relative">
          <div className="space-y-1">
            <DialogTitle className="text-3xl font-bold text-on-surface tracking-tight">
              New Asset
            </DialogTitle>
            <DialogDescription className="text-base text-on-surface-variant">
              Add a new product to your library
            </DialogDescription>
          </div>
          <DialogClose className="absolute right-8 top-8 opacity-70 transition-opacity hover:opacity-100 outline-none">
            <X className="h-6 w-6 text-on-surface" />
          </DialogClose>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="px-10 pb-10 space-y-7 md:min-w-145">
          {/* Product Name */}
          <div className="space-y-2">
            <Label htmlFor="name" className="text-[11px] font-bold uppercase tracking-[0.15em] text-secondary">
              Product Name
            </Label>
            <Input
              id="name"
              placeholder="e.g. Minimalist Linen Shirt"
              value={formData.name || ""}
              onChange={(e) => handleInputChange("name", e.target.value)}
              className="h-12 bg-surface-container-lowest border-outline-variant/40 rounded-md px-4 text-on-surface placeholder:text-on-surface-variant/30 focus-visible:ring-1 focus-visible:ring-primary-container"
            />
            {errors.name && <p className="text-sm text-destructive">{errors.name}</p>}
          </div>

          {/* Price and URL Row */}
          <div className="grid grid-cols-2 gap-6">
            <div className="space-y-2">
              <Label htmlFor="price" className="text-[11px] font-bold uppercase tracking-[0.15em] text-secondary">
                Price
              </Label>
              <div className="relative">
                <span className="absolute left-4 top-1/2 -translate-y-1/2 text-on-surface-variant text-sm">
                  $
                </span>
                <Input
                  id="price"
                  type="number"
                  step="0.01"
                  placeholder="0.00"
                  value={formData.price || ""}
                  onChange={(e) => handleInputChange("price", e.target.value)}
                  className="h-12 bg-surface-container-lowest border-outline-variant/40 rounded-md pl-8 pr-4 text-on-surface placeholder:text-on-surface-variant/30 focus-visible:ring-1 focus-visible:ring-primary-container"
                />
              </div>
              {errors.price && <p className="text-sm text-destructive">{errors.price}</p>}
            </div>

            <div className="space-y-2">
              <Label htmlFor="productUrl" className="text-[11px] font-bold uppercase tracking-[0.15em] text-secondary">
                Product URL
              </Label>
              <Input
                id="productUrl"
                placeholder="vidmerce.com/products/..."
                value={formData.productUrl || ""}
                onChange={(e) => handleInputChange("productUrl", e.target.value)}
                className="h-12 bg-surface-container-lowest border-outline-variant/40 rounded-md px-4 text-on-surface placeholder:text-on-surface-variant/30 focus-visible:ring-1 focus-visible:ring-primary-container"
              />
              {errors.productUrl && <p className="text-sm text-destructive">{errors.productUrl}</p>}
            </div>
          </div>

          {/* Images Section */}
          <div className="space-y-3">
            <Label className="text-[11px] font-bold uppercase tracking-[0.15em] text-secondary">
              Images
            </Label>
            <ImageUploadZone
              value={images}
              onChange={setImages}
              maxFiles={MAX_IMAGES}
            />
            {errors.images && <p className="text-sm text-destructive">{errors.images}</p>}
          </div>

          {/* Footer Buttons */}
          <div className="pt-4 flex gap-4 justify-end">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isSubmitting}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={isSubmitting}
            >
              {isSubmitting ? "Saving..." : "Save & Upload"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}