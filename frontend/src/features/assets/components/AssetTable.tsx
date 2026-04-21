import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { AssetTableRow, type AssetTableRowData } from "./AssetTableRow";
import { AssetPagination } from "./AssetPagination";

interface AssetTableProps {
  assets: AssetTableRowData[];
  currentPage?: number;
  totalItems?: number;
  itemsPerPage?: number;
}

export function AssetTable({ 
  assets, 
  currentPage = 1, 
  totalItems = 1248,
  itemsPerPage = 10
}: AssetTableProps) {
  return (
    <div className="bg-white custom-shadow overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow className="bg-muted/20 border-border/30">
            <TableHead className="px-8 py-6 text-[10px] uppercase tracking-widest text-muted-foreground font-bold">Thumbnail</TableHead>
            <TableHead className="px-8 py-6 text-[10px] uppercase tracking-widest text-muted-foreground font-bold">Asset Name</TableHead>
            <TableHead className="px-8 py-6 text-[10px] uppercase tracking-widest text-muted-foreground font-bold">Price</TableHead>
            <TableHead className="px-8 py-6 text-[10px] uppercase tracking-widest text-muted-foreground font-bold">Inventory</TableHead>
            <TableHead className="px-8 py-6 text-[10px] uppercase tracking-widest text-muted-foreground font-bold text-right">Action</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody className="divide-y divide-muted">
          {assets.map((asset) => (
            <AssetTableRow key={asset.id} asset={asset} />
          ))}
        </TableBody>
      </Table>
      
      <AssetPagination currentPage={currentPage} totalItems={totalItems} itemsPerPage={itemsPerPage} />
    </div>
  );
}