import { Sparkles } from "lucide-react";
import { Button } from "@/components/ui/button";
import type { ToneType } from "../models/studio.model";

interface EditorialVoiceProps {
  selectedTone: ToneType | null;
  onSelectTone: (tone: ToneType) => void;
}

export function EditorialVoice({ selectedTone, onSelectTone }: EditorialVoiceProps) {
  return (
    <section className="space-y-6">
      <div className="flex items-center space-x-3">
        <span className="w-8 h-8 rounded-full bg-muted flex items-center justify-center text-xs font-bold text-muted-foreground border border-border/20">03</span>
        <h2 className="text-xl text-foreground font-serif">Editorial Voice</h2>
      </div>
      <div className="flex flex-wrap gap-4">
        <Button 
          variant="outline" 
          className="px-8 h-12 rounded-full border-border/30 bg-muted/20"
          onClick={() => onSelectTone('professional')}
        >
          Professional
        </Button>
        <Button 
          className="px-8 h-12 rounded-full bronze-gradient text-white shadow-lg shadow-primary/20 gap-2"
          onClick={() => onSelectTone('luxury')}
        >
          <Sparkles className="w-4 h-4" />
          Luxury
        </Button>
        <Button 
          variant="outline" 
          className="px-8 h-12 rounded-full border-border/30 bg-muted/20"
          onClick={() => onSelectTone('high-energy')}
        >
          High-Energy
        </Button>
      </div>
    </section>
  );
}