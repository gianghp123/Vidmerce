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
import { useState } from "react";
import { toast } from "sonner";
import {
  confirmAssetUpload,
  createAssetGetUploadUrl,
  uploadImageToS3,
} from "../services/asset.action";

interface ImageFile {
  id: string;
  file: File;
  preview: string;
}

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
  const [images, setImages] = useState<ImageFile[]>([]);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const [isPreparing, setIsPreparing] = useState(false);


  const { dispatch: confirmAction, isPending: isConfirming } = useServerAction<
    IAsset,
    ConfirmPayload
  >(
    // Wrapper to adapt the existing action to the (prevState, payload) signature
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
        // Reset form and close modal ONLY on actual success
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

  const handleSubmit = async (e : React.SubmitEvent) => {
    e.preventDefault();
    if (!validate()) return;

    setIsPreparing(true);

    try {
      const uploadUrlResponse = await createAssetGetUploadUrl({
        imageCount: images.length,
      });

      // Assuming your action returns standard BaseResponse wrapper
      if (uploadUrlResponse.error || !uploadUrlResponse.data) {
        throw new Error(uploadUrlResponse.error?.message || "Failed to get upload URLs");
      }

      const { assetId, uploads } = uploadUrlResponse.data;

      // Step 2: Upload images to S3
      const results = await Promise.allSettled(
        images.map((img, index) => {
          const upload = uploads[index];
          if (!upload)
            return Promise.reject(new Error(`Missing upload URL for image ${index}`));
          return uploadImageToS3(upload.uploadUrl, img.file);
        })
      );

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

      // Step 3: Pass off to the Next.js Server Action Hook for final confirmation
      // This will automatically toggle `isConfirming` to true, and handle success/error toasts
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
        className="w-145! max-w-none! rounded-3xl bg-background border-none p-0 shadow-2xl"
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

        <form onSubmit={handleSubmit} className="px-10 pb-10 space-y-7">
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
              className="border-dashed border-2 border-outline-variant/50 rounded-xl bg-transparent min-h-55"
            />
            {errors.images && <p className="text-sm text-destructive">{errors.images}</p>}
          </div>

          {/* Footer Buttons */}
          <div className="pt-4 flex gap-4">
            <Button
              type="button"
              variant="ghost"
              onClick={() => onOpenChange(false)}
              className="flex-1 h-14 rounded-xl font-bold text-on-surface-variant bg-surface-container-high hover:bg-surface-container-highest transition-colors"
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={isSubmitting}
              variant={"default"}
            >
              {isSubmitting ? "Saving..." : "Save & Upload"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}