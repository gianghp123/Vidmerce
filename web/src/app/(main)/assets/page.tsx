import { PageLayout } from "@/components/layout/PageLayout";
import { AssetTable } from "@/features/assets/components/AssetTable";
import { fetchAssets } from "@/features/assets/services/asset.get";
import { AssetsHeaderActions } from "@/features/assets/components/AssetsHeaderActions";

interface AssetsPageProps {
  searchParams: Promise<{
    lastKey?: string;
    limit?: string;
  }>;
}

export default async function AssetLibraryPage({
  searchParams,
}: AssetsPageProps) {
  const params = await searchParams;
  const limit = params.limit ? parseInt(params.limit) : 10;

  const lastKey = params.lastKey || undefined;
  const assetRes = await fetchAssets({ lastKey, limit });
  const assets = assetRes.data || [];
  const meta = assetRes.meta;

  return (
    <PageLayout
      title="Asset Library"
      description="Manage raw product information and intelligence assets."
      headerButtons={<AssetsHeaderActions />}
    >
      <AssetTable
        assets={assets}
        lastKey={meta?.lastKey}
        hasMore={meta?.hasMore || false}
        limit={meta?.limit || limit}
      />
    </PageLayout>
  );
}