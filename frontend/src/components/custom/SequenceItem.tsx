import { useState } from "react";
import { Image, Delete } from "lucide-react";
import { cn } from "@/lib/utils";
import type { SequenceItem, TransitionType } from "@/features/videos/models/video-sequence.model";

interface SequenceItemComponentProps {
  item: SequenceItem;
  onChange?: (item: SequenceItem) => void;
  onRemove?: (id: string) => void;
  assets?: { id: string; name: string; imageUrl: string }[];
  className?: string;
}

const transitionOptions: { value: TransitionType; label: string }[] = [
  { value: "FADE", label: "FADE" },
  { value: "SLIDE", label: "SLIDE" },
  { value: "ZOOM", label: "ZOOM" },
];

export function SequenceItemComponent({
  item,
  onChange,
  onRemove,
  assets = [],
  className,
}: SequenceItemComponentProps) {
  const [showAssetPicker, setShowAssetPicker] = useState(false);

  const selectedAsset = assets.find((a) => a.id === item.assetId);

  const handleAssetSelect = (assetId: string) => {
    const asset = assets.find((a) => a.id === assetId);
    onChange?.({
      ...item,
      assetId,
      assetName: asset?.name,
      assetImageUrl: asset?.imageUrl,
    });
    setShowAssetPicker(false);
  };

  return (
    <div
      className={cn(
        "bg-surface-container-lowest p-5 rounded-2xl flex items-center gap-6 shadow-sm group transition-all hover:shadow-md",
        className
      )}
    >
      <div
        className="w-24 h-24 rounded-xl overflow-hidden flex-shrink-0 bg-surface-container cursor-pointer"
        onClick={() => setShowAssetPicker(!showAssetPicker)}
      >
        {selectedAsset?.imageUrl ? (
          <img
            src={selectedAsset.imageUrl}
            alt={selectedAsset.name}
            className="w-full h-full object-cover"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center">
            <Image className="w-6 h-6 text-on-surface-variant" />
          </div>
        )}
      </div>

      <div className="flex-1 grid grid-cols-3 gap-6">
        <div className="space-y-1">
          <label className="text-[10px] font-bold text-outline uppercase tracking-wider">
            Asset
          </label>
          <select
            value={item.assetId}
            onChange={(e) => handleAssetSelect(e.target.value)}
            className="w-full border-none bg-surface-container-high rounded-lg text-sm font-medium py-1.5 focus:ring-0"
          >
            <option value="">Select asset...</option>
            {assets.map((asset) => (
              <option key={asset.id} value={asset.id}>
                {asset.name}
              </option>
            ))}
          </select>
        </div>

        <div className="space-y-1">
          <label className="text-[10px] font-bold text-outline uppercase tracking-wider">
            Duration (s)
          </label>
          <input
            type="number"
            value={item.duration}
            onChange={(e) =>
              onChange?.({ ...item, duration: parseFloat(e.target.value) || 0 })
            }
            className="w-full border-none bg-surface-container-high rounded-lg text-sm font-medium py-1.5 focus:ring-0"
            min={0.5}
            step={0.5}
          />
        </div>

        <div className="space-y-1">
          <label className="text-[10px] font-bold text-outline uppercase tracking-wider">
            Transition
          </label>
          <select
            value={item.transition}
            onChange={(e) =>
              onChange?.({ ...item, transition: e.target.value as TransitionType })
            }
            className="w-full border-none bg-surface-container-high rounded-lg text-sm font-medium py-1.5 focus:ring-0"
          >
            {transitionOptions.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </select>
        </div>
      </div>

      <button
        onClick={() => onRemove?.(item.id)}
        className="p-2 text-outline-variant hover:text-error transition-colors"
      >
        <Delete className="w-5 h-5" />
      </button>

      {showAssetPicker && assets.length > 0 && (
        <div className="absolute top-full left-0 mt-2 w-48 p-2 rounded-lg bg-surface-container-high border border-outline-variant shadow-lg z-10">
          {assets.map((asset) => (
            <button
              key={asset.id}
              onClick={() => handleAssetSelect(asset.id)}
              className="w-full flex items-center gap-2 p-2 rounded hover:bg-surface-container-low text-left"
            >
              <img
                src={asset.imageUrl}
                alt={asset.name}
                className="w-8 h-8 rounded object-cover"
              />
              <span className="text-sm truncate">{asset.name}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
