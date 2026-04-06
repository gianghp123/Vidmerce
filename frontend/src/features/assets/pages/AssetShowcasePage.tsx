import { useState, useEffect, useCallback } from "react";
import { useParams } from "react-router-dom";
import { toast } from "sonner";
import { Upload, Trash2, ArrowLeft } from "lucide-react";
import { LibraryLayout } from "@/components/custom/LibraryLayout";
import { NotFound } from "@/components/custom/NotFound";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { fetchAsset, deleteAssetImage } from "../api/asset.api";
import { confirmAssetUpload } from "../api/create-asset.api";
import type { Asset } from "../models/asset.model";
import { useAssetImages } from "../hooks/useAssetImages";
import { MAX_IMAGES } from "@/lib/constants";

export function AssetShowcasePage() {
  const { id } = useParams<{ id: string }>();
  const [asset, setAsset] = useState<Asset | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchAssetData = useCallback(async () => {
    if (!id) return;
    try {
      const response = await fetchAsset(id);
      if (response.error) {
        toast.error(response.error.message);
        setError(response.error.message);
      } else {
        setAsset(response.data);
        setError(null);
      }
    } catch {
      setError("Failed to load asset");
    } finally {
      setIsLoading(false);
    }
  }, [id]);

  useEffect(() => {
    fetchAssetData();
  }, [fetchAssetData]);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    uploadImages(e.target.files);
    e.target.value = "";
  };

  const handleDeleteImage = async (imageId: string) => {
    if (!id) return;

    try {
      const deleteResponse = await deleteAssetImage(id, imageId);
      if (deleteResponse.error) {
        throw new Error(deleteResponse.error.message);
      }

      await confirmAssetUpload(id);
      toast.success("Image deleted successfully");
      await fetchAssetData();
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to delete image";
      toast.error(message);
    }
  };

  const {
    isUploading,
    canUpload: canUploadFromHook,
    uploadImages,
  } = useAssetImages({
    asset,
    assetId: id || "",
    onRefresh: fetchAssetData,
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
      </div>
    );
  }

  if (error || !asset) {
    return (
      <LibraryLayout
        title="Error"
        actions={
          <Button onClick={() => window.history.back()}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back
          </Button>
        }
      >
        <NotFound
          title={error ? "Error" : "Asset Not Found"}
          description={error || "The asset you're looking for doesn't exist."}
          actionLabel="Go Back"
          onAction={() => window.history.back()}
        />
      </LibraryLayout>
    );
  }

  return (
    <LibraryLayout
      title={asset.name}
      description={asset.description || "Asset details"}
      actions={
        <Button onClick={() => window.history.back()}>
          <ArrowLeft className="w-4 h-4 mr-2" />
          Back
        </Button>
      }
    >
      <div className="space-y-8 pt-6">
        <Card>
          <CardContent className="p-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <h3 className="text-sm font-medium text-on-surface-variant mb-1">Price</h3>
                <p className="text-2xl font-bold">${asset.price.toFixed(2)}</p>
              </div>
              <div>
                <h3 className="text-sm font-medium text-on-surface-variant mb-1">Status</h3>
                <p className="text-lg font-medium">{asset.status}</p>
              </div>
              {asset.productUrl && (
                <div className="md:col-span-2">
                  <h3 className="text-sm font-medium text-on-surface-variant mb-1">Product URL</h3>
                  <a
                    href={asset.productUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-primary hover:underline"
                  >
                    {asset.productUrl}
                  </a>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        <div>
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold">Images ({asset.images.length}/{MAX_IMAGES})</h2>
            {canUploadFromHook && (
              <label>
                <input
                  type="file"
                  accept="image/*"
                  multiple
                  onChange={handleFileChange}
                  className="hidden"
                  disabled={isUploading}
                />
                <Button asChild disabled={isUploading}>
                  <span className="cursor-pointer gap-2">
                    <Upload className="w-4 h-4" />
                    {isUploading ? "Uploading..." : "Upload Images"}
                  </span>
                </Button>
              </label>
            )}
          </div>

          {asset.images.length === 0 ? (
            <div className="text-center py-10 text-on-surface-variant">
              No images yet. Upload up to {MAX_IMAGES} images.
            </div>
          ) : (
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
              {asset.images.map((image, index) => (
                <div
                  key={index}
                  className="relative aspect-square rounded-lg overflow-hidden bg-surface-container-high group"
                >
                  <img
                    src={image.imageUrl}
                    alt={`Asset image ${index + 1}`}
                    className="w-full h-full object-cover"
                  />
                  <button
                    onClick={() => handleDeleteImage(image.order.toString())}
                    className="absolute top-2 right-2 p-2 rounded-full bg-error text-on-error opacity-0 group-hover:opacity-100 transition-opacity hover:bg-error/80"
                  >
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </LibraryLayout>
  );
}