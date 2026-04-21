import { useState, useEffect } from "react";
import { toast } from "sonner";
import { AssetPageHeader } from "../components/AssetPageHeader";
import { AssetStats } from "../components/AssetStats";
import { AssetTable } from "../components/AssetTable";
import { AssetInfoCards } from "../components/AssetInfoCards";
import { useAssets } from "../hooks/useAssets";
import type { AssetTableRowData } from "../components/AssetTableRow";

const MOCK_TABLE_DATA: AssetTableRowData[] = [
  { id: "1", name: "Terra Chronograph v.2", sku: "TR-CH-2024", price: 420, stock: "IN STOCK", imageUrl: "https://picsum.photos/seed/watch1/200/200" },
  { id: "2", name: "Crimson Kinetic Runner", sku: "CR-KR-101", price: 185, stock: "IN STOCK", imageUrl: "https://picsum.photos/seed/shoe/200/200" },
  { id: "3", name: "Aura Soundscape Pro", sku: "AU-SP-90", price: 349.50, stock: "LOW STOCK", imageUrl: "https://picsum.photos/seed/headphone/200/200" },
  { id: "4", name: "Nomad Satchel Limited", sku: "NM-SL-22", price: 890, stock: "IN STOCK", imageUrl: "https://picsum.photos/seed/bag/200/200" },
];

export function AssetLibraryPage() {
  const { assets, isLoading, error } = useAssets();
  
  useEffect(() => {
    if (error) {
      toast.error(error.message);
    }
  }, [error]);

  return (
    <div className="py-12 px-12 max-w-[1720px] mx-auto space-y-12">
      <AssetPageHeader />
      
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <AssetStats />
      </div>
      
      <AssetTable 
        assets={MOCK_TABLE_DATA}
        currentPage={1}
        totalItems={1248}
      />
      
      <AssetInfoCards />
    </div>
  );
}