"use client";

import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { ROUTES } from "@/lib/routes";
import {
  Compass,
  FileBox,
  HelpCircle,
  History,
  LayoutGrid,
  Sparkles
} from "lucide-react";
import { usePathname, useRouter } from "next/navigation";


const menuItems = [
  {
    name: "Asset Library",
    href: ROUTES.MAIN.ASSETS.LIST,
    icon: <LayoutGrid className="w-5 h-5" />
  },
  {
    name: "Architect",
    href: ROUTES.MAIN.STUDIO.LIST,
    icon: <Compass className="w-5 h-5" />
  },
  {
    name: "Composition",
    href: ROUTES.MAIN.STUDIO.LIST,
    icon: <Sparkles className="w-5 h-5" />
  },
  {
    name: "Logistics",
    href: ROUTES.MAIN.STUDIO.LIST,
    icon: <FileBox className="w-5 h-5" />
  },
];

export function Sidebar() {
  const pathname = usePathname();
  const router = useRouter();

  return (
    <aside className="h-screen w-80 flex-col bg-surface-container-low flex p-6 z-50 border-r border-border/20">
      <div className="mb-8 p-6 -mx-6">
        <div className="flex items-center space-x-3 mb-4">
          <div className="w-8 h-8 bg-primary rounded-sm flex items-center justify-center text-white">
            <Compass className="w-5 h-5" />
          </div>
          <div>
            <div className="font-serif text-lg font-bold text-foreground">Vidmerce Studio</div>
            <div className="text-[10px] uppercase tracking-widest text-muted-foreground/70 font-semibold">Creative Intelligence</div>
          </div>
        </div>
      </div>

      <nav className="flex-1 space-y-1">
        {menuItems.map((item, index) => {
          // Determine if this item is active based on href
          const isActive = pathname === item.href || pathname.includes(item.href);

          return (
            <button
              key={index}
              onClick={() => router.push(item.href)}
              className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${isActive
                ? "bg-primary/10 text-primary font-semibold"
                : "text-muted-foreground hover:bg-muted hover:text-foreground"
                }`}
            >
              {item.icon}
              <span className="text-[11px] uppercase tracking-wider font-bold">
                {item.name === "Asset Library" ? "Library" : item.name}
              </span>
            </button>
          );
        })}
      </nav>

      <div className="mt-auto">
        <Button className="w-full bronze-gradient text-white py-6 mb-6">
          Compose Asset
        </Button>
        <Separator className="bg-border/20 mb-4" />
        <div className="space-y-1">
          <button className="flex items-center space-x-3 text-muted-foreground/70 w-full px-2 py-2 hover:text-primary transition-colors text-left">
            <HelpCircle className="w-4 h-4" />
            <span className="text-[10px] uppercase tracking-wider font-semibold">Help</span>
          </button>
          <button className="flex items-center space-x-3 text-muted-foreground/70 w-full px-2 py-2 hover:text-primary transition-colors text-left">
            <History className="w-4 h-4" />
            <span className="text-[10px] uppercase tracking-wider font-semibold">Archive</span>
          </button>
        </div>
      </div>
    </aside>
  );
}