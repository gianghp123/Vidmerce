import { MainLayout } from "@/components/layout/MainLayout";
import { TopBar } from "@/components/layout/TopBar";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { SequenceTimeline } from "../components/SequenceTimeline";
import { VideoConfigCard } from "../components/VideoConfigCard";
import { VideoPreview } from "../components/VideoPreview";
import { useCreateVideo } from "../hooks/useCreateVideo";
import type { SequenceItem, VideoStyle } from "../models/video-sequence.model";

const mockAssets = [
  { id: "1", name: "Minimalist Clay Vase", imageUrl: "https://lh3.googleusercontent.com/aida-public/AB6AXuBuTZrsM5VRVj5dQeyxi9_7Uuk5mKKJzkUtLK3wpqM4y1r-p3Fr_X-3guCLZGxnQloCeMYEcO-GT_5CLAcfTWl4b_dDr6rbVLw9fbOalwrN27K0nN0TPGGuJXuoMc80KO4DVDDk_84NXh4CD7U7PtdkQYTfNlUikXuL8MGj60SjGgpf8xWiL7S8iUR5YGrx2VRaeSQQfWuBYl7YAaAxsWnuWl5UbvkxaYES4FybrOWEiZrQUSHnnbDzTSa2bEU8b_Adw0g-9-Yjeejt" },
  { id: "2", name: "Artisan Leather Watch", imageUrl: "https://lh3.googleusercontent.com/aida-public/AB6AXuAka8Yt0yod9MRY5ik9TRRrodtN6ZBnpPlkFyXYYzd93_Ft1s14evy_UmmzDkBXqKQbBTSg74ITOHfmnLMoHAD_tgafvDdDasbDiQoXEDK2rfB_a8xPUoAXbwynjFZ9yES0VUcyhS5JbTGxF2IqpSnOSL2O23sq9c_qfJKkCSKbQ5THVBVdOlk7bl6ju6zMGmXXoviOWXg68W_sjZnXWFPMETuKDMVi5SqOtCXAw42Jgh5wqqczXOLSF4QTTDvUjAd1w6D-nzOiEOwe" },
  { id: "3", name: "Cedar & Amber Candle", imageUrl: "https://lh3.googleusercontent.com/aida-public/AB6AXuBaeZig8jemDV41r7DTALl_soOJ2Sz-AvDkWI-x3XeMoYQMw_O_vXebCil1lJMPE4mfNWJHLziLu2KPcBYBDuDpJPkHhCyg4lk2wbMFj0WOqa4iNGkTJAiLUI83ONnAld_1lcfNcR7YqNq-oj2WmB3xJWcnD4zbSKE4Of5e3l0j8LuzpmX-nro867PDoaix7Bk2v4fPgF8shIuHQFg8nhfYt7IuM1n7trUAzbp6pfHVtChbQAd4lhj1W1eIaduH-kVlAoXfCdAAKMy" },
  { id: "4", name: "Botanical Face Oil", imageUrl: "https://lh3.googleusercontent.com/aida-public/AB6AXuBr-dggq-4mOkngADpHdVixeCMy7qspFU6KJjv5StBacTMh9Wg78U6VjjhTOIWIwRxY9l99fUI81r_q2N-LUyhhNxdg_kAMUEtk_jeAcR81dZk94bQ3YsMewRep0KpQ5cSvG1gqIbKdHSQuXaHqCnbYilTZDZQw28jnGM7gRwX_QLyFGFvwY8bcAqkOh8JWM_UrLPXKdELOSqJqbTqxM45mCw4V2qzg_ng1eJS0o1Xpz4vnB_lGNnsnHiaxtavNB-Vnw3HL5g9dHTe1" },
];

export function VideoBuilderPage() {
  const navigate = useNavigate();
  const [title, setTitle] = useState("");
  const [style, setStyle] = useState<VideoStyle>("kenburns");
  const [sequences, setSequences] = useState<SequenceItem[]>([]);

  const { create, isLoading } = useCreateVideo({
    onSuccess: () => {
      navigate("/videos");
    },
    onError: (error) => {
      console.error("Failed to create video:", error);
    },
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
    <MainLayout>
      <TopBar />

      <main className="p-12 max-w-6xl mx-auto bg-background">
        <div className="mb-12">
          <h2 className="text-4xl font-bold text-on-surface tracking-tight mb-2">
            Create Video
          </h2>
          <p className="text-on-surface-variant text-[15px] max-w-xl leading-relaxed">
            Curate your brand story by stitching high-quality assets into a seamless editorial experience.
          </p>
        </div>

        <div className="grid grid-cols-12 gap-10 items-start">
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

          <div className="col-span-12 lg:col-span-4 sticky top-32 space-y-6">
            <VideoPreview sequences={sequences} />

            <div className="p-6 rounded-3xl bg-surface-container-low/50">
              <div className="space-y-4 mb-6">
                <div className="flex justify-between items-center">
                  <span className="text-sm font-bold text-on-surface-variant">Estimated Length</span>
                  <span className="text-sm font-extrabold text-on-surface">15.0 Seconds</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm font-bold text-on-surface-variant">Resolution</span>
                  <span className="text-sm font-extrabold text-on-surface">1080 × 1920 (9:16)</span>
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
                For a professional editorial look, keep transitions consistent and durations between 3-5 seconds per asset.
              </p>
            </div>
          </div>
        </div>
      </main>
    </MainLayout>
  );
}
