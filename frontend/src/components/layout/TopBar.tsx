import { Search, Bell, Plus } from "lucide-react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

interface TopBarProps {
  onSearch?: (query: string) => void;
  onNewAsset?: () => void;
}

export function TopBar({ onSearch, onNewAsset }: TopBarProps) {
  return (
    <header className="w-full sticky top-0 z-40 bg-surface-background/80 backdrop-blur-md shadow-[0px_20px_40px_rgba(70,33,7,0.06)] flex justify-between items-center px-12 py-6">
      <div className="flex items-center gap-8">
        <div className="relative group w-full max-w-md">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-outline-variant group-focus-within:text-primary w-4 h-4" />
          <Input
            className="pl-10 pr-4 py-2 bg-surface-container-high/50 border-none rounded-full text-sm w-full focus:ring-2 focus:ring-primary/20 placeholder:text-outline-variant"
            placeholder="Search library..."
            type="text"
            onChange={(e) => onSearch?.(e.target.value)}
          />
        </div>
        <nav className="hidden md:flex items-center gap-6">
          <Link
            to="/analytics"
            className="font-heading text-sm font-medium text-primary/60 hover:text-[#4B2E2B] transition-colors"
          >
            Analytics
          </Link>
          <Link
            to="/library"
            className="font-heading text-sm font-medium text-[#4B2E2B] border-b-2 border-primary pb-1"
          >
            Library
          </Link>
          <Link
            to="/templates"
            className="font-heading text-sm font-medium text-primary/60 hover:text-[#4B2E2B] transition-colors"
          >
            Templates
          </Link>
        </nav>
      </div>
      <div className="flex items-center gap-4">
        <Button
          className="hidden lg:flex items-center px-6 py-2.5 text-sm font-heading font-bold text-primary bg-surface-container-high hover:bg-surface-container-highest rounded-xl transition-colors active:scale-95"
        >
          Export
        </Button>
        <Button
          onClick={onNewAsset}
          className="bg-linear-to-br from-secondary to-primary-container text-primary-foreground px-4 py-3 rounded-xl font-heading font-bold text-sm flex items-center justify-center gap-2 shadow-lg active:opacity-80 transition-opacity"
        >
          <Plus className="w-4 h-4" />
          New Asset
        </Button>
        <div className="flex items-center gap-2 border-l border-outline-variant/20 ml-4 pl-4">
          <Button variant="ghost" size="icon" className="text-primary hover:bg-surface-container-high rounded-full">
            <Bell className="w-5 h-5" />
          </Button>
          <div className="w-10 h-10 rounded-full bg-surface-container-low border-2 border-surface-container-high overflow-hidden">
            <img
              alt="User Profile"
              className="w-full h-full object-cover"
              src="https://lh3.googleusercontent.com/aida-public/AB6AXuDMkASMTSMdWp2M4FvEAA3oMbgs1UngV7ehWYM-mWtTzqpSPZmLls8IjcoZRM4sFVe8nVlMGgAoEY1L9R68KgCuyGgnQcj6hdmJRqgWKDA287Xm6HP_2a_jXuaoQ3ZAT6_NwMd-9g2wTh5mz7FTJ2rdPRYwDn6xJTWMKQgMFDRYMvYcoGiB3KBNCNaRd0s0NrcXZKBMsa7bTKFFytY2iqKF5cekOBo5TZam3uD4sPw5ffS-l1Ol4-TBH-TZC4BO5dQkTFKt2ZCF47Mq"
            />
          </div>
        </div>
      </div>
    </header>
  );
}