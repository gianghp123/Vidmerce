import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";

type EmptyStateProps = {
  title: string;
  description?: string;
  actionLabel?: string;
  onAction?: () => void;
};

export function EmptyState({
  title,
  description,
  actionLabel,
  onAction,
}: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-20 text-center">
      <div className="mb-4 text-muted-foreground">
        <svg
          width="48"
          height="48"
          fill="none"
          stroke="currentColor"
          strokeWidth="1.5"
        >
          <rect x="6" y="6" width="36" height="36" rx="8" />
          <path d="M16 24h16M24 16v16" />
        </svg>
      </div>

      <h3 className="text-xl font-semibold">{title}</h3>
      {description && (
        <p className="text-muted-foreground mt-2 max-w-sm">
          {description}
        </p>
      )}

      {actionLabel && onAction && (
        <Button onClick={onAction} className="mt-6 gap-2">
          <Plus className="w-4 h-4" />
          {actionLabel}
        </Button>
      )}
    </div>
  );
}