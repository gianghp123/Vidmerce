import { LibraryLayout } from "@/components/custom/LibraryLayout";
import { CursorPagination } from "@/components/custom/Pagination";
import { Button } from "@/components/ui/button";
import { Filter, Plus } from "lucide-react";
import { useState } from "react";
import { AssetGrid } from "../components/AssetGrid";
import { CreateAssetModal } from "../components/CreateAssetModal";
import { useAssets } from "../hooks/useAssets";
import { EmptyState } from "@/components/custom/EmptyState";

export function AssetLibraryPage() {
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const limit = 12;

  const { assets, isLoading, hasMore, fetchNext } = useAssets({ limit });
  const isEmpty = !isLoading && assets.length === 0;

  return (
    <>
      <LibraryLayout
        title="Asset Library"
        description="Manage your shoppable product catalog and creative assets."
        actions={
          <div className="flex items-center gap-2">
            <Button onClick={() => setIsCreateModalOpen(true)} className="gap-2 bg-linear-to-br from-secondary to-primary-container text-primary-foreground rounded-xl font-heading font-bold text-sm shadow-lg active:opacity-80 transition-opacity">
              <Plus className="w-4 h-4 mr-2" />
              New Asset
            </Button>
            <div className="bg-surface-container-low px-4 py-2 rounded-full flex items-center gap-2 text-sm font-medium text-on-surface-variant">
              <Filter className="w-4 h-4" />
              <span>Sort: Recently Added</span>
            </div>
          </div>
        }
      >
        <div className="pt-6">
          {isLoading && assets.length === 0 ? (
            <div className="flex items-center justify-center py-20">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
          ) : isEmpty ? (
            <EmptyState
              title="No asset yet"
              description="Start by creating your first asset."
            />
          ) : (
            <>
              <AssetGrid assets={assets} />

              <CursorPagination
                hasMore={hasMore}
                onLoadMore={fetchNext}
                isLoading={isLoading}
                entityName="assets"
              />
            </>
          )}
        </div>
      </LibraryLayout>
      <CreateAssetModal
        open={isCreateModalOpen}
        onOpenChange={setIsCreateModalOpen}
        onSuccess={() => { }}
      />
    </>
  );
}