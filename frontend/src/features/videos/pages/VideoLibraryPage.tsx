import { LibraryLayout } from "@/components/custom/LibraryLayout";
import { CursorPagination } from "@/components/custom/Pagination";
import { Button } from "@/components/ui/button";
import { Grid3X3, List, Plus } from "lucide-react";
import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { VideoGrid } from "../components/VideoGrid";
import { VideoMetadataSidebar } from "../components/VideoMetadataSidebar";
import { VideoPlayer } from "../components/VideoPlayer";
import { useVideos } from "../hooks/useVideos";
import type { Video } from "../models/video.model";
import { EmptyState } from "@/components/custom/EmptyState";

export function VideoLibraryPage() {
  const navigate = useNavigate();
  const [selectedVideo, setSelectedVideo] = useState<Video | null>(null);
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");
  const limit = 12;

  const { videos, isLoading, hasMore, fetchNext, error } = useVideos({ limit });

  useEffect(() => {
    if (error) {
      toast.error(error.message);
    }
  }, [error]);

  const handleVideoClick = (video: Video) => {
    setSelectedVideo(video);
  };

  const handleEditLayers = () => {
    console.log("Edit interactive layers for:", selectedVideo?.title);
  };

  const displayVideo = selectedVideo ?? videos[0] ?? null;
  const isEmpty = !isLoading && videos.length === 0;

  return (
    <>
      {/* Active Video Section */}
      {displayVideo && (
        <section className="grid grid-cols-12 gap-8">
          <div className="col-span-12">
            <div className="flex items-baseline gap-4 mb-2">
              <h2 className="font-headline text-4xl font-extrabold text-on-surface tracking-tight">
                {displayVideo.title}
              </h2>
              {displayVideo.status === "COMPLETED" && (
                <span className="bg-primary-fixed text-on-primary-fixed-variant px-3 py-1 rounded-full text-xs font-bold uppercase tracking-widest">
                  Live Now
                </span>
              )}
            </div>
            <p className="text-on-surface-variant font-body text-lg">
              {displayVideo.description || "Shoppable video campaign for premium products."}
            </p>
          </div>

          <div className="col-span-12 lg:col-span-8">
            <VideoPlayer video={displayVideo} />
          </div>

          <div className="col-span-12 lg:col-span-4">
            <VideoMetadataSidebar video={displayVideo} onEdit={handleEditLayers} />
          </div>
        </section>
      )}

      <LibraryLayout
        title="Video Library"
        description="Browse and manage your generated shoppable videos."
        actions={
          <div className="flex items-center gap-2">
            <Button
              onClick={() => navigate("/videos/create")}
              className="gap-2 bg-linear-to-br from-secondary to-primary-container text-primary-foreground rounded-xl font-heading font-bold text-sm shadow-lg active:opacity-80 transition-opacity"
            >
              <Plus className="w-4 h-4 mr-2" />
              Create Video
            </Button>

            <button
              onClick={() => setViewMode("grid")}
              className={`p-2 rounded-lg ${viewMode === "grid"
                ? "bg-surface-container-high text-on-surface-variant"
                : "text-outline-variant hover:text-primary"
                } transition-colors`}
            >
              <Grid3X3 className="w-5 h-5" />
            </button>

            <button
              onClick={() => setViewMode("list")}
              className={`p-2 rounded-lg ${viewMode === "list"
                ? "bg-surface-container-high text-on-surface-variant"
                : "text-outline-variant hover:text-primary"
                } transition-colors`}
            >
              <List className="w-5 h-5" />
            </button>
          </div>
        }
      >
        {/* Content */}
        <div className="pt-6">
          {isLoading && videos.length === 0 ? (
            <div className="flex items-center justify-center py-20">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
          ) : isEmpty ? (
            <EmptyState
              title="No videos yet"
              description="Start by creating your first shoppable video."
            />
          ) : (
            <>
              <VideoGrid videos={videos} onVideoClick={handleVideoClick} />

              <CursorPagination
                hasMore={hasMore}
                onLoadMore={fetchNext}
                isLoading={isLoading}
                entityName="videos"
              />
            </>
          )}
        </div>
      </LibraryLayout>


    </>
  );
}