import { useLocation } from "react-router-dom";
import { Settings } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { cn } from "@/lib/utils";

export function TopBar() {
  const location = useLocation();
  const isStudio = location.pathname.includes("/studio");

  return (
    <header className="fixed top-0 left-0 right-0 h-[72px] bg-background border-b border-border/30 z-50">
      <div className="flex justify-between items-center w-full px-12 max-w-[1920px] mx-auto h-full">
        <div className="flex items-center gap-12">
          <div className="text-2xl font-serif italic text-primary font-bold">Vidmerce</div>
          <nav className="hidden md:flex items-center space-x-10 h-full">
            <a 
              href="/assets"
              className={cn(
                "h-full text-sm font-medium transition-colors relative",
                !isStudio ? "text-primary font-semibold" : "text-foreground/60 hover:text-primary"
              )}
            >
              Library
              {!isStudio && <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-accent" />}
            </a>
            <a 
              href="/studio"
              className={cn(
                "h-full text-sm font-medium transition-colors relative",
                isStudio ? "text-primary font-semibold" : "text-foreground/60 hover:text-primary"
              )}
            >
              Studio
              {isStudio && <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-accent" />}
            </a>
          </nav>
        </div>
        
        <div className="flex items-center space-x-6">
          <div className="flex items-center space-x-3">
            <Button variant="ghost" className="text-foreground font-medium text-sm">+ Add Product</Button>
            <Button className="bronze-gradient text-white border-none shadow-md">Export</Button>
          </div>
          <div className="flex items-center space-x-4 border-l border-border/20 pl-6">
            <Settings className="w-5 h-5 text-muted-foreground cursor-pointer hover:text-primary transition-colors" />
            <Avatar className="w-8 h-8 rounded-full border border-border/30">
              <AvatarImage src="https://lh3.googleusercontent.com/aida-public/AB6AXuBfmOSGCp2Q2eyWy2aoy8mcmo2GwOAnXJyFSru_QawdEqBXn7pCl5WGSKb76YvcLrJKcGKI23LUIdmaDyQp7EzP3BVHGLl2EGKLqxqsaONMwQRaSM5Ha4GWSay3GX3Vn6eVsQmEZO8pGRWkehn-BShyNHRCR-9cQOr_VjxvOlRP0dILPNLs_SvzAqszi4zFgCwnyiulsHLoaoKsc5bhI17lmQ44iEynf9UWkYvwSbZe4h6i0VSSbN56sDmh9_rg7QXmxAmZ2N1wN9nA" />
              <AvatarFallback>CC</AvatarFallback>
            </Avatar>
          </div>
        </div>
      </div>
    </header>
  );
}