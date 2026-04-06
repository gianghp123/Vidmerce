import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AssetLibraryPage } from "@/features/assets/pages/AssetLibraryPage";
import { AssetShowcasePage } from "@/features/assets/pages/AssetShowcasePage";
import { VideoLibraryPage } from "@/features/videos/pages/VideoLibraryPage";
import { VideoBuilderPage } from "@/features/videos/pages/VideoBuilderPage";
import { Toaster } from "@/components/ui/sonner";

function App() {
  return (
    <BrowserRouter>
      <Toaster />
      <Routes>
        <Route path="/" element={<Navigate to="/library" replace />} />
        <Route path="/library" element={<AssetLibraryPage />} />
        <Route path="/assets/:id" element={<AssetShowcasePage />} />
        <Route path="/videos" element={<VideoLibraryPage />} />
        <Route path="/videos/create" element={<VideoBuilderPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;