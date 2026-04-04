import { useState } from "react";
import { Filter } from "lucide-react";
import { MainLayout } from "@/components/layout/MainLayout";
import { TopBar } from "@/components/layout/TopBar";
import { AssetGrid } from "../components/AssetGrid";
import { AssetPagination } from "../components/AssetPagination";
import { useAssets } from "../hooks/useAssets";
import type { Asset } from "../models/asset.model";

const mockAssets: Asset[] = [
  {
    assetId: "1",
    name: "Minimalist Clay Vase",
    description: "High-end studio photography of a designer ceramic vase",
    price: 89,
    images: [{ imageUrl: "https://lh3.googleusercontent.com/aida-public/AB6AXuBuTZrsM5VRVj5dQeyxi9_7Uuk5mKKJzkUtLK3wpqM4y1r-p3Fr_X-3guCLZGxnQloCeMYEcO-GT_5CLAcfTWl4b_dDr6rbVLw9fbOalwrN27K0nN0TPGGuJXuoMc80KO4DVDDk_84NXh4CD7U7PtdkQYTfNlUikXuL8MGj60SjGgpf8xWiL7S8iUR5YGrx2VRaeSQQfWuBYl7YAaAxsWnuWl5UbvkxaYES4FybrOWEiZrQUSHnnbDzTSa2bEU8b_Adw0g-9-Yjeejt", order: 0 }],
    productUrl: "#",
    status: "COMPLETED",
    createdAt: "2024-01-15T10:00:00Z",
  },
  {
    assetId: "2",
    name: "Artisan Leather Watch",
    description: "Minimalist wrist watch with leather strap",
    price: 245,
    images: [{ imageUrl: "https://lh3.googleusercontent.com/aida-public/AB6AXuAka8Yt0yod9MRY5ik9TRRrodtN6ZBnpPlkFyXYYzd93_Ft1s14evy_UmmzDkBXqKQbBTSg74ITOHfmnLMoHAD_tgafvDdDasbDiQoXEDK2rfB_a8xPUoAXbwynjFZ9yES0VUcyhS5JbTGxF2IqpSnOSL2O23sq9c_qfJKkCSKbQ5THVBVdOlk7bl6ju6zMGmXXoviOWXg68W_sjZnXWFPMETuKDMVi5SqOtCXAw42Jgh5wqqczXOLSF4QTTDvUjAd1w6D-nzOiEOwe", order: 0 }],
    productUrl: "#",
    status: "COMPLETED",
    createdAt: "2024-01-14T09:30:00Z",
  },
  {
    assetId: "3",
    name: "Cedar & Amber Candle",
    description: "Luxury scented candle in a frosted glass jar",
    price: 42,
    images: [{ imageUrl: "https://lh3.googleusercontent.com/aida-public/AB6AXuBaeZig8jemDV41r7DTALl_soOJ2Sz-AvDkWI-x3XeMoYQMw_O_vXebCil1lJMPE4mfNWJHLziLu2KPcBYBDuDpJPkHhCyg4lk2wbMFj0WOqa4iNGkTJAiLUI83ONnAld_1lcfNcR7YqNq-oj2WmB3xJWcnD4zbSKE4Of5e3l0j8LuzpmnX-nro867PDoaix7Bk2v4fPgF8shIuHQFg8nhfYt7IuM1n7trUAzbp6pfHVtChbQAd4lhj1W1eIaduH-kVlAoXfCdAAKMy", order: 0 }],
    productUrl: "#",
    status: "DRAFT",
    createdAt: "2024-01-13T14:20:00Z",
  },
  {
    assetId: "4",
    name: "Botanical Face Oil",
    description: "Premium organic face serum bottle",
    price: 64,
    images: [{ imageUrl: "https://lh3.googleusercontent.com/aida-public/AB6AXuBr-dggq-4mOkngADpHdVixeCMy7qspFU6KJjv5StBacTMh9Wg78U6VjjhTOIWIwRxY9l99fUI81r_q2N-LUyhhNxdg_kAMUEtk_jeAcR81dZk94bQ3YsMewRep0KpQ5cSvG1gqIbKdHSQuXaHqCnbYilTZDZQw28jnGM7gRwX_QLyFGFvwY8bcAqkOh8JWM_UrLPXKdELOSqJqbTqxM45mCw4V2qzg_ng1eJS0o1Xpz4vnB_lGNnsnHiaxtavNB-Vnw3HL5g9dHTe1", order: 0 }],
    productUrl: "#",
    status: "COMPLETED",
    createdAt: "2024-01-12T11:45:00Z",
  },
];

export function AssetLibraryPage() {
  const [currentPage, setCurrentPage] = useState(1);
  const limit = 12;
  
  const { assets: fetchedAssets, isLoading } = useAssets({ page: currentPage, limit });
  
  const assets = fetchedAssets.length > 0 ? fetchedAssets : mockAssets;
  const total = mockAssets.length;
  const totalPages = Math.ceil(total / limit);

  return (
    <MainLayout>
      <TopBar />
      
      <main className="p-12 space-y-16 max-w-7xl mx-auto">
        {/* Header Section */}
        <div className="mb-12 flex flex-col md:flex-row md:items-end justify-between gap-6">
          <div>
            <h2 className="text-4xl font-extrabold font-heading text-on-surface tracking-tight mb-2">
              Asset Library
            </h2>
            <p className="text-on-surface-variant font-body">
              Manage your shoppable product catalog and creative assets.
            </p>
          </div>
          <div className="flex items-center gap-3">
            <div className="bg-surface-container-low px-4 py-2 rounded-full flex items-center gap-2 text-sm font-medium text-on-surface-variant">
              <Filter className="w-4 h-4" />
              <span>Sort: Recently Added</span>
            </div>
          </div>
        </div>

        {/* Loading State */}
        {isLoading ? (
          <div className="flex items-center justify-center py-20">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
          </div>
        ) : (
          <>
            {/* Asset Grid */}
            <AssetGrid assets={assets} />

            {/* Pagination */}
            <AssetPagination
              currentPage={currentPage}
              totalPages={totalPages}
              total={total}
              limit={limit}
              onPageChange={setCurrentPage}
            />
          </>
        )}
      </main>
    </MainLayout>
  );
}