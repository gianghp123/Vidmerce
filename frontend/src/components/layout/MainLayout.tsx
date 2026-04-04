import type { ReactNode } from "react";
import { Link, useLocation } from "react-router-dom";
import { Package, Video, Settings, HelpCircle, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";

interface MainLayoutProps {
  children: ReactNode;
}

export function MainLayout({ children }: MainLayoutProps) {
  const location = useLocation();
  const isVideosActive = location.pathname.startsWith("/videos");
  const isAssetsActive = location.pathname.startsWith("/library") || location.pathname === "/";

  return (
    <div className="flex min-h-screen">
      {/* Sidebar */}
      <aside className="h-screen w-64 flex-col fixed left-0 top-0 bg-surface-container-low flex p-6 z-50 border-none">
        <div className="mb-10">
          <h1 className="text-2xl font-bold text-[#4B2E2B] font-heading">Vidmerce</h1>
          <p className="text-xs font-medium text-primary/70 tracking-tight">The Editorial Merchant</p>
        </div>
        <nav className="flex-1 space-y-2">
          <Link
            to="/library"
            className={`flex items-center gap-3 px-4 py-3 ${isAssetsActive ? "bg-[#EEE7DF] text-[#4B2E2B] font-semibold" : "text-primary/70 hover:bg-[#EEE7DF]"} transition-colors duration-200 active:scale-95 rounded-lg`}
          >
            <Package className="w-5 h-5" />
            <span className="font-heading text-sm tracking-tight">Assets</span>
          </Link>
          <Link
            to="/videos"
            className={`flex items-center gap-3 px-4 py-3 ${isVideosActive ? "bg-[#EEE7DF] text-[#4B2E2B] font-semibold" : "text-primary/70 hover:bg-[#EEE7DF]"} transition-colors duration-200 active:scale-95 rounded-lg`}
          >
            <Video className="w-5 h-5" />
            <span className="font-heading text-sm tracking-tight">Videos</span>
          </Link>
        </nav>
        <div className="mt-auto space-y-2 pt-6">
          <Button className="w-full flex items-center justify-center gap-2 bg-gradient-to-br from-secondary to-primary-container text-primary-foreground py-3 px-4 rounded-xl font-heading font-bold text-sm shadow-lg active:opacity-80 transition-opacity mb-6">
            <Plus className="w-4 h-4" />
            New Asset
          </Button>
          <a
            className="flex items-center gap-3 px-4 py-2 text-primary/70 hover:bg-[#EEE7DF] transition-colors"
            href="#"
          >
            <Settings className="w-5 h-5" />
            <span className="font-heading text-sm tracking-tight">Settings</span>
          </a>
          <a
            className="flex items-center gap-3 px-4 py-2 text-primary/70 hover:bg-[#EEE7DF] transition-colors"
            href="#"
          >
            <HelpCircle className="w-5 h-5" />
            <span className="font-heading text-sm tracking-tight">Support</span>
          </a>
        </div>
      </aside>
      
      {/* Main Content */}
      <div className="ml-64 min-h-screen flex flex-col flex-1">
        {children}
      </div>
    </div>
  );
}