import { ShoppingBag } from "lucide-react";
import type { Hotspot as HotspotType } from "../models/video.model";

interface HotspotProps {
  hotspot: HotspotType;
}

export function Hotspot({ hotspot }: HotspotProps) {
  const { x, y, asset } = hotspot;

  return (
    <div
      className="absolute z-20"
      style={{ left: `${x}%`, top: `${y}%`, transform: "translate(-50%, -50%)" }}
    >
      <div className="relative flex items-center justify-center group">
        <div className="w-8 h-8 bg-white/90 backdrop-blur-sm rounded-full hotspot-pulse cursor-pointer flex items-center justify-center text-primary-container">
          <ShoppingBag className="w-4 h-4" />
        </div>
        
        <div className="absolute top-10 left-0 w-64 bg-surface-container-lowest/90 backdrop-blur-xl rounded-2xl p-4 shadow-[0px_20px_40px_rgba(70,33,7,0.12)] opacity-0 group-hover:opacity-100 transition-all duration-500 translate-y-4 group-hover:translate-y-0">
          <div className="flex gap-4 mb-4">
            <div className="w-16 h-16 rounded-lg overflow-hidden bg-surface-container-low shrink-0">
              <img
                className="w-full h-full object-cover"
                alt={asset.name}
                src={asset.image_url}
              />
            </div>
            <div>
              <p className="text-xs font-bold text-primary tracking-wider uppercase mb-1">New Arrival</p>
              <h4 className="font-heading font-bold text-on-surface leading-tight line-clamp-2">
                {asset.name}
              </h4>
              <p className="text-sm font-semibold text-secondary-fixed-dim">
                ${asset.price.toFixed(2)}
              </p>
            </div>
          </div>
          <button className="w-full bg-primary text-on-primary py-2.5 rounded-xl font-heading font-bold text-xs hover:bg-primary-container transition-colors">
            Buy Now
          </button>
        </div>
      </div>
    </div>
  );
}
