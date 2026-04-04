import { Badge } from "@/components/ui/badge";
import type { AssetStatus } from "../models/asset.model";

interface AssetStatusBadgeProps {
  status: AssetStatus;
}

export function AssetStatusBadge({ status }: AssetStatusBadgeProps) {
  const isPublished = status === "COMPLETED";
  const isDraft = status === "DRAFT";
  
  const className = isPublished
    ? "bg-primary-fixed text-on-primary-fixed-variant text-[10px] font-bold tracking-widest uppercase"
    : isDraft
      ? "bg-tertiary-container text-on-tertiary-container text-[10px] font-bold tracking-widest uppercase"
      : "text-[10px] font-bold tracking-widest uppercase";

  const label = isPublished ? "Published" : isDraft ? "Draft" : "Uploading";
  
  return (
    <Badge className={`px-3 py-1 rounded-full ${className}`}>
      {label}
    </Badge>
  );
}