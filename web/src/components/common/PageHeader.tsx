"use client";

import { ReactNode } from "react";

interface PageHeaderProps {
  title: string;
  description?: string;

  // Right side content (search, buttons, filters, etc.)
  headerButtons?: ReactNode;

  // Optional bottom content (tabs, breadcrumbs...)
  children?: ReactNode;
}

export function PageHeader({
  title,
  description,
  headerButtons,
  children,
}: PageHeaderProps) {
  return (
    <header className="space-y-4">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-end gap-6">
        {/* Left */}
        <div className="space-y-1">
          <h1 className="text-3xl font-bold tracking-tight">{title}</h1>
          {description && (
            <p className="text-muted-foreground text-sm">{description}</p>
          )}
        </div>

        {/* Right */}
        {headerButtons && (
          <div className="flex gap-3 w-full md:w-auto items-center">
            {headerButtons}
          </div>
        )}
      </div>

      {/* Bottom (tabs, filters, etc.) */}
      {children && <div>{children}</div>}
    </header>
  );
}