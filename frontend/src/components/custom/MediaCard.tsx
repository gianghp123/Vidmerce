import { cn } from "@/lib/utils";
import { Card, CardContent } from "@/components/ui/card";
import { AspectRatio } from "@/components/ui/aspect-ratio";
import type { ReactNode } from "react";

type AspectRatioValue = "4/3" | "4/5" | "16/9" | "1/1";

const aspectRatioMap: Record<AspectRatioValue, number> = {
  "4/3": 4 / 3,
  "4/5": 4 / 5,
  "16/9": 16 / 9,
  "1/1": 1,
};

interface MediaCardProps {
  children: ReactNode;
  onClick?: () => void;
  aspectRatio?: AspectRatioValue;
  className?: string;
  thumbnail?: ReactNode;
  badge?: ReactNode;
}

export function MediaCard({
  children,
  onClick,
  aspectRatio = "4/3",
  className,
  thumbnail,
  badge,
}: MediaCardProps) {
  return (
    <Card
      className={cn(
        "bg-surface-container-lowest rounded-2xl overflow-hidden shadow-sm hover:shadow-xl transition-all duration-300 cursor-pointer",
        onClick && "cursor-pointer",
        className
      )}
      onClick={onClick}
    >
      <div className="relative bg-surface-container-low">
        {aspectRatio && (
          <AspectRatio ratio={aspectRatioMap[aspectRatio]}>
            {thumbnail}
          </AspectRatio>
        )}
        {badge && (
          <div className="absolute top-4 left-4">
            {badge}
          </div>
        )}
      </div>
      <CardContent className="p-6">
        {children}
      </CardContent>
    </Card>
  );
}
