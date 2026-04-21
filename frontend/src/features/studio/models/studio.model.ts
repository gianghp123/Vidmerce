import type { ReactNode } from "react";

export type TemplateType = 'carousel' | 'single' | 'comparison' | 'deepdive' | 'gallery' | 'video';

export type ToneType = 'professional' | 'luxury' | 'high-energy';

export interface ITemplate {
  id: string;
  title: string;
  description: string;
  icon: ReactNode;
  type: TemplateType;
}

export interface IArchitectState {
  selectedTemplate: number;
  selectedTone: ToneType | null;
  selectedAsset: {
    id: string;
    name: string;
    imageUrl: string;
  } | null;
}