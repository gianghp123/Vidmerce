/**
 * @license
 * SPDX-License-Identifier: Apache-2.0
 */

import React, { useState } from 'react';
import { 
  SidebarProvider, 
  Sidebar, 
  SidebarContent, 
  SidebarHeader, 
  SidebarFooter, 
  SidebarGroup, 
  SidebarGroupContent, 
  SidebarMenu, 
  SidebarMenuItem, 
  SidebarMenuButton,
  SidebarInset,
  SidebarTrigger
} from "@/components/ui/sidebar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Separator } from "@/components/ui/separator";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { 
  LayoutGrid, 
  Compass, 
  Sparkles, 
  Package2, 
  Search, 
  Plus, 
  Settings, 
  HelpCircle, 
  History, 
  FileBox, 
  ChevronRight,
  Filter,
  ExternalLink,
  ChevronLeft,
  LayoutDashboard,
  ColumnsIcon,
  Video,
  Monitor,
  CheckCircle2,
  Maximize2
} from "lucide-react";
import { motion, AnimatePresence } from "motion/react";
import { cn } from "@/lib/utils";

// --- Types ---
type ViewState = 'studio' | 'library';

// --- Components ---

const TopNav = ({ currentView, setView }: { currentView: ViewState, setView: (v: ViewState) => void }) => (
  <header className="fixed top-0 left-0 right-0 h-[72px] bg-background border-b border-border/30 z-50">
    <div className="flex justify-between items-center w-full px-12 max-w-[1920px] mx-auto h-full">
      <div className="flex items-center gap-12">
        <div className="text-2xl font-serif italic text-primary font-bold">AdMint</div>
        <nav className="hidden md:flex items-center space-x-10 h-full">
          <button 
            onClick={() => setView('library')}
            className={cn(
              "h-full text-sm font-medium transition-colors relative",
              currentView === 'library' ? "text-primary font-semibold" : "text-foreground/60 hover:text-primary"
            )}
          >
            Library
            {currentView === 'library' && <motion.div layoutId="underline" className="absolute bottom-0 left-0 right-0 h-0.5 bg-accent" />}
          </button>
          <button 
            onClick={() => setView('studio')}
            className={cn(
              "h-full text-sm font-medium transition-colors relative",
              currentView === 'studio' ? "text-primary font-semibold" : "text-foreground/60 hover:text-primary"
            )}
          >
            Studio
            {currentView === 'studio' && <motion.div layoutId="underline" className="absolute bottom-0 left-0 right-0 h-0.5 bg-accent" />}
          </button>
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

const AppSidebar = ({ currentView, setView }: { currentView: ViewState, setView: (v: ViewState) => void }) => (
  <Sidebar className="mt-[72px] border-r border-border/20 bg-muted/50">
    <SidebarHeader className="p-6 mb-4">
      <div className="flex items-center space-x-3">
        <div className="w-8 h-8 bg-primary rounded-sm flex items-center justify-center text-white">
          <Compass className="w-5 h-5" />
        </div>
        <div>
          <div className="font-serif text-lg font-bold text-foreground">AdMint Studio</div>
          <div className="text-[10px] uppercase tracking-widest text-muted-foreground/70 font-semibold">Creative Intelligence</div>
        </div>
      </div>
    </SidebarHeader>
    <SidebarContent>
      <SidebarGroup>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton 
                isActive={currentView === 'library'}
                onClick={() => setView('library')}
                className="py-6 px-4"
              >
                <LayoutGrid className="w-5 h-5" />
                <span className="text-[11px] uppercase tracking-wider font-bold">Library</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <SidebarMenuButton 
                isActive={currentView === 'studio'}
                onClick={() => setView('studio')}
                className="py-6 px-4"
              >
                <Compass className="w-5 h-5" />
                <span className="text-[11px] uppercase tracking-wider font-bold">Architect</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <SidebarMenuButton className="py-6 px-4 text-foreground/60">
                <Sparkles className="w-5 h-5" />
                <span className="text-[11px] uppercase tracking-wider font-bold">Composition</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
            <SidebarMenuItem>
              <SidebarMenuButton className="py-6 px-4 text-foreground/60">
                <FileBox className="w-5 h-5" />
                <span className="text-[11px] uppercase tracking-wider font-bold">Logistics</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContent>
    <SidebarFooter className="p-6">
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
    </SidebarFooter>
  </Sidebar>
);

const ArchitectView = () => {
  const [selectedTemplate, setSelectedTemplate] = useState(0);

  const templates = [
    { 
      title: "3-Slide Carousel", 
      desc: "Sequential storytelling for core product features.", 
      icon: <Monitor className="w-5 h-5" />,
      type: 'carousel'
    },
    { 
      title: "Single Image Ad", 
      desc: "High-impact hero visual with minimal copy.", 
      icon: <Maximize2 className="w-5 h-5" />,
      type: 'single'
    },
    { 
      title: "Feature Comparison", 
      desc: "Data-driven layout for competitive edge.", 
      icon: <ColumnsIcon className="w-5 h-5" />,
      type: 'comparison'
    },
    { 
      title: "Product Deep-Dive", 
      desc: "Comprehensive look at specifications and utility.", 
      icon: <Search className="w-5 h-5" />,
      type: 'deepdive'
    },
    { 
      title: "Gallery Grid", 
      desc: "Multi-angle showcase for lifestyle photography.", 
      icon: <LayoutDashboard className="w-5 h-5" />,
      type: 'gallery'
    },
    { 
      title: "Dynamic Video Reveal", 
      desc: "Motion-first framework for digital billboards.", 
      icon: <Video className="w-5 h-5" />,
      type: 'video'
    },
  ];

  return (
    <div className="space-y-16 py-12 px-12 max-w-5xl mx-auto">
      <header className="flex justify-between items-end">
        <div className="space-y-2">
          <span className="technical-label text-primary">Campaign Architecture</span>
          <h1 className="text-4xl font-serif text-foreground">The Creative Blueprint</h1>
          <p className="text-muted-foreground max-w-xl text-sm leading-relaxed">
            Configure your visual narrative and structural flow. AdMint's intelligence will synthesize these parameters into refined marketing assets.
          </p>
        </div>
        <Button variant="outline" size="sm" className="bg-white/50 backdrop-blur-sm shadow-sm gap-2 uppercase text-[10px] tracking-wider font-bold">
          <Package2 className="w-4 h-4" />
          View Wireframe
        </Button>
      </header>

      <div className="grid grid-cols-1 gap-12">
        <section className="space-y-6">
          <div className="flex items-center space-x-3">
            <span className="w-8 h-8 rounded-full bg-muted flex items-center justify-center text-xs font-bold text-muted-foreground border border-border/20">01</span>
            <h2 className="text-xl text-foreground font-serif">Source Material</h2>
          </div>
          <div className="relative max-w-2xl">
            <div className="flex items-center space-x-4 bg-white border-b-2 border-border p-4 group focus-within:border-primary transition-all">
              <Search className="w-5 h-5 text-muted-foreground" />
              <Input 
                className="flex-1 bg-transparent border-none focus:ring-0 text-lg font-light p-0 h-auto focus-visible:ring-0 shadow-none" 
                placeholder="Select an asset from the Library..." 
              />
            </div>
            
            <div className="mt-4 flex items-center space-x-4 p-4 bg-muted/30 rounded-lg border border-primary/20">
              <div className="w-12 h-12 rounded bg-muted overflow-hidden">
                <img 
                  alt="Vanguard White Watch" 
                  className="w-full h-full object-cover" 
                  src="https://picsum.photos/seed/watch/200/200" 
                  referrerPolicy="no-referrer"
                />
              </div>
              <div className="flex-1">
                <p className="font-medium text-sm text-foreground">Vanguard Series: Chrono White</p>
                <p className="text-[10px] uppercase tracking-tighter text-muted-foreground font-bold">Selected Product Asset</p>
              </div>
              <CheckCircle2 className="w-5 h-5 text-primary" />
            </div>
          </div>
        </section>

        <section className="space-y-6">
          <div className="flex items-center space-x-3">
            <span className="w-8 h-8 rounded-full bg-muted flex items-center justify-center text-xs font-bold text-muted-foreground border border-border/20">02</span>
            <h2 className="text-xl text-foreground font-serif">Structural Framework</h2>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {templates.map((t, i) => (
              <motion.div 
                key={i} 
                className="group cursor-pointer"
                whileHover={{ y: -4 }}
                onClick={() => setSelectedTemplate(i)}
              >
                <div className={cn(
                  "aspect-video bg-muted border p-0 mb-3 transition-all overflow-hidden relative",
                  selectedTemplate === i ? "border-accent ring-4 ring-accent/10" : "border-border/40 hover:border-primary/50"
                )}>
                  {selectedTemplate === i && (
                    <div className="absolute top-2 right-2 z-20">
                      <CheckCircle2 className="w-5 h-5 text-accent stroke-[3]" />
                    </div>
                  )}
                  
                  {/* Template Visual Representations */}
                  <div className="h-full w-full flex items-center justify-center relative">
                    {t.type === 'carousel' && (
                      <div className="h-full w-full linen-texture flex flex-col items-center justify-center gap-2">
                        <div className="flex space-x-1.5 z-10">
                          <div className="w-10 h-14 bg-white/40 backdrop-blur-sm shadow-sm border border-white/50"></div>
                          <div className="w-12 h-16 bg-white shadow-md border border-white scale-110 z-20"></div>
                          <div className="w-10 h-14 bg-white/40 backdrop-blur-sm shadow-sm border border-white/50"></div>
                        </div>
                        <div className="flex space-x-1 mt-1">
                          <div className="w-1 h-1 rounded-full bg-accent"></div>
                          <div className="w-1 h-1 rounded-full bg-black/10"></div>
                          <div className="w-1 h-1 rounded-full bg-black/10"></div>
                        </div>
                      </div>
                    )}
                    {t.type === 'single' && (
                      <div className="h-full w-full bg-atelier-shadow flex items-center justify-center">
                         <div className="w-4/5 h-4/5 bg-white shadow-xl flex items-center justify-center">
                            <div className="w-1/2 h-1/2 bg-muted/20 border border-border/10 flex items-center justify-center">
                              <Sparkles className="w-4 h-4 text-primary/20" />
                            </div>
                         </div>
                      </div>
                    )}
                    {t.type === 'comparison' && (
                       <div className="h-full w-full bg-[#fdfaf7] flex p-4 space-x-4">
                          <div className="w-1/2 h-full bg-white shadow-sm flex items-center justify-center">
                            <div className="w-10 h-10 border border-accent/10 bg-muted/20" />
                          </div>
                          <div className="w-1/2 flex flex-col justify-center space-y-3">
                            <div className="w-full h-1 bg-accent/20"></div>
                            <div className="w-3/4 h-1 bg-accent/10"></div>
                            <div className="w-full h-1 bg-accent/20"></div>
                          </div>
                       </div>
                    )}
                    {t.type === 'deepdive' && (
                      <div className="h-full w-full bg-[#eee7df] flex flex-col p-3 space-y-2">
                        <div className="h-3/5 w-full bg-white shadow-sm flex items-center justify-center">
                          <Search className="w-6 h-6 text-border/30" />
                        </div>
                        <div className="flex-1 grid grid-cols-2 gap-2">
                          <div className="bg-white/80 border border-white shadow-inner"></div>
                          <div className="bg-white/80 border border-white shadow-inner"></div>
                        </div>
                      </div>
                    )}
                    {t.type === 'gallery' && (
                      <div className="h-full w-full linen-texture p-4 flex flex-col">
                        <div className="w-1/3 h-1 bg-accent/30 mb-4"></div>
                        <div className="flex-1 grid grid-cols-3 gap-2">
                          <div className="bg-white shadow-sm border border-white/50"></div>
                          <div className="bg-white shadow-sm border border-white/50 translate-y-1"></div>
                          <div className="bg-white shadow-sm border border-white/50 -translate-y-1"></div>
                        </div>
                      </div>
                    )}
                    {t.type === 'video' && (
                      <div className="h-full w-full bg-atelier-shadow flex items-center justify-center">
                        <div className="w-2/3 h-4/5 bg-white shadow-2xl flex flex-col items-center justify-center space-y-4 border border-white relative overflow-hidden">
                          <motion.div 
                            animate={{ rotate: 360 }}
                            transition={{ repeat: Infinity, duration: 2, ease: "linear" }}
                            className="w-10 h-10 rounded-full border-2 border-accent/20 border-t-accent" 
                          />
                          <div className="w-1/2 h-1 bg-accent/10"></div>
                          <div className="absolute inset-0 bg-gradient-to-b from-transparent via-white/5 to-white/20" />
                        </div>
                      </div>
                    )}
                  </div>
                </div>
                <div className="flex flex-col">
                  <p className="font-semibold text-sm text-foreground">{t.title}</p>
                  <p className="text-xs text-muted-foreground leading-relaxed">{t.desc}</p>
                </div>
              </motion.div>
            ))}
          </div>
        </section>

        <section className="space-y-6">
          <div className="flex items-center space-x-3">
            <span className="w-8 h-8 rounded-full bg-muted flex items-center justify-center text-xs font-bold text-muted-foreground border border-border/20">03</span>
            <h2 className="text-xl text-foreground font-serif">Editorial Voice</h2>
          </div>
          <div className="flex flex-wrap gap-4">
            <Button variant="outline" className="px-8 h-12 rounded-full border-border/30 bg-muted/20">Professional</Button>
            <Button className="px-8 h-12 rounded-full bronze-gradient text-white shadow-lg shadow-primary/20 gap-2">
              <Sparkles className="w-4 h-4" />
              Luxury
            </Button>
            <Button variant="outline" className="px-8 h-12 rounded-full border-border/30 bg-muted/20">High-Energy</Button>
          </div>
        </section>
      </div>

      <footer className="pt-12 border-t border-border/10">
        <div className="bg-muted/30 p-8 rounded-xl flex flex-col md:flex-row items-center justify-between gap-6">
          <div className="space-y-1">
            <h3 className="text-lg text-foreground font-bold font-serif">Ready for Synthesis?</h3>
            <p className="text-xs text-muted-foreground italic font-serif">Estimated processing time: 14 seconds</p>
          </div>
          <Button size="lg" className="bronze-gradient text-white font-bold tracking-widest uppercase px-10 h-16 group">
            Compose Marketing Assets
            <ChevronRight className="ml-2 w-5 h-5 group-hover:translate-x-1 transition-transform" />
          </Button>
        </div>
      </footer>
    </div>
  );
};

const LibraryView = () => {
  const assets = [
    { name: "Terra Chronograph v.2", sku: "TR-CH-2024", price: 420, stock: "IN STOCK", img: "https://picsum.photos/seed/watch1/200/200" },
    { name: "Crimson Kinetic Runner", sku: "CR-KR-101", price: 185, stock: "IN STOCK", img: "https://picsum.photos/seed/shoe/200/200" },
    { name: "Aura Soundscape Pro", sku: "AU-SP-90", price: 349.50, stock: "LOW STOCK", img: "https://picsum.photos/seed/headphone/200/200" },
    { name: "Nomad Satchel Limited", sku: "NM-SL-22", price: 890, stock: "IN STOCK", img: "https://picsum.photos/seed/bag/200/200" },
  ];

  return (
    <div className="py-12 px-12 max-w-[1720px] mx-auto space-y-12">
      <header className="flex flex-col md:flex-row justify-between items-start md:items-end gap-6">
        <div className="space-y-1">
          <h1 className="text-3xl font-bold text-foreground tracking-tight">Product Knowledge Base</h1>
          <p className="text-muted-foreground text-sm">Manage raw product information and intelligence assets.</p>
        </div>
        <div className="flex gap-4 w-full md:w-auto">
          <div className="flex-1 md:flex-none flex items-center bg-muted/30 px-4 py-2 rounded-lg border-b-2 border-border/50 focus-within:border-primary transition-all group">
            <Search className="w-5 h-5 text-muted-foreground mr-3" />
            <Input 
              className="bg-transparent border-none focus-visible:ring-0 shadow-none p-0 h-8 md:w-64" 
              placeholder="Search knowledge base..." 
            />
          </div>
          <Button variant="ghost" className="gap-2 text-muted-foreground hover:text-primary transition-colors">
            <Filter className="w-4 h-4" />
            <span className="text-[10px] uppercase font-bold tracking-widest hidden md:inline">Filter</span>
          </Button>
        </div>
      </header>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {[
          { label: "Total Assets", val: "1,248", color: "text-primary" },
          { label: "Active Campaigns", val: "42", color: "text-secondary" },
          { label: "Library Value", val: "$248.5k", color: "text-foreground" },
        ].map((stat, i) => (
          <Card key={i} className="bg-white border-border/10 custom-shadow">
            <CardContent className="p-6">
              <span className="text-[10px] uppercase tracking-widest text-muted-foreground font-bold block mb-2">{stat.label}</span>
              <div className={cn("text-3xl font-serif font-bold", stat.color)}>{stat.val}</div>
            </CardContent>
          </Card>
        ))}
        <Card className="bg-muted border-border/10 flex items-center justify-between p-6 cursor-pointer hover:bg-muted/80 transition-colors">
          <div>
            <span className="text-[10px] uppercase tracking-widest text-muted-foreground font-bold block mb-1">New Intent</span>
            <div className="text-sm font-medium">Capture Studio Data</div>
          </div>
          <div className="w-10 h-10 rounded-full bronze-gradient flex items-center justify-center text-white">
            <ChevronRight className="w-5 h-5" />
          </div>
        </Card>
      </div>

      <Card className="bg-white custom-shadow overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow className="bg-muted/20 border-border/30">
              <TableHead className="px-8 py-6 text-[10px] uppercase tracking-widest text-muted-foreground font-bold">Thumbnail</TableHead>
              <TableHead className="px-8 py-6 text-[10px] uppercase tracking-widest text-muted-foreground font-bold">Asset Name</TableHead>
              <TableHead className="px-8 py-6 text-[10px] uppercase tracking-widest text-muted-foreground font-bold">Price</TableHead>
              <TableHead className="px-8 py-6 text-[10px] uppercase tracking-widest text-muted-foreground font-bold">Inventory</TableHead>
              <TableHead className="px-8 py-6 text-[10px] uppercase tracking-widest text-muted-foreground font-bold text-right">Action</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody className="divide-y divide-muted">
            {assets.map((asset, i) => (
              <TableRow key={i} className="group hover:bg-muted/5 transition-colors">
                <TableCell className="px-8 py-6">
                  <div className="w-16 h-16 rounded bg-muted overflow-hidden border border-border/10">
                    <img 
                      className="w-full h-full object-cover" 
                      src={asset.img} 
                      referrerPolicy="no-referrer"
                    />
                  </div>
                </TableCell>
                <TableCell className="px-8 py-6">
                  <div className="font-semibold text-foreground text-sm">{asset.name}</div>
                  <div className="text-[10px] text-muted-foreground mt-1 font-bold uppercase tracking-wider">SKU: {asset.sku}</div>
                </TableCell>
                <TableCell className="px-8 py-6 text-sm font-medium">${asset.price.toFixed(2)}</TableCell>
                <TableCell className="px-8 py-6">
                  <Badge variant="outline" className={cn(
                    "text-[10px] font-bold tracking-tighter px-2",
                    asset.stock === 'IN STOCK' ? "bg-primary/5 text-primary border-primary/20" : "bg-muted text-muted-foreground"
                  )}>
                    {asset.stock}
                  </Badge>
                </TableCell>
                <TableCell className="px-8 py-6 text-right">
                  <Button variant="link" className="text-primary font-bold text-sm h-auto p-0">Create Campaign</Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        
        <div className="px-8 py-6 bg-muted/10 border-t border-border/10 flex justify-between items-center">
          <span className="text-[10px] text-muted-foreground font-bold uppercase tracking-widest">Showing 1-4 of 1,248 assets</span>
          <div className="flex gap-2">
            <Button variant="outline" size="icon" className="w-8 h-8 rounded border-border/30 animate-pulse-slow">
              <ChevronLeft className="w-4 h-4" />
            </Button>
            <Button variant="outline" className="w-8 h-8 p-0 rounded border-primary text-primary font-bold text-xs bg-white">1</Button>
            <Button variant="ghost" className="w-8 h-8 p-0 rounded text-xs">2</Button>
            <Button variant="ghost" className="w-8 h-8 p-0 rounded text-xs">3</Button>
            <Button variant="outline" size="icon" className="w-8 h-8 rounded border-border/30">
              <ChevronRight className="w-4 h-4" />
            </Button>
          </div>
        </div>
      </Card>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mt-4">
        <div className="bg-muted p-8 rounded-xl flex flex-col gap-4">
          <h3 className="text-lg font-bold font-serif">Data Intelligence Update</h3>
          <p className="text-sm text-muted-foreground leading-relaxed">Your library has been synchronized with the master inventory. 12 new product assets were detected and auto-tagged with AI metadata for immediate campaign generation.</p>
          <div className="mt-2">
            <button className="text-primary text-sm font-bold flex items-center gap-1 hover:underline">
              View Changelog
              <ExternalLink className="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
        
        <Card className="bg-muted-foreground bg-opacity-5 p-8 border border-primary/10 rounded-xl">
          <CardContent className="p-0 flex items-start justify-between">
            <div className="space-y-4">
              <h3 className="text-lg font-bold font-serif">Studio Integration</h3>
              <p className="text-sm text-muted-foreground">Ready to transform these raw assets into high-converting visual stories?</p>
              <Button className="bronze-gradient text-white px-8 h-12 shadow-lg shadow-primary/10 font-bold">
                Launch Studio Designer
              </Button>
            </div>
            <Sparkles className="w-12 h-12 text-primary/20" />
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

// --- Main App ---

export default function App() {
  const [view, setView] = useState<ViewState>('studio');

  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full bg-background font-sans selection:bg-accent/30">
        <TopNav currentView={view} setView={setView} />
        
        <div className="flex flex-1 pt-[72px]">
          <AppSidebar currentView={view} setView={setView} />
          
          <SidebarInset className="bg-background flex flex-col">
            <main className="flex-1 overflow-y-auto">
              {/* Sidebar trigger for mobile */}
              <div className="md:hidden p-4">
                <SidebarTrigger />
              </div>

              <AnimatePresence mode="wait">
                <motion.div
                  key={view}
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -10 }}
                  transition={{ duration: 0.2 }}
                >
                  {view === 'studio' ? <ArchitectView /> : <LibraryView />}
                </motion.div>
              </AnimatePresence>
            </main>
          </SidebarInset>
        </div>
      </div>
    </SidebarProvider>
  );
}
