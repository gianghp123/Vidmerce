import { Badge } from "@/components/ui/badge";
import type { VideoStatus } from "@/features/videos/models/video.model";
import type { AssetStatus } from "@/features/assets/models/asset.model";

export type EntityStatus = VideoStatus | AssetStatus;

interface StatusConfig {
  className: string;
  label: string;
}

const defaultStatusConfig: Record<string, StatusConfig> = {
  COMPLETED: {
    className: "bg-green-100 text-green-800 text-[10px] font-extrabold px-2.5 py-1 rounded-full uppercase tracking-widest",
    label: "Completed",
  },
  PROCESSING: {
    className: "bg-amber-100 text-amber-800 text-[10px] font-extrabold px-2.5 py-1 rounded-full uppercase tracking-widest",
    label: "Processing",
  },
  PENDING: {
    className: "bg-tertiary-container text-on-tertiary-container text-[10px] font-extrabold px-2.5 py-1 rounded-full uppercase tracking-widest",
    label: "Pending",
  },
  FAILED: {
    className: "bg-error-container text-on-error-container text-[10px] font-extrabold px-2.5 py-1 rounded-full uppercase tracking-widest",
    label: "Failed",
  },
  DRAFT: {
    className: "bg-tertiary-container text-on-tertiary-container text-[10px] font-bold tracking-widest uppercase",
    label: "Draft",
  },
  UPLOADING: {
    className: "bg-muted text-muted-foreground text-[10px] font-bold tracking-widest uppercase",
    label: "Uploading",
  },
};

interface StatusBadgeProps {
  status: EntityStatus;
  labelMap?: Partial<Record<EntityStatus, string>>;
  className?: string;
}

export function StatusBadge({ status, labelMap, className }: StatusBadgeProps) {
  const config = defaultStatusConfig[status] || {
    className: "bg-muted text-muted-foreground text-[10px] font-extrabold px-2.5 py-1 rounded-full uppercase tracking-widest",
    label: status,
  };

  const label = labelMap?.[status] || config.label;

  return (
    <Badge className={className || config.className}>
      {label}
    </Badge>
  );
}
