import type { UserRole } from "@/lib/enums/user-role.enum";
import { ROUTES } from "@/lib/routes";
import { useUser } from "@clerk/react";
import { Navigate } from "react-router-dom";

export function HomeRedirect() {
  const { user } = useUser();

  const role = user?.publicMetadata?.role as UserRole;

  if (role === "ADMIN") return <Navigate to={ROUTES.ADMIN.MAIN} replace />;
  return <Navigate to={ROUTES.MAIN.ASSETS.LIST} replace />;
}