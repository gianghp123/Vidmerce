"use client";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious
} from "@/components/ui/pagination";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { IAsset } from "@/types/models/asset.model";
import { AssetTableRow } from "./AssetTableRow";

interface AssetTableProps {
  assets: IAsset[];
  limit: number;
  hasMore: boolean;
}

export function AssetTable({
  assets,
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
        {
          assets.length > 0 ? <TableBody className="divide-y divide-muted">
            {assets.map((asset) => (
              <AssetTableRow key={asset.id} asset={asset} />
            ))}
          </TableBody> : <TableBody className="divide-y divide-muted">
            <TableRow className="bg-muted/20 border-border/30">
              <TableCell colSpan={5} className="px-8 py-6 text-center text-[10px] uppercase tracking-widest text-muted-foreground font-bold">No assets found</TableCell>
            </TableRow>
          </TableBody>
        }
      </Table>

      <Pagination>
        <PaginationContent>
          <PaginationItem>
            <PaginationPrevious />
          </PaginationItem>
          <PaginationItem>
            <PaginationNext />
          </PaginationItem>
        </PaginationContent>
      </Pagination>
    </div>
  );
}