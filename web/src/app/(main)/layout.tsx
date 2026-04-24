import { Sidebar } from "@/components/common/Sidebar";
import { Topbar } from "@/components/common/Topbar";
import { RedirectToSignIn, Show } from "@clerk/nextjs";


export default function MainLayout({ children }: { children: React.ReactNode }) {

  return (
    <div className="flex min-h-screen">
      <Sidebar />

      <div className="flex flex-col w-full">
        <Topbar />
        <Show when="signed-in">
          {children}
        </Show>
        <Show when="signed-out">
          <RedirectToSignIn />
        </Show>
      </div>
    </div>
  );
}