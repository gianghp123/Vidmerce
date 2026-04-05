import { Play, MoreVertical, Hourglass, Clock, RefreshCw, VideoOff } from "lucide-react";
import type { Video } from "../models/video.model";
import { StatusBadge } from "@/components/custom/StatusBadge";
import { MediaCard } from "@/components/custom/MediaCard";

interface VideoCardProps {
  video: Video;
  onClick?: () => void;
}

export function VideoCard({ video, onClick }: VideoCardProps) {
  const formatViews = (views: number) => {
    if (views >= 1000) {
      return `${(views / 1000).toFixed(1)}k Views`;
    }
    return `${views} Views`;
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const hours = Math.floor(diff / (1000 * 60 * 60));
    const days = Math.floor(hours / 24);

    if (hours < 1) return "Just now";
    if (hours < 24) return `Modified ${hours} hours ago`;
    if (days === 1) return "Modified yesterday";
    return `Modified ${days} days ago`;
  };

  const isFailed = video.status === "FAILED";
  const isProcessing = video.status === "PROCESSING";
  const isPending = video.status === "PENDING";

  return (
    <MediaCard
      onClick={onClick}
      aspectRatio="4/3"
      badge={<StatusBadge status={video.status} />}
      thumbnail={
        isFailed ? (
          <div className="w-full h-full flex flex-col items-center justify-center text-error/30 gap-2">
            <VideoOff className="w-14 h-14" />
            <p className="text-[10px] font-bold tracking-widest uppercase">Generation Failed</p>
          </div>
        ) : (
          <>
            <img
              src={video.thumbnail_url}
              alt={video.title}
              className={`w-full h-full object-cover group-hover:scale-105 transition-transform duration-500 ${isPending ? "opacity-60" : ""} ${isProcessing ? "grayscale-[0.5]" : ""}`}
            />
            <div className="absolute inset-0 bg-black/20 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
              <div className="w-12 h-12 bg-white/30 backdrop-blur-md rounded-full flex items-center justify-center text-white">
                <Play className="w-6 h-6 fill-current" />
              </div>
            </div>
          </>
        )
      }
    >
      <h5 className="font-headline font-bold text-on-surface mb-1 truncate">
        {video.title}
      </h5>
      
      {isFailed ? (
        <p className="text-xs text-on-surface-variant mb-4">Error Code: ERR_092</p>
      ) : isProcessing ? (
        <p className="text-xs text-on-surface-variant mb-4">Uploading: 67%</p>
      ) : isPending ? (
        <p className="text-xs text-on-surface-variant mb-4">Draft created yesterday</p>
      ) : (
        <p className="text-xs text-on-surface-variant mb-4">{formatDate(video.created_at)}</p>
      )}

      <div className="flex items-center justify-between">
        {isFailed ? (
          <span className="text-xs font-semibold text-error">Retry Required</span>
        ) : isProcessing ? (
          <span className="text-xs font-semibold text-outline-variant italic">Waiting for assets...</span>
        ) : isPending ? (
          <span className="text-xs font-semibold text-outline-variant">Not yet published</span>
        ) : (
          <span className="text-xs font-semibold text-primary">{formatViews(video.views || 0)}</span>
        )}

        {isFailed ? (
          <RefreshCw className="w-4 h-4 text-error cursor-pointer hover:rotate-180 transition-transform" />
        ) : isProcessing ? (
          <Hourglass className="w-4 h-4 text-outline-variant" />
        ) : isPending ? (
          <Clock className="w-4 h-4 text-outline-variant" />
        ) : (
          <MoreVertical className="w-4 h-4 text-outline-variant cursor-pointer hover:text-primary" />
        )}
      </div>
    </MediaCard>
  );
}
