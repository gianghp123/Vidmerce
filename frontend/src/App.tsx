import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AssetLibraryPage } from "@/features/assets/pages/AssetLibraryPage";
import { AssetShowcasePage } from "@/features/assets/pages/AssetShowcasePage";
import { VideoLibraryPage } from "@/features/videos/pages/VideoLibraryPage";
import { VideoBuilderPage } from "@/features/videos/pages/VideoBuilderPage";
import { Toaster } from "@/components/ui/sonner";
import { PageLayout } from "@/components/custom/PageLayout";

function App() {
  return (
    <BrowserRouter>
      <Toaster />
      <PageLayout>
        <Routes>
          <Route path="/" element={<Navigate to="/assets" replace />} />
          <Route path="/assets" element={<AssetLibraryPage />} />
          <Route path="/assets/:id" element={<AssetShowcasePage />} />
          <Route path="/videos" element={<VideoLibraryPage />} />
          <Route path="/videos/create" element={<VideoBuilderPage />} />
        </Routes>
      </PageLayout>
    </BrowserRouter>
  );
}

export default App;