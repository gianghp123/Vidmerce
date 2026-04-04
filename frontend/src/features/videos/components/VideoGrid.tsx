import type { Video } from "../models/video.model";
import { VideoCard } from "./VideoCard";

interface VideoGridProps {
  videos: Video[];
  onVideoClick?: (video: Video) => void;
}

export function VideoGrid({ videos, onVideoClick }: VideoGridProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
      {videos.map((video) => (
        <VideoCard
          key={video.id}
          video={video}
          onClick={() => onVideoClick?.(video)}
        />
      ))}
    </div>
  );
}
