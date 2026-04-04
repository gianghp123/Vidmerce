import type { Asset } from "../models/asset.model";
import { AssetCard } from "./AssetCard";

interface AssetGridProps {
  assets: Asset[];
}

export function AssetGrid({ assets }: AssetGridProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
      {assets.map((asset) => (
        <AssetCard key={asset.assetId} asset={asset} />
      ))}
    </div>
  );
}