import { Toaster } from "@/components/ui/sonner";
import { AssetLibraryPage } from "@/features/assets/pages/AssetLibraryPage";
import { AssetShowcasePage } from "@/features/assets/pages/AssetShowcasePage";
import { ArchitectPage } from "@/features/studio/pages/ArchitectPage";
import { BrowserRouter, Navigate, Route, Routes } from "react-router";
import { CenterLayout } from "./components/layout/CenterLayout";
import { MainLayout } from "./components/layout/MainLayout";
import { RequireAuth } from "./features/auth/components/RequireAuth";
import { RequireRole } from "./features/auth/components/RequireRole";
import SignInPage from "./features/auth/pages/SignInPage";
import SignUpPage from "./features/auth/pages/SignUpPage";
import { ROUTES } from "./lib/routes";

function App() {
  return (
    <BrowserRouter>
      <Toaster />
      <Routes>

        <Route element={<CenterLayout />}>
          <Route path="/sign-in" element={<SignInPage />} />
          <Route path="/sign-up" element={<SignUpPage />} />
        </Route>

        <Route element={<RequireAuth />}>

          <Route element={<RequireRole roles={["USER", "ADMIN"]} />}>
            <Route element={<MainLayout />}>
              <Route index element={<Navigate to={ROUTES.MAIN.ASSETS.LIST} replace />} />
              <Route path={ROUTES.MAIN.ASSETS.LIST} element={<AssetLibraryPage />} />
              <Route path={ROUTES.MAIN.ASSETS.SHOWCASE} element={<AssetShowcasePage />} />
              <Route path="/studio" element={<ArchitectPage />} />
            </Route>
          </Route>

          <Route element={<RequireRole roles={["ADMIN"]} />}>
            <Route path="/admin" element={<div>ADMIN</div>} />
          </Route>

        </Route>

      </Routes>
    </BrowserRouter>
  );
}

export default App;