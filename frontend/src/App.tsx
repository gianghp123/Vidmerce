import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AssetLibraryPage } from "@/features/assets/pages/AssetLibraryPage";
import { VideoLibraryPage } from "@/features/videos/pages/VideoLibraryPage";
import { VideoBuilderPage } from "@/features/videos/pages/VideoBuilderPage";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/library" replace />} />
        <Route path="/library" element={<AssetLibraryPage />} />
        <Route path="/videos" element={<VideoLibraryPage />} />
        <Route path="/videos/create" element={<VideoBuilderPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;