import { Outlet } from "react-router";

export function CenterLayout() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen">
      <Outlet />
    </div>
  );
}