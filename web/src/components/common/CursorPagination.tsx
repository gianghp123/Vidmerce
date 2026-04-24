"use client";

import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { useRouter, useSearchParams } from "next/navigation";
import { useCallback, useMemo } from "react";
import { buildBackParams, buildForwardParams } from "@/lib/utils/pagination.util";

interface CursorPaginationProps {
  lastKey?: string | null;
  hasMore: boolean;
  limit: number;
}

export function CursorPagination({ lastKey, hasMore, limit }: CursorPaginationProps) {
  const router = useRouter();
  const searchParams = useSearchParams();

  const canGoBack = useMemo(() => !!searchParams.get("lastKey"), [searchParams]);

  const handleNext = useCallback(() => {
    if (!lastKey) return;
    const nextParams = buildForwardParams(searchParams, lastKey, limit);
    router.push(`?${nextParams.toString()}`);
  }, [searchParams, lastKey, limit, router]);

  const handlePrevious = useCallback(() => {
    const backParams = buildBackParams(searchParams);
    if (!backParams) return;
    router.push(`?${backParams.toString()}`);
  }, [searchParams, router]);

  return (
    <Pagination className="mt-4">
      <PaginationContent>
        <PaginationItem>
          <PaginationPrevious
            onClick={handlePrevious}
            className={canGoBack ? "cursor-pointer" : "pointer-events-none opacity-50"}
          />
        </PaginationItem>
        <PaginationItem>
          <PaginationNext
            onClick={handleNext}
            className={hasMore ? "cursor-pointer" : "pointer-events-none opacity-50"}
          />
        </PaginationItem>
      </PaginationContent>
    </Pagination>
  );
}