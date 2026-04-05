import {
  Pagination as ShadcnPagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
} from "@/components/ui/pagination";
import { Button } from "@/components/ui/button";
import { ChevronLeft, ChevronRight } from "lucide-react";

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  totalItems: number;
  pageSize: number;
  onPageChange: (page: number) => void;
  entityName?: string;
  className?: string;
}

export function Pagination({
  currentPage,
  totalPages,
  totalItems,
  pageSize,
  onPageChange,
  entityName = "items",
  className,
}: PaginationProps) {
  const start = (currentPage - 1) * pageSize + 1;
  const end = Math.min(currentPage * pageSize, totalItems);

  const getPageNumbers = (): (number | string)[] => {
    const pages: (number | string)[] = [];
    const maxVisible = 5;

    if (totalPages <= maxVisible) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      if (currentPage <= 3) {
        for (let i = 1; i <= 4; i++) pages.push(i);
        pages.push("...");
        pages.push(totalPages);
      } else if (currentPage >= totalPages - 2) {
        pages.push(1);
        pages.push("...");
        for (let i = totalPages - 3; i <= totalPages; i++) pages.push(i);
      } else {
        pages.push(1);
        pages.push("...");
        for (let i = currentPage - 1; i <= currentPage + 1; i++) pages.push(i);
        pages.push("...");
        pages.push(totalPages);
      }
    }

    return pages;
  };

  return (
    <footer className="mt-16 flex items-center justify-between py-6 border-t border-outline/10">
      <div className="text-sm text-on-surface-variant">
        Showing <span className="font-bold text-on-surface">{start}</span> to{" "}
        <span className="font-bold text-on-surface">{end}</span> of{" "}
        <span className="font-bold text-on-surface">{totalItems}</span>{" "}
        {entityName}
      </div>
      <div className="flex items-center gap-4">
        <Button
          variant="outline"
          size="default"
          onClick={() => onPageChange(currentPage - 1)}
          disabled={currentPage === 1}
          className="flex items-center gap-2 px-6 py-2.5 rounded-full border border-outline-variant text-sm font-medium text-on-surface-variant hover:bg-surface-container-low"
        >
          <ChevronLeft className="w-4 h-4" />
          Prev
        </Button>
        <ShadcnPagination className={className}>
          <PaginationContent>
            {getPageNumbers().map((page, idx) =>
              typeof page === "number" ? (
                <PaginationItem key={idx}>
                  <PaginationLink
                    onClick={() => onPageChange(page)}
                    isActive={page === currentPage}
                    className={`w-10 h-10 rounded-lg ${
                      page === currentPage
                        ? "bg-secondary text-primary-foreground font-bold"
                        : "hover:bg-surface-container-high"
                    }`}
                  >
                    {page}
                  </PaginationLink>
                </PaginationItem>
              ) : (
                <PaginationItem key={idx}>
                  <PaginationEllipsis />
                </PaginationItem>
              )
            )}
          </PaginationContent>
        </ShadcnPagination>
        <Button
          variant="outline"
          size="default"
          onClick={() => onPageChange(currentPage + 1)}
          disabled={currentPage === totalPages}
          className="flex items-center gap-2 px-6 py-2.5 rounded-full bg-surface-container-high text-sm font-medium text-on-surface hover:bg-surface-container-highest transition-colors"
        >
          Next
          <ChevronRight className="w-4 h-4" />
        </Button>
      </div>
    </footer>
  );
}
