import { ExternalLink, Sparkles } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export function AssetInfoCards() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mt-4">
      <div className="bg-muted p-8 rounded-xl flex flex-col gap-4">
        <h3 className="text-lg font-bold font-serif">Data Intelligence Update</h3>
        <p className="text-sm text-muted-foreground leading-relaxed">
          Your library has been synchronized with the master inventory. 12 new product assets were detected and auto-tagged with AI metadata for immediate campaign generation.
        </p>
        <div className="mt-2">
          <button className="text-primary text-sm font-bold flex items-center gap-1 hover:underline">
            View Changelog
            <ExternalLink className="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
      
      <Card className="bg-muted-foreground bg-opacity-5 p-8 border border-primary/10 rounded-xl">
        <CardContent className="p-0 flex items-start justify-between">
          <div className="space-y-4">
            <h3 className="text-lg font-bold font-serif">Studio Integration</h3>
            <p className="text-sm text-muted-foreground">Ready to transform these raw assets into high-converting visual stories?</p>
            <Button className="bronze-gradient text-white px-8 h-12 shadow-lg shadow-primary/10 font-bold">
              Launch Studio Designer
            </Button>
          </div>
          <Sparkles className="w-12 h-12 text-primary/20" />
        </CardContent>
      </Card>
    </div>
  );
}