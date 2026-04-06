import type { ReactNode } from "react";
import { MainLayout } from "@/components/layout/MainLayout";
import { TopBar } from "@/components/layout/TopBar";

interface PageLayoutProps {
  children: ReactNode;
}

export function PageLayout({ children }: PageLayoutProps) {
  return (
    <MainLayout>
      <TopBar />
      {children}
    </MainLayout>
  );
}