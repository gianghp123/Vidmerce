import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AssetLibraryPage } from "@/features/assets/pages/AssetLibraryPage";
import { VideoLibraryPage } from "@/features/videos/pages/VideoLibraryPage";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Navigate to="/library" replace />} />
        <Route path="/library" element={<AssetLibraryPage />} />
        <Route path="/videos" element={<VideoLibraryPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;