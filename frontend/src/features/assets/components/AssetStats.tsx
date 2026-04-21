import { Card, CardContent } from "@/components/ui/card";
import { ChevronRight } from "lucide-react";

interface AssetStatsProps {
  totalAssets?: number;
  activeCampaigns?: number;
  libraryValue?: string;
}

export function AssetStats({ 
  totalAssets = 1248, 
  activeCampaigns = 42, 
  libraryValue = "$248.5k" 
}: AssetStatsProps) {
  const stats = [
    { label: "Total Assets", val: totalAssets.toLocaleString(), color: "text-primary" },
    { label: "Active Campaigns", val: activeCampaigns.toString(), color: "text-secondary" },
    { label: "Library Value", val: libraryValue, color: "text-foreground" },
  ];

  return (
    <>
      {stats.map((stat, i) => (
        <Card key={i} className="bg-white border-border/10 custom-shadow">
          <CardContent className="p-6">
            <span className="text-[10px] uppercase tracking-widest text-muted-foreground font-bold block mb-2">{stat.label}</span>
            <div className={`text-3xl font-serif font-bold ${stat.color}`}>{stat.val}</div>
          </CardContent>
        </Card>
      ))}
      <Card className="bg-muted border-border/10 flex items-center justify-between p-6 cursor-pointer hover:bg-muted/80 transition-colors">
        <div>
          <span className="text-[10px] uppercase tracking-widest text-muted-foreground font-bold block mb-1">New Intent</span>
          <div className="text-sm font-medium">Capture Studio Data</div>
        </div>
        <div className="w-10 h-10 rounded-full bronze-gradient flex items-center justify-center text-white">
          <ChevronRight className="w-5 h-5" />
        </div>
      </Card>
    </>
  );
}