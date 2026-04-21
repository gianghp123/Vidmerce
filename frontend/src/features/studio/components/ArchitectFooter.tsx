import { ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";

export function ArchitectFooter() {
  return (
    <footer className="pt-12 border-t border-border/10">
      <div className="bg-muted/30 p-8 rounded-xl flex flex-col md:flex-row items-center justify-between gap-6">
        <div className="space-y-1">
          <h3 className="text-lg text-foreground font-bold font-serif">Ready for Synthesis?</h3>
          <p className="text-xs text-muted-foreground italic font-serif">Estimated processing time: 14 seconds</p>
        </div>
        <Button size="lg" className="bronze-gradient text-white font-bold tracking-widest uppercase px-10 h-16 group">
          Compose Marketing Assets
          <ChevronRight className="ml-2 w-5 h-5 group-hover:translate-x-1 transition-transform" />
        </Button>
      </div>
    </footer>
  );
}