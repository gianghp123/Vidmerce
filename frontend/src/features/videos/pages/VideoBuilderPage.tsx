import { LibraryLayout } from "@/components/custom/LibraryLayout";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { SequenceTimeline } from "../components/SequenceTimeline";
import { VideoConfigCard } from "../components/VideoConfigCard";
import { VideoPreview } from "../components/VideoPreview";
import { useCreateVideo } from "../hooks/useCreateVideo";
import type { SequenceItem, VideoStyle } from "../models/video-sequence.model";

const mockAssets = [
  { id: "1", name: "Minimalist Clay Vase", imageUrl: "..." },
  { id: "2", name: "Artisan Leather Watch", imageUrl: "..." },
  { id: "3", name: "Cedar & Amber Candle", imageUrl: "..." },
  { id: "4", name: "Botanical Face Oil", imageUrl: "..." },
];

export function VideoBuilderPage() {
  const navigate = useNavigate();
  const [title, setTitle] = useState("");
  const [style, setStyle] = useState<VideoStyle>("kenburns");
  const [sequences, setSequences] = useState<SequenceItem[]>([]);

  const { create, isLoading } = useCreateVideo({
    onSuccess: () => navigate("/videos"),
    onError: (error) => console.error("Failed to create video:", error),
  });

  const handleGenerateVideo = async () => {
    if (!title.trim() || sequences.length === 0) return;

    await create({
      title,
      style,
      sequences: sequences
        .filter((seq) => seq.assetId)
        .map((seq) => ({
          assetId: seq.assetId,
          duration: seq.duration,
          transition: seq.transition,
        })),
    });
  };

  const canGenerate = title.trim().length > 0 && sequences.length > 0;

  return (
    <LibraryLayout
      title="Create Video"
      description="Curate your brand story by stitching high-quality assets into a seamless editorial experience."
    >
      <div className="grid grid-cols-12 gap-10 items-start">
        
        {/* LEFT */}
        <div className="col-span-12 lg:col-span-8 space-y-8">
          <VideoConfigCard
            title={title}
            style={style}
            onTitleChange={setTitle}
            onStyleChange={setStyle}
          />

          <SequenceTimeline
            sequences={sequences}
            onSequencesChange={setSequences}
            availableAssets={mockAssets}
          />
        </div>

        {/* RIGHT */}
        <div className="col-span-12 lg:col-span-4 sticky top-32 space-y-6">
          <VideoPreview sequences={sequences} />

          <div className="p-6 rounded-3xl bg-surface-container-low/50">
            <div className="space-y-4 mb-6">
              <div className="flex justify-between items-center">
                <span className="text-sm font-bold text-on-surface-variant">
                  Estimated Length
                </span>
                <span className="text-sm font-extrabold text-on-surface">
                  15.0 Seconds
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm font-bold text-on-surface-variant">
                  Resolution
                </span>
                <span className="text-sm font-extrabold text-on-surface">
                  1080 × 1920 (9:16)
                </span>
              </div>
            </div>

            <Button
              onClick={handleGenerateVideo}
              disabled={!canGenerate || isLoading}
              className="w-full h-14 bg-secondary text-white rounded-2xl font-bold text-lg shadow-lg hover:opacity-90 transition-all"
            >
              {isLoading ? "Generating..." : "Generate Video"}
            </Button>

            <p className="text-[10px] text-center mt-4 text-on-surface-variant/60 italic leading-relaxed px-4">
              Rendering usually takes 2-3 minutes. You will be notified when your video is ready in the library.
            </p>
          </div>

          <div className="p-6 rounded-2xl bg-surface-container-low/50">
            <h4 className="text-sm font-bold text-on-surface mb-2 flex items-center gap-2">
              <span className="text-secondary text-lg">💡</span>
              Curator's Tip
            </h4>
            <p className="text-xs text-on-surface-variant leading-relaxed">
              For a professional editorial look, keep transitions consistent and durations between 3–5 seconds per asset.
            </p>
          </div>
        </div>
      </div>
    </LibraryLayout>
  );
}