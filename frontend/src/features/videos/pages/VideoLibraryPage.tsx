import { useState } from "react";
import { Grid3X3, List } from "lucide-react";
import { MainLayout } from "@/components/layout/MainLayout";
import { TopBar } from "@/components/layout/TopBar";
import { VideoPlayer } from "../components/VideoPlayer";
import { VideoMetadataSidebar } from "../components/VideoMetadataSidebar";
import { VideoGrid } from "../components/VideoGrid";
import { useVideos } from "../hooks/useVideos";
import type { Video } from "../models/video.model";

const mockVideos: Video[] = [
  {
    id: "1",
    title: "Autumn Collection Showcase",
    description: "Shoppable video campaign for Nordic lifestyle furniture.",
    status: "COMPLETED",
    video_url: "https://example.com/video1.mp4",
    thumbnail_url: "https://lh3.googleusercontent.com/aida-public/AB6AXuADsWg5pdyhkA0gIczJHQ5pPetMqnzonjxxnPybAwfLS5vO8xunMPw1zvajyVWloEiYcQNX2QIew0A0urgJnybfRhH4uowZto1ns9W3zVQYMiI519N6AzzTKgD0Fz88sVrhbTpXIlKpMvoij96U_HKmr0TRWTOYY5pTzUbKPaijljPekPwngj-hHym-I-NF6cXHk8xfjF1_EhP1OH4uDD7s5rIAG0AgJdyJj2AGDHF_C2gBNs-LMf8I0TJLoPLnXCsE7Ed6a2xQhbWt",
    created_at: "2024-01-15T10:00:00Z",
    updated_at: "2024-01-15T10:00:00Z",
    views: 12500,
    hotspots: [
      {
        x: 50,
        y: 50,
        asset: {
          asset_id: "a1",
          name: "Nordic Elm Dining Chair",
          price: 349,
          image_url: "https://lh3.googleusercontent.com/aida-public/AB6AXuB-HUHJDyDbRe-SM1OiPCt9bHYpOE_nY15HRwFY9iMogZ73zL5cjzYtXnzmC6jUDnxicgkI_ng7TCOkoyEnCL-5uCsFxdjIws4KiYIOqL11G_J2OW-XmYe0OHIrfqgs7w6Cb-EmWiJhgE7jB9zuXfY_CjZJgwoDZ-yEq-dVy5atEv1xuq_iXI-s7Ba7bHqR5hEaqZSI_kMiOjrJN74uJiB9x89SGMIW3W0sTcBKG3WLEOjPs8Sk1Tx8u30AR7QiUVjMLnjF_4Ul7vW0",
          product_url: "#",
        },
      },
    ],
  },
  {
    id: "2",
    title: "Spring Minimalist Vibes",
    status: "COMPLETED",
    video_url: "https://example.com/video2.mp4",
    thumbnail_url: "https://lh3.googleusercontent.com/aida-public/AB6AXuA6kUPrJpjspvqBDmrt08aNBBVP_ak_XBCEj-cRukw6KKx2Ia6Cbj9xkuPORJjZazp3ZW2hrJALHD8UBOnprgbSgzH-J9_lYuQ2gxO8_t8zbPVfusqO_d7yPIZuevCJpxvd37g2t4UcwOcO-gqyVVezJlL1wnP80YmhHAgy5YRCRHyx6VUB7sHl1igJn3sDikT9E0NRDDD8qdgebydhrWQeIDuUuSUGz64RZWtCvIxhs-mWHVZ9I1zgtXmujKj-Gr3XRS8bN1eD3DIp",
    created_at: "2024-01-14T09:30:00Z",
    views: 8900,
    hotspots: [],
  },
  {
    id: "3",
    title: "Home Office Essentials",
    status: "PROCESSING",
    video_url: "",
    thumbnail_url: "https://lh3.googleusercontent.com/aida-public/AB6AXuAAbRnag5jvhWnxo8BWRMfEt3wyv2k2jVEyL_p93Pi3lIrZw5kjXNNwQ4A_FbnCx-TX8ljNkHKgKV6u_tUnnpRNgoJh1XtaIyMjPxyQfZZ6RvXc1Z66bk9zk_-_24TZ-_k6eCu6ogGhZkIhDtrlzVIoWzv0HMjFuLdmoa9xiJvV2Mwhczq-Hz7wQerGdOTgQh1vPVKJvmS4_FSIdV3YftFHslrudTxbGQTfyrB2A8GVw7w6rrGpBpOTY_DPjlawWZ1zF52-d2egTHIh",
    created_at: "2024-01-13T14:20:00Z",
    hotspots: [],
  },
  {
    id: "4",
    title: "Kitchenware Collection",
    status: "PENDING",
    video_url: "",
    thumbnail_url: "https://lh3.googleusercontent.com/aida-public/AB6AXuB58VZawlj7OukaEGsgWQP8f7JSInANlwXnA52ImQm1nC13GZZlcFOcd1FbYYTiQ8orce0H24vJji68fiy6pMcoQskM5MZj86D6AFeOD3kqWm15Nao45ubTiEbijAgHQUNNIq2Thk275OFfUIWCbFCkaPMTS1MqZApM7jWGXpK6AtQZnqZdTIhKu_FtxTgNEX7HE4xQ3-MXXFdrnp4aSAuz-QbwwHTPm7G_mkOHvTvqZVsy685lmkMlgPe3afzczyvR8JleKDbwVNwd",
    created_at: "2024-01-12T11:45:00Z",
    hotspots: [],
  },
  {
    id: "5",
    title: "Holiday Gift Guide",
    status: "FAILED",
    video_url: "",
    thumbnail_url: "",
    created_at: "2024-01-11T08:00:00Z",
    hotspots: [],
  },
];

export function VideoLibraryPage() {
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedVideo, setSelectedVideo] = useState<Video>(mockVideos[0]);
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");
  const limit = 12;

  const { videos, total, totalPages, isLoading } = useVideos({
    page: currentPage,
    limit,
  });

  const displayVideos = videos.length > 0 ? videos : mockVideos;
  const displayTotal = total || mockVideos.length;
  const displayTotalPages = totalPages || Math.ceil(mockVideos.length / limit);

  const handleVideoClick = (video: Video) => {
    setSelectedVideo(video);
  };

  const handleEditLayers = () => {
    console.log("Edit interactive layers for:", selectedVideo.title);
  };

  return (
    <MainLayout>
      <TopBar />
      
      <main className="p-12 space-y-16 max-w-7xl mx-auto">
        {/* Active Video Section */}
        <section className="grid grid-cols-12 gap-8">
          <div className="col-span-12">
            <div className="flex items-baseline gap-4 mb-2">
              <h2 className="font-headline text-4xl font-extrabold text-on-surface tracking-tight">
                {selectedVideo.title}
              </h2>
              {selectedVideo.status === "COMPLETED" && (
                <span className="bg-primary-fixed text-on-primary-fixed-variant px-3 py-1 rounded-full text-xs font-bold uppercase tracking-widest">
                  Live Now
                </span>
              )}
            </div>
            <p className="text-on-surface-variant font-body text-lg">
              {selectedVideo.description || "Shoppable video campaign for premium products."}
            </p>
          </div>

          <div className="col-span-12 lg:col-span-8">
            <VideoPlayer video={selectedVideo} />
          </div>

          <div className="col-span-12 lg:col-span-4">
            <VideoMetadataSidebar video={selectedVideo} onEdit={handleEditLayers} />
          </div>
        </section>

        {/* Video Library Gallery */}
        <section className="space-y-8">
          <div className="flex justify-between items-end border-b border-outline-variant/10 pb-6">
            <div>
              <h3 className="font-headline text-2xl font-bold text-on-surface">
                Recent Projects
              </h3>
              <p className="text-on-surface-variant font-body">
                Browse your recently generated shoppable content.
              </p>
            </div>
            <div className="flex gap-3">
              <button
                onClick={() => setViewMode("grid")}
                className={`p-2 rounded-lg ${
                  viewMode === "grid"
                    ? "bg-surface-container-high text-on-surface-variant"
                    : "text-outline-variant hover:text-primary"
                } transition-colors`}
              >
                <Grid3X3 className="w-5 h-5" />
              </button>
              <button
                onClick={() => setViewMode("list")}
                className={`p-2 rounded-lg ${
                  viewMode === "list"
                    ? "bg-surface-container-high text-on-surface-variant"
                    : "text-outline-variant hover:text-primary"
                } transition-colors`}
              >
                <List className="w-5 h-5" />
              </button>
            </div>
          </div>

          {isLoading ? (
            <div className="flex items-center justify-center py-20">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
          ) : (
            <VideoGrid videos={displayVideos} onVideoClick={handleVideoClick} />
          )}

          {/* Pagination */}
          <footer className="mt-16 flex items-center justify-between py-6 border-t border-outline-variant/10">
            <div className="text-sm text-on-surface-variant">
              Showing <span className="font-bold text-on-surface">1</span> to{" "}
              <span className="font-bold text-on-surface">
                {Math.min(limit, displayTotal)}
              </span>{" "}
              of <span className="font-bold text-on-surface">{displayTotal}</span>{" "}
              videos
            </div>
            <div className="flex items-center gap-4">
              <button
                className="flex items-center gap-2 px-6 py-2.5 rounded-full border border-outline-variant text-sm font-medium text-on-surface-variant hover:bg-surface-container-low disabled:opacity-50"
                disabled={currentPage === 1}
                onClick={() => setCurrentPage(currentPage - 1)}
              >
                Prev
              </button>
              <div className="flex items-center gap-2">
                {Array.from({ length: displayTotalPages }, (_, i) => i + 1).map(
                  (page) => (
                    <button
                      key={page}
                      className={`w-10 h-10 flex items-center justify-center rounded-full ${
                        page === currentPage
                          ? "bg-primary text-on-primary font-bold"
                          : "hover:bg-surface-container-high"
                      }`}
                      onClick={() => setCurrentPage(page)}
                    >
                      {page}
                    </button>
                  )
                )}
              </div>
              <button
                className="flex items-center gap-2 px-6 py-2.5 rounded-full bg-surface-container-high text-sm font-medium text-on-surface hover:bg-surface-container-highest transition-colors disabled:opacity-50"
                disabled={currentPage === displayTotalPages}
                onClick={() => setCurrentPage(currentPage + 1)}
              >
                Next
              </button>
            </div>
          </footer>
        </section>
      </main>
    </MainLayout>
  );
}
