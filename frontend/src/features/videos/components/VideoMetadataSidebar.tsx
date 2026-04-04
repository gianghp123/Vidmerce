import { Edit, ExternalLink } from "lucide-react";
import type { Video } from "../models/video.model";

interface VideoMetadataSidebarProps {
  video: Video;
  onEdit?: () => void;
}

export function VideoMetadataSidebar({ video, onEdit }: VideoMetadataSidebarProps) {
  const conversions = 1204;
  const ctr = 8.4;

  return (
    <div className="col-span-12 lg:col-span-4 flex flex-col gap-6">
      <div className="bg-surface-container-low p-8 rounded-2xl space-y-6 flex-1">
        <div className="space-y-2">
          <label className="text-[10px] font-bold text-outline-variant uppercase tracking-widest">
            Campaign Details
          </label>
          <p className="text-on-surface font-body leading-relaxed">
            {video.description || "This video integrates high-engagement hotspots directly linked to your Shopify inventory. Performance is tracking 12% above benchmark."}
          </p>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div className="bg-surface-container-lowest p-4 rounded-xl border border-outline-variant/10">
            <span className="text-[10px] font-bold text-secondary uppercase tracking-widest">
              Conversions
            </span>
            <p className="text-2xl font-headline font-extrabold text-on-surface">
              {conversions.toLocaleString()}
            </p>
          </div>
          <div className="bg-surface-container-lowest p-4 rounded-xl border border-outline-variant/10">
            <span className="text-[10px] font-bold text-secondary uppercase tracking-widest">
              CTR
            </span>
            <p className="text-2xl font-headline font-extrabold text-on-surface">
              {ctr}%
            </p>
          </div>
        </div>

        <div className="pt-4 space-y-3">
          <label className="text-[10px] font-bold text-outline-variant uppercase tracking-widest">
            Linked Products
          </label>
          {video.hotspots.slice(0, 2).map((hotspot, idx) => (
            <div
              key={idx}
              className="flex items-center justify-between p-3 bg-surface-container-highest/30 rounded-lg"
            >
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded bg-surface-container-high overflow-hidden">
                  <img
                    className="w-full h-full object-cover"
                    alt={hotspot.asset.name}
                    src={hotspot.asset.image_url}
                  />
                </div>
                <span className="text-xs font-medium text-on-surface line-clamp-1">
                  {hotspot.asset.name}
                </span>
              </div>
              <ExternalLink className="w-4 h-4 text-outline-variant cursor-pointer hover:text-primary" />
            </div>
          ))}
        </div>
      </div>

      <button
        onClick={onEdit}
        className="bg-surface-container-high text-primary py-4 rounded-2xl font-headline font-bold hover:bg-surface-container-highest transition-all flex items-center justify-center gap-2"
      >
        <Edit className="w-5 h-5" />
        Edit Interactive Layers
      </button>
    </div>
  );
}
