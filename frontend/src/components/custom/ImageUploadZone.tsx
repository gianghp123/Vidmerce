import { useCallback, useState } from "react";
import { Upload, X } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";

interface ImageFile {
  id: string;
  file: File;
  preview: string;
}

interface ImageUploadZoneProps {
  value?: ImageFile[];
  onChange?: (files: ImageFile[]) => void;
  maxFiles?: number;
  className?: string;
}

export function ImageUploadZone({
  value = [],
  onChange,
  maxFiles = 10,
  className,
}: ImageUploadZoneProps) {
  const [isDragging, setIsDragging] = useState(false);

  const addFiles = useCallback((files: File[]) => {
    const remainingSlots = maxFiles - value.length;
    const filesToAdd = files.slice(0, remainingSlots);
    const newFiles: ImageFile[] = filesToAdd.map((file) => ({
      id: Math.random().toString(36).substring(7),
      file,
      preview: URL.createObjectURL(file),
    }));
    onChange?.([...value, ...newFiles]);
  }, [value, maxFiles, onChange]);

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      setIsDragging(false);
      const files = Array.from(e.dataTransfer.files).filter((file) =>
        file.type.startsWith("image/")
      );
      addFiles(files);
    },
    [addFiles]
  );

  const handleFileInput = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      if (e.target.files) {
        const files = Array.from(e.target.files);
        addFiles(files);
      }
    },
    [addFiles]
  );

  const removeFile = (id: string) => {
    const newFiles = value.filter((f) => f.id !== id);
    onChange?.(newFiles);
  };

  return (
    <div className={cn("space-y-3", className)}>
      <div
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        className={cn(
          "group relative flex flex-col items-center justify-center border-2 border-dashed border-outline-variant/40 rounded-xl bg-surface-container-lowest/50 py-12 px-6 hover:bg-surface-container-lowest transition-colors cursor-pointer",
          isDragging && "bg-surface-container-lowest"
        )}
      >
        <input
          type="file"
          accept="image/*"
          multiple
          onChange={handleFileInput}
          className="absolute inset-0 opacity-0 cursor-pointer"
          id="image-upload"
          disabled={value.length >= maxFiles}
        />
        <div className="w-16 h-16 rounded-full bg-primary-fixed flex items-center justify-center text-primary mb-4 group-hover:scale-110 transition-transform">
          <Upload className="w-6 h-6" />
        </div>
        <p className="text-on-surface font-semibold mb-1">Upload 1 to 10 images</p>
        <p className="text-sm text-on-surface-variant">Drag and drop or click to browse</p>
      </div>

      {value.length > 0 && (
        <div className="flex gap-3 overflow-x-auto pb-2">
          {value.map((image) => (
            <div
              key={image.id}
              className="w-16 h-16 rounded-lg bg-surface-container-high border border-outline-variant/20 flex items-center justify-center shrink-0 relative group"
            >
              <img
                src={image.preview}
                alt="Preview"
                className="w-full h-full object-cover rounded-lg"
              />
              <Button
                variant="secondary"
                size="icon-sm"
                className="absolute -top-2 -right-2 opacity-0 group-hover:opacity-100 transition-opacity w-6 h-6"
                onClick={() => removeFile(image.id)}
              >
                <X className="w-3 h-3" />
              </Button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
