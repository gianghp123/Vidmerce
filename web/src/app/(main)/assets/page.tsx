import { PageLayout } from "@/components/layout/PageLayout";
import { AssetTable } from "@/features/assets/components/AssetTable";
import { fetchAssets } from "@/features/assets/services/asset.get";

export default async function AssetLibraryPage() {
  const  assetRes = await fetchAssets();
  const assets = assetRes.data || [];
  const meta = assetRes.meta || {};

  return (
    <PageLayout
      title="Asset Library"
      description="Manage raw product information and intelligence assets."
    >
      <AssetTable
        assets={assets}
        limit={10}
        hasMore={true}
      />
    </PageLayout>
  );
}