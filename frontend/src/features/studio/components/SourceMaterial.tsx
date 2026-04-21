import { Search, CheckCircle2 } from "lucide-react";
import { Input } from "@/components/ui/input";

interface SourceMaterialProps {
  selectedAsset?: {
    id: string;
    name: string;
    imageUrl: string;
  } | null;
}

export function SourceMaterial({ selectedAsset }: SourceMaterialProps) {
  return (
    <section className="space-y-6">
      <div className="flex items-center space-x-3">
        <span className="w-8 h-8 rounded-full bg-muted flex items-center justify-center text-xs font-bold text-muted-foreground border border-border/20">01</span>
        <h2 className="text-xl text-foreground font-serif">Source Material</h2>
      </div>
      <div className="relative max-w-2xl">
        <div className="flex items-center space-x-4 bg-white border-b-2 border-border p-4 group focus-within:border-primary transition-all">
          <Search className="w-5 h-5 text-muted-foreground" />
          <Input 
            className="flex-1 bg-transparent border-none focus:ring-0 text-lg font-light p-0 h-auto focus-visible:ring-0 shadow-none" 
            placeholder="Select an asset from the Library..." 
          />
        </div>
        
        {selectedAsset && (
          <div className="mt-4 flex items-center space-x-4 p-4 bg-muted/30 rounded-lg border border-primary/20">
            <div className="w-12 h-12 rounded bg-muted overflow-hidden">
              <img 
                alt={selectedAsset.name}
                className="w-full h-full object-cover" 
                src={selectedAsset.imageUrl}
                referrerPolicy="no-referrer"
              />
            </div>
            <div className="flex-1">
              <p className="font-medium text-sm text-foreground">{selectedAsset.name}</p>
              <p className="text-[10px] uppercase tracking-tighter text-muted-foreground font-bold">Selected Product Asset</p>
            </div>
            <CheckCircle2 className="w-5 h-5 text-primary" />
          </div>
        )}
      </div>
    </section>
  );
}