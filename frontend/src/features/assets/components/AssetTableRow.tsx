import { TableCell, TableRow } from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export interface AssetTableRowData {
  id: string;
  name: string;
  sku: string;
  price: number;
  stock: 'IN STOCK' | 'LOW STOCK' | 'OUT OF STOCK';
  imageUrl: string;
}

interface AssetTableRowProps {
  asset: AssetTableRowData;
}

export function AssetTableRow({ asset }: AssetTableRowProps) {
  return (
    <TableRow className="group hover:bg-muted/5 transition-colors">
      <TableCell className="px-8 py-6">
        <div className="w-16 h-16 rounded bg-muted overflow-hidden border border-border/10">
          <img 
            className="w-full h-full object-cover" 
            src={asset.imageUrl} 
            alt={asset.name}
            referrerPolicy="no-referrer"
          />
        </div>
      </TableCell>
      <TableCell className="px-8 py-6">
        <div className="font-semibold text-foreground text-sm">{asset.name}</div>
        <div className="text-[10px] text-muted-foreground mt-1 font-bold uppercase tracking-wider">SKU: {asset.sku}</div>
      </TableCell>
      <TableCell className="px-8 py-6 text-sm font-medium">${asset.price.toFixed(2)}</TableCell>
      <TableCell className="px-8 py-6">
        <Badge variant="outline" className={cn(
          "text-[10px] font-bold tracking-tighter px-2",
          asset.stock === 'IN STOCK' ? "bg-primary/5 text-primary border-primary/20" : "bg-muted text-muted-foreground"
        )}>
          {asset.stock}
        </Badge>
      </TableCell>
      <TableCell className="px-8 py-6 text-right">
        <Button variant="link" className="text-primary font-bold text-sm h-auto p-0">Create Campaign</Button>
      </TableCell>
    </TableRow>
  );
}