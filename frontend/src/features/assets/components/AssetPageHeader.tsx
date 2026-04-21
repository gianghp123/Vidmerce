import { Search } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";

interface AssetPageHeaderProps {
  title?: string;
  description?: string;
  searchPlaceholder?: string;
  onSearch?: (query: string) => void;
}

export function AssetPageHeader({
  title = "Product Knowledge Base",
  description = "Manage raw product information and intelligence assets.",
  searchPlaceholder = "Search knowledge base...",
  onSearch
}: AssetPageHeaderProps) {
  return (
    <header className="flex flex-col md:flex-row justify-between items-start md:items-end gap-6">
      <div className="space-y-1">
        <h1 className="text-3xl font-bold text-foreground tracking-tight">{title}</h1>
        <p className="text-muted-foreground text-sm">{description}</p>
      </div>
      <div className="flex gap-4 w-full md:w-auto">
        <div className="flex-1 md:flex-none flex items-center bg-muted/30 px-4 py-2 rounded-lg border-b-2 border-border/50 focus-within:border-primary transition-all group">
          <Search className="w-5 h-5 text-muted-foreground mr-3" />
          <Input 
            className="bg-transparent border-none focus-visible:ring-0 shadow-none p-0 h-8 md:w-64" 
            placeholder={searchPlaceholder}
            onChange={(e) => onSearch?.(e.target.value)}
          />
        </div>
        <Button variant="ghost" className="gap-2 text-muted-foreground hover:text-primary transition-colors">
          <Search className="w-4 h-4" className="rotate-90" />
          <span className="text-[10px] uppercase font-bold tracking-widest hidden md:inline">Filter</span>
        </Button>
      </div>
    </header>
  );
}