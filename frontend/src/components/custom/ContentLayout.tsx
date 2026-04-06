import type { ReactNode } from "react";

interface ContentLayoutProps {
  title?: string;
  description?: string;
  actions?: ReactNode;
  children: ReactNode;
}

export function ContentLayout({
  title,
  description,
  actions,
  children,
}: ContentLayoutProps) {
  return (
    <main className="p-12 space-y-16 max-w-7xl mx-auto w-full">
      <div className="mb-12 flex flex-col md:flex-row md:items-end justify-between gap-6">
        <div>
          {title && (
            <h2 className="text-4xl font-extrabold font-heading text-on-surface tracking-tight mb-2">
              {title}
            </h2>
          )}
          {description && (
            <p className="text-on-surface-variant font-body">
              {description}
            </p>
          )}
        </div>
        {actions && (
          <div className="flex items-center gap-3">{actions}</div>
        )}
      </div>

      {children}
    </main>
  );
}