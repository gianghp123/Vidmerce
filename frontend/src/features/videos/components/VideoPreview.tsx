import { Play } from "lucide-react";
import type { SequenceItem } from "../models/video-sequence.model";

interface VideoPreviewProps {
  sequences: SequenceItem[];
}

export function VideoPreview({  }: VideoPreviewProps) {
  return (
    <div className="relative aspect-9/16 bg-black rounded-[32px] overflow-hidden shadow-2xl ring-8 ring-on-surface/5">
      {/* Placeholder for video content */}
      <div className="absolute inset-0 bg-linear-to-t from-black/60 to-transparent z-10" />

      <div className="absolute inset-0 flex items-center justify-center z-20">
        <div className="w-16 h-12 bg-white/20 backdrop-blur-md rounded-xl flex items-center justify-center border border-white/30 cursor-pointer">
          <Play className="w-8 h-8 text-white fill-white" />
        </div>
      </div>

      <div className="absolute bottom-6 left-6 right-6 z-20 flex justify-between items-end">
        <div className="bg-black/60 backdrop-blur-md px-3 py-1 rounded-full text-[10px] text-white font-mono">
          00:07 / 00:15
        </div>
        <div className="flex gap-1.5 pb-1">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className={`h-1 rounded-full transition-all ${i === 2 ? 'w-8 bg-white' : 'w-4 bg-white/30'}`} />
          ))}
        </div>
      </div>
    </div>
  );
}