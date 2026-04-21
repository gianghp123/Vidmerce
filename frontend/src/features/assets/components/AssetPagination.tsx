import { ChevronLeft, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";

interface AssetPaginationProps {
  currentPage?: number;
  totalItems?: number;
  itemsPerPage?: number;
}

export function AssetPagination({ 
  currentPage = 1, 
  totalItems = 1248, 
  itemsPerPage = 10 
}: AssetPaginationProps) {
  const totalPages = Math.ceil(totalItems / itemsPerPage);
  const startItem = (currentPage - 1) * itemsPerPage + 1;
  const endItem = Math.min(currentPage * itemsPerPage, totalItems);

  return (
    <div className="px-8 py-6 bg-muted/10 border-t border-border/10 flex justify-between items-center">
      <span className="text-[10px] text-muted-foreground font-bold uppercase tracking-widest">
        Showing {startItem}-{endItem} of {totalItems.toLocaleString()} assets
      </span>
      <div className="flex gap-2">
        <Button variant="outline" size="icon" className="w-8 h-8 rounded border-border/30">
          <ChevronLeft className="w-4 h-4" />
        </Button>
        <Button variant="outline" className="w-8 h-8 p-0 rounded border-primary text-primary font-bold text-xs bg-white">1</Button>
        <Button variant="ghost" className="w-8 h-8 p-0 rounded text-xs">2</Button>
        <Button variant="ghost" className="w-8 h-8 p-0 rounded text-xs">3</Button>
        <Button variant="outline" size="icon" className="w-8 h-8 rounded border-border/30">
          <ChevronRight className="w-4 h-4" />
        </Button>
      </div>
    </div>
  );
}