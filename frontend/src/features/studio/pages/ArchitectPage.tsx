import { useState } from "react";
import { Package2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SourceMaterial } from "../components/SourceMaterial";
import { TemplateGrid } from "../components/TemplateGrid";
import { EditorialVoice } from "../components/EditorialVoice";
import { ArchitectFooter } from "../components/ArchitectFooter";
import type { ToneType } from "../models/studio.model";

export function ArchitectPage() {
  const [selectedTemplate, setSelectedTemplate] = useState(0);
  const [selectedTone, setSelectedTone] = useState<ToneType | null>('luxury');

  return (
    <div className="space-y-16 py-12 px-12 max-w-5xl mx-auto">
      <header className="flex justify-between items-end">
        <div className="space-y-2">
          <span className="technical-label text-primary">Campaign Architecture</span>
          <h1 className="text-4xl font-serif text-foreground">The Creative Blueprint</h1>
          <p className="text-muted-foreground max-w-xl text-sm leading-relaxed">
            Configure your visual narrative and structural flow. AdMint's intelligence will synthesize these parameters into refined marketing assets.
          </p>
        </div>
        <Button variant="outline" size="sm" className="bg-white/50 backdrop-blur-sm shadow-sm gap-2 uppercase text-[10px] tracking-wider font-bold">
          <Package2 className="w-4 h-4" />
          View Wireframe
        </Button>
      </header>

      <div className="grid grid-cols-1 gap-12">
        <section className="space-y-6">
          <div className="flex items-center space-x-3">
            <span className="w-8 h-8 rounded-full bg-muted flex items-center justify-center text-xs font-bold text-muted-foreground border border-border/20">02</span>
            <h2 className="text-xl text-foreground font-serif">Structural Framework</h2>
          </div>
          <TemplateGrid 
            selectedIndex={selectedTemplate} 
            onSelectTemplate={setSelectedTemplate} 
          />
        </section>

        <SourceMaterial selectedAsset={null} />
        
        <EditorialVoice 
          selectedTone={selectedTone}
          onSelectTone={setSelectedTone}
        />
      </div>

      <ArchitectFooter />
    </div>
  );
}