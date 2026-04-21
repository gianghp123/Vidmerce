import { Outlet, useLocation, useNavigate } from "react-router-dom";
import { 
  Compass, 
  FileBox, 
  HelpCircle, 
  History, 
  LayoutGrid, 
  Sparkles 
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";


export function MainLayout() {
  const location = useLocation();
  const navigate = useNavigate();
  const isAssets = location.pathname === "/" || location.pathname.includes("/assets");
  const isStudio = location.pathname.includes("/studio");

  return (
    <div className="flex min-h-screen">
      <aside className="h-screen w-64 flex-col fixed left-0 top-0 bg-surface-container-low flex p-6 z-50 border-r border-border/20">
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
          <button
            onClick={() => navigate("/assets")}
            className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${
              isAssets 
                ? "bg-primary/10 text-primary font-semibold" 
                : "text-muted-foreground hover:bg-muted hover:text-foreground"
            }`}
          >
            <LayoutGrid className="w-5 h-5" />
            <span className="text-[11px] uppercase tracking-wider font-bold">Library</span>
          </button>
          
          <button
            onClick={() => navigate("/studio")}
            className={`w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${
              isStudio 
                ? "bg-primary/10 text-primary font-semibold" 
                : "text-muted-foreground hover:bg-muted hover:text-foreground"
            }`}
          >
            <Compass className="w-5 h-5" />
            <span className="text-[11px] uppercase tracking-wider font-bold">Architect</span>
          </button>
          
          <button
            className="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-muted-foreground/60 hover:bg-muted hover:text-foreground transition-colors"
          >
            <Sparkles className="w-5 h-5" />
            <span className="text-[11px] uppercase tracking-wider font-bold">Composition</span>
          </button>
          
          <button
            className="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-muted-foreground/60 hover:bg-muted hover:text-foreground transition-colors"
          >
            <FileBox className="w-5 h-5" />
            <span className="text-[11px] uppercase tracking-wider font-bold">Logistics</span>
          </button>
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
      
      <div className="ml-64 min-h-screen flex flex-col flex-1 pt-[72px]">
        <Outlet />
      </div>
    </div>
  );
}