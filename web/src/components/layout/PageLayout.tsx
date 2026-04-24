"use client";

import { ReactNode } from "react";
import { PageHeader } from "@/components/common/PageHeader";

interface PageLayoutProps {
  title: string;
  description?: string;
  headerButtons?: ReactNode;
  children: ReactNode;
}

export function PageLayout({
  title,
  description,
  headerButtons,
  children,
}: PageLayoutProps) {
  return (
    <div className="py-12 px-12 md:px-24 w-full mx-auto space-y-12">
      <PageHeader
        title={title}
        description={description}
        headerButtons={headerButtons}
      />
      <main>{children}</main>
    </div>
  );
}