"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import { CreateAssetModal } from "./CreateAssetModal";
import { useRouter } from "next/navigation";

export function AssetsHeaderActions() {
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const router = useRouter();

  const handleSuccess = () => {
    router.refresh();
  };

  return (
    <>
      <Button
        onClick={() => setIsCreateModalOpen(true)}
        className="flex items-center gap-2"
      >
        <Plus className="h-4 w-4" />
        Create Asset
      </Button>

      <CreateAssetModal
        open={isCreateModalOpen}
        onOpenChange={setIsCreateModalOpen}
        onSuccess={handleSuccess}
      />
    </>
  );
}
