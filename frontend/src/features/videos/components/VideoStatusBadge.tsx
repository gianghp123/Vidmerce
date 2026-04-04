import { Badge } from "@/components/ui/badge";
import type { VideoStatus } from "../models/video.model";

interface VideoStatusBadgeProps {
  status: VideoStatus;
}

export function VideoStatusBadge({ status }: VideoStatusBadgeProps) {
  const getStatusConfig = () => {
    switch (status) {
      case "COMPLETED":
        return {
          className: "bg-green-100 text-green-800 text-[10px] font-extrabold px-2.5 py-1 rounded-full uppercase tracking-widest",
          label: "Completed",
        };
      case "PROCESSING":
        return {
          className: "bg-amber-100 text-amber-800 text-[10px] font-extrabold px-2.5 py-1 rounded-full uppercase tracking-widest",
          label: "Processing",
        };
      case "PENDING":
        return {
          className: "bg-tertiary-container text-on-tertiary-container text-[10px] font-extrabold px-2.5 py-1 rounded-full uppercase tracking-widest",
          label: "Pending",
        };
      case "FAILED":
        return {
          className: "bg-error-container text-on-error-container text-[10px] font-extrabold px-2.5 py-1 rounded-full uppercase tracking-widest",
          label: "Failed",
        };
      default:
        return {
          className: "bg-muted text-muted-foreground text-[10px] font-extrabold px-2.5 py-1 rounded-full uppercase tracking-widest",
          label: status,
        };
    }
  };

  const { className, label } = getStatusConfig();

  return <Badge className={className}>{label}</Badge>;
}
