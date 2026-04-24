"use client";

import { Button } from "@/components/ui/button";
import { TableCell, TableRow } from "@/components/ui/table";
import { IAsset } from "@/types/models/asset.model";

interface AssetTableRowProps {
  asset: IAsset;
}

export function AssetTableRow({ asset }: AssetTableRowProps) {
  return (
    <TableRow className="group hover:bg-muted/5 transition-colors">
      <TableCell className="px-8 py-6">
        <div className="w-16 h-16 rounded bg-muted overflow-hidden border border-border/10">
          <img
            className="w-full h-full object-cover"
            src={asset.images?.length || 0 > 0 ? asset?.images?.[0].fileKey: ""}
            alt={asset.name}
            referrerPolicy="no-referrer"
          />
        </div>
      </TableCell>
      <TableCell className="px-8 py-6">
        <div className="font-semibold text-foreground text-sm">{asset.name}</div>
        <div className="text-[10px] text-muted-foreground mt-1 font-bold uppercase tracking-wider">SKU: {asset.price}</div>
      </TableCell>
      <TableCell className="px-8 py-6 text-sm font-medium">${asset.price.toFixed(2)}</TableCell>
      <TableCell className="px-8 py-6 text-right">
        <Button variant="link" className="text-primary font-bold text-sm h-auto p-0">Create Campaign</Button>
      </TableCell>
    </TableRow>
  );
}