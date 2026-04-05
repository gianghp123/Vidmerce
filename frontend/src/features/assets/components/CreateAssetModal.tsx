import { useState } from "react";
import { X } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogClose,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { ImageUploadZone } from "@/components/custom/ImageUploadZone";
import { useCreateAsset } from "../hooks/useCreateAsset";
import type { CreateAssetPayload } from "../api/create-asset.api";

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

export function CreateAssetModal({
  open,
  onOpenChange,
  onSuccess,
}: CreateAssetModalProps) {
  const [formData, setFormData] = useState({
    name: "",
    price: "",
    productUrl: "",
  });
  const [images, setImages] = useState<ImageFile[]>([]);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const { createAsset, isLoading } = useCreateAsset({
    onSuccess: () => {
      setFormData({ name: "", price: "", productUrl: "" });
      setImages([]);
      onOpenChange(false);
      onSuccess?.();
    },
    onError: (error) => {
      setErrors({ submit: error.message });
    },
  });

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};
    if (!formData.name.trim()) newErrors.name = "Product name is required";
    if (!formData.price || parseFloat(formData.price) <= 0) newErrors.price = "Valid price is required";
    if (!formData.productUrl.trim()) newErrors.productUrl = "Product URL is required";
    if (images.length === 0) newErrors.images = "At least one image is required";

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;

    const payload: CreateAssetPayload = {
      name: formData.name.trim(),
      price: parseFloat(formData.price),
      productUrl: formData.productUrl.trim(),
    };

    await createAsset(payload, images);
  };

  const handleInputChange = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: "" }));
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange} >
      <DialogContent className="w-145! max-w-none! rounded-3xl bg-background border-none p-0 shadow-2xl" showCloseButton={false}>
        <DialogHeader className="px-10 pt-10 pb-6 relative">
          <div className="space-y-1">
            <DialogTitle className="text-3xl font-bold text-on-surface tracking-tight">
              New Asset
            </DialogTitle>
            <DialogDescription className="text-base text-on-surface-variant">
              Add a new product to your merchant library
            </DialogDescription>
          </div>
          <DialogClose className="absolute right-8 top-8 opacity-70 transition-opacity hover:opacity-100 outline-none">
            <X className="h-6 w-6 text-on-surface" />
          </DialogClose>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="px-10 pb-10 space-y-7">
          {/* Product Name */}
          <div className="space-y-2">
            <Label
              htmlFor="name"
              className="text-[11px] font-bold uppercase tracking-[0.15em] text-secondary"
            >
              Product Name
            </Label>
            <Input
              id="name"
              placeholder="e.g. Minimalist Linen Shirt"
              value={formData.name}
              onChange={(e) => handleInputChange("name", e.target.value)}
              className="h-12 bg-surface-container-lowest border-outline-variant/40 rounded-md px-4 text-on-surface placeholder:text-on-surface-variant/30 focus-visible:ring-1 focus-visible:ring-primary-container"
            />
          </div>

          {/* Price and URL Row */}
          <div className="grid grid-cols-2 gap-6">
            <div className="space-y-2">
              <Label
                htmlFor="price"
                className="text-[11px] font-bold uppercase tracking-[0.15em] text-secondary"
              >
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
                  value={formData.price}
                  onChange={(e) => handleInputChange("price", e.target.value)}
                  className="h-12 bg-surface-container-lowest border-outline-variant/40 rounded-md pl-8 pr-4 text-on-surface placeholder:text-on-surface-variant/30 focus-visible:ring-1 focus-visible:ring-primary-container"
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label
                htmlFor="productUrl"
                className="text-[11px] font-bold uppercase tracking-[0.15em] text-secondary"
              >
                Product URL
              </Label>
              <Input
                id="productUrl"
                placeholder="vidmerce.com/products/..."
                value={formData.productUrl}
                onChange={(e) => handleInputChange("productUrl", e.target.value)}
                className="h-12 bg-surface-container-lowest border-outline-variant/40 rounded-md px-4 text-on-surface placeholder:text-on-surface-variant/30 focus-visible:ring-1 focus-visible:ring-primary-container"
              />
            </div>
          </div>

          {/* Images Section */}
          <div className="space-y-3">
            <Label className="text-[11px] font-bold uppercase tracking-[0.15em] text-secondary">
              Images
            </Label>
            
            {/* The ImageUploadZone style is inherited from your tailwind theme */}
            <ImageUploadZone
              value={images}
              onChange={setImages}
              maxFiles={10}
              className="border-dashed border-2 border-outline-variant/50 rounded-xl bg-transparent min-h-[220px]"
            />
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
              disabled={isLoading}
              className="flex-[1.8] h-14 rounded-xl font-bold text-primary-foreground bg-primary-container hover:bg-primary-container/90 shadow-none transition-all active:scale-[0.98]"
            >
              {isLoading ? "Saving..." : "Save & Upload"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}