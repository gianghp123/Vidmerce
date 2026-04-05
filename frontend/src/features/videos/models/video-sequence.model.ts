export type TransitionType = "FADE" | "SLIDE" | "ZOOM";
export type VideoStyle = "kenburns" | "cinematic" | "static" | "fastcut";

export interface SequenceItem {
  id: string;
  assetId: string;
  assetName?: string;
  assetImageUrl?: string;
  duration: number;
  transition: TransitionType;
}

export interface VideoBuilder {
  title: string;
  style: VideoStyle;
  sequences: SequenceItem[];
}