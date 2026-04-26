"use client";

import * as React from "react";
import { Upload, X } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  FileUpload,
  FileUploadDropzone,
  FileUploadItem,
  FileUploadItemDelete,
  FileUploadItemPreview,
  FileUploadList,
  FileUploadTrigger,
} from "@/components/ui/file-upload";
import { MAX_IMAGES } from "@/lib/constants";

interface ImageUploadZoneProps {
  value?: File[];
  onChange?: (files: File[]) => void;
  maxFiles?: number;
  className?: string;
}

export function ImageUploadZone({
  value = [],
  onChange,
  maxFiles = MAX_IMAGES,
  className,
}: ImageUploadZoneProps) {
  
  // Sync internal state with external onChange
  const handleValueChange = (newFiles: File[]) => {
    onChange?.(newFiles);
  };

  return (
    <FileUpload
      accept="image/*"
      maxFiles={maxFiles}
      maxSize={5 * 1024 * 1024} // 5MB limit
      value={value}
      onValueChange={handleValueChange}
      multiple
      className={cn("w-full space-y-3", className)}
    >
      <FileUploadDropzone className="rounded-xl bg-surface-container-lowest/50 py-12 px-6 hover:bg-surface-container-lowest transition-colors cursor-pointer">
        <div className="flex flex-col items-center justify-center text-center">
          <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center text-primary mb-4 group-hover:scale-110 transition-transform">
            <Upload className="w-6 h-6" />
          </div>
          <p className="text-sm font-semibold mb-1">
            Upload 1 to {maxFiles} images
          </p>
          <p className="text-xs text-muted-foreground">
            Drag and drop or click to browse
          </p>
          
          <FileUploadTrigger asChild>
            <Button variant="outline" size="sm" className="mt-4">
              Select Images
            </Button>
          </FileUploadTrigger>
        </div>
      </FileUploadDropzone>

      {value.length > 0 && (
        <FileUploadList orientation="horizontal" className="pb-2 gap-3">
          {value.map((file, index) => (
            <FileUploadItem
              key={`${file.name}-${index}`}
              value={file}
              className="relative w-20 h-20 shrink-0 p-0 border rounded-lg overflow-visible"
            >
              <FileUploadItemPreview 
                className="w-full h-full object-cover rounded-lg" 
              />
              
              {/* Optional: Add name overlay or tooltip if needed */}
              
              <FileUploadItemDelete asChild>
                <Button
                  variant="destructive"
                  size="icon"
                  className="absolute -top-2 -right-2 size-6 rounded-full shadow-sm"
                >
                  <X className="size-3" />
                </Button>
              </FileUploadItemDelete>
            </FileUploadItem>
          ))}
        </FileUploadList>
      )}
    </FileUpload>
  );
}

export default ImageUploadZone;