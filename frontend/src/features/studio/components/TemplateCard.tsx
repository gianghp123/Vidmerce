import { CheckCircle2, ColumnsIcon, LayoutDashboard, Maximize2, Monitor, Search, Sparkles, Video } from "lucide-react";
import { cn } from "@/lib/utils";
import type { TemplateType } from "../models/studio.model";

interface TemplateCardProps {
  title: string;
  description: string;
  type: TemplateType;
  isSelected: boolean;
  onClick: () => void;
}

export function TemplateCard({ title, description, type, isSelected, onClick }: TemplateCardProps) {
  return (
    <motion.div
      className="group cursor-pointer"
      whileHover={{ y: -4 }}
      onClick={onClick}
    >
      <div className={cn(
        "aspect-video bg-muted border p-0 mb-3 transition-all overflow-hidden relative",
        isSelected ? "border-accent ring-4 ring-accent/10" : "border-border/40 hover:border-primary/50"
      )}>
        {isSelected && (
          <div className="absolute top-2 right-2 z-20">
            <CheckCircle2 className="w-5 h-5 text-accent stroke-[3]" />
          </div>
        )}
        
        <div className="h-full w-full flex items-center justify-center relative">
          {type === 'carousel' && (
            <div className="h-full w-full linen-texture flex flex-col items-center justify-center gap-2">
              <div className="flex space-x-1.5 z-10">
                <div className="w-10 h-14 bg-white/40 backdrop-blur-sm shadow-sm border border-white/50"></div>
                <div className="w-12 h-16 bg-white shadow-md border border-white scale-110 z-20"></div>
                <div className="w-10 h-14 bg-white/40 backdrop-blur-sm shadow-sm border border-white/50"></div>
              </div>
              <div className="flex space-x-1 mt-1">
                <div className="w-1 h-1 rounded-full bg-accent"></div>
                <div className="w-1 h-1 rounded-full bg-black/10"></div>
                <div className="w-1 h-1 rounded-full bg-black/10"></div>
              </div>
            </div>
          )}
          {type === 'single' && (
            <div className="h-full w-full bg-atelier-shadow flex items-center justify-center">
              <div className="w-4/5 h-4/5 bg-white shadow-xl flex items-center justify-center">
                <div className="w-1/2 h-1/2 bg-muted/20 border border-border/10 flex items-center justify-center">
                  <Sparkles className="w-4 h-4 text-primary/20" />
                </div>
              </div>
            </div>
          )}
          {type === 'comparison' && (
            <div className="h-full w-full bg-[#fdfaf7] flex p-4 space-x-4">
              <div className="w-1/2 h-full bg-white shadow-sm flex items-center justify-center">
                <div className="w-10 h-10 border border-accent/10 bg-muted/20" />
              </div>
              <div className="w-1/2 flex flex-col justify-center space-y-3">
                <div className="w-full h-1 bg-accent/20"></div>
                <div className="w-3/4 h-1 bg-accent/10"></div>
                <div className="w-full h-1 bg-accent/20"></div>
              </div>
            </div>
          )}
          {type === 'deepdive' && (
            <div className="h-full w-full bg-[#eee7df] flex flex-col p-3 space-y-2">
              <div className="h-3/5 w-full bg-white shadow-sm flex items-center justify-center">
                <Search className="w-6 h-6 text-border/30" />
              </div>
              <div className="flex-1 grid grid-cols-2 gap-2">
                <div className="bg-white/80 border border-white shadow-inner"></div>
                <div className="bg-white/80 border border-white shadow-inner"></div>
              </div>
            </div>
          )}
          {type === 'gallery' && (
            <div className="h-full w-full linen-texture p-4 flex flex-col">
              <div className="w-1/3 h-1 bg-accent/30 mb-4"></div>
              <div className="flex-1 grid grid-cols-3 gap-2">
                <div className="bg-white shadow-sm border border-white/50"></div>
                <div className="bg-white shadow-sm border border-white/50 translate-y-1"></div>
                <div className="bg-white shadow-sm border border-white/50 -translate-y-1"></div>
              </div>
            </div>
          )}
          {type === 'video' && (
            <div className="h-full w-full bg-atelier-shadow flex items-center justify-center">
              <div className="w-2/3 h-4/5 bg-white shadow-2xl flex flex-col items-center justify-center space-y-4 border border-white relative overflow-hidden">
                <motion.div
                  animate={{ rotate: 360 }}
                  transition={{ repeat: Infinity, duration: 2, ease: "linear" }}
                  className="w-10 h-10 rounded-full border-2 border-accent/20 border-t-accent"
                />
                <div className="w-1/2 h-1 bg-accent/10"></div>
                <div className="absolute inset-0 bg-gradient-to-b from-transparent via-white/5 to-white/20" />
              </div>
            </div>
          )}
        </div>
      </div>
      <div className="flex flex-col">
        <p className="font-semibold text-sm text-foreground">{title}</p>
        <p className="text-xs text-muted-foreground leading-relaxed">{description}</p>
      </div>
    </motion.div>
  );
}

export const TEMPLATES = [
  {
    id: 'carousel',
    title: "3-Slide Carousel",
    description: "Sequential storytelling for core product features.",
    icon: <Monitor className="w-5 h-5" />,
    type: 'carousel' as TemplateType
  },
  {
    id: 'single',
    title: "Single Image Ad",
    description: "High-impact hero visual with minimal copy.",
    icon: <Maximize2 className="w-5 h-5" />,
    type: 'single' as TemplateType
  },
  {
    id: 'comparison',
    title: "Feature Comparison",
    description: "Data-driven layout for competitive edge.",
    icon: <ColumnsIcon className="w-5 h-5" />,
    type: 'comparison' as TemplateType
  },
  {
    id: 'deepdive',
    title: "Product Deep-Dive",
    description: "Comprehensive look at specifications and utility.",
    icon: <Search className="w-5 h-5" />,
    type: 'deepdive' as TemplateType
  },
  {
    id: 'gallery',
    title: "Gallery Grid",
    description: "Multi-angle showcase for lifestyle photography.",
    icon: <LayoutDashboard className="w-5 h-5" />,
    type: 'gallery' as TemplateType
  },
  {
    id: 'video',
    title: "Dynamic Video Reveal",
    description: "Motion-first framework for digital billboards.",
    icon: <Video className="w-5 h-5" />,
    type: 'video' as TemplateType
  },
];