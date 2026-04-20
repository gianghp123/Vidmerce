import { MediaCard } from "@/components/custom/MediaCard";
import { StatusBadge } from "@/components/custom/StatusBadge";
import { ArrowRight, MoreVertical } from "lucide-react";
import { useNavigate } from "react-router-dom";
import type { IAsset } from "../../../lib/models/asset.model";

interface AssetCardProps {
  asset: IAsset;
}

export function AssetCard({ asset }: AssetCardProps) {
  const navigate = useNavigate();
  const thumbnailUrl = asset.images?.[0]?.imageUrl || "/placeholder-image.jpg";
  const formattedPrice = new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  }).format(asset.price);

  return (
    <MediaCard
      aspectRatio="4/5"
      badge={<StatusBadge status={asset.status} />}
      thumbnail={
        <img
          src={thumbnailUrl}
          alt={asset.name}
          className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110"
        />
      }
      onClick={() => navigate(`/assets/${asset.id}`)}
    >
      <div className="flex justify-between items-start mb-2">
        <h3 className="font-headline font-bold text-lg text-on-surface line-clamp-1">
          {asset.name}
        </h3>
        <button className="text-outline cursor-pointer hover:text-primary transition-colors p-1">
          <MoreVertical className="w-5 h-5" />
        </button>
      </div>
      <p className="text-primary font-bold text-lg mb-4">{formattedPrice}</p>
      <a
        href={asset.productUrl}
        target="_blank"
        rel="noopener noreferrer"
        className="flex items-center gap-2 text-sm font-semibold text-secondary hover:underline group/link"
      >
        View Product
        <ArrowRight className="w-4 h-4 group-hover/link:translate-x-1 transition-transform" />
      </a>
    </MediaCard>
  );
}
