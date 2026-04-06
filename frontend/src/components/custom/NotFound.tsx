import { Button } from "@/components/ui/button";
import { ArrowLeft, TriangleAlert } from "lucide-react";

type NotFoundProps = {
  title?: string;
  description?: string;
  actionLabel?: string;
  onAction?: () => void;
};

export function NotFound({
  title = "Not Found",
  description = "The page or resource you're looking for doesn't exist.",
  actionLabel,
  onAction,
}: NotFoundProps) {
  return (
    <div className="flex flex-col items-center justify-center py-20 text-center">
      <div className="mb-4 text-error">
        <TriangleAlert className="w-12 h-12" />
      </div>

      <h3 className="text-xl font-semibold text-on-surface">{title}</h3>
      {description && (
        <p className="text-on-surface-variant mt-2 max-w-sm">
          {description}
        </p>
      )}

      {actionLabel && onAction && (
        <Button onClick={onAction} className="mt-6 gap-2">
          <ArrowLeft className="w-4 h-4" />
          {actionLabel}
        </Button>
      )}
    </div>
  );
}