import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { VideoStyle } from "../models/video-sequence.model";

interface VideoConfigCardProps {
  title: string;
  style: VideoStyle;
  onTitleChange: (title: string) => void;
  onStyleChange: (style: VideoStyle) => void;
}

const styleOptions: { value: VideoStyle; label: string }[] = [
  { value: "kenburns", label: "Ken Burns" },
  { value: "cinematic", label: "Cinematic Minimalist" },
  { value: "static", label: "Static Editorial" },
  { value: "fastcut", label: "Fast Cut Grid" },
];

export function VideoConfigCard({ title, style, onTitleChange, onStyleChange }: VideoConfigCardProps) {
  return (
    <div className="bg-surface-container-low p-10 rounded-[32px] space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div className="space-y-3">
          <Label className="text-[11px] font-bold uppercase tracking-[0.15em] text-secondary">
            Video Title
          </Label>
          <Input
            value={title}
            onChange={(e) => onTitleChange(e.target.value)}
            placeholder="Summer Collection 2024"
            className="h-12 bg-white border-none rounded-xl px-4 text-on-surface focus-visible:ring-1 focus-visible:ring-outline-variant shadow-sm"
          />
        </div>

        <div className="space-y-3">
          <Label className="text-[11px] font-bold uppercase tracking-[0.15em] text-secondary">
            Video Style
          </Label>
          <Select value={style} onValueChange={onStyleChange}>
            <SelectTrigger className="bg-white border-none rounded-xl px-4 text-on-surface shadow-sm">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {styleOptions.map((opt) => (
                <SelectItem key={opt.value} value={opt.value}>{opt.label}</SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>
    </div>
  );
}