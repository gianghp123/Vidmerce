import { useUser } from "@clerk/react";
import { Navigate, Outlet } from "react-router";
import type { UserRole } from "@/lib/enums/user-role.enum";

type Props = {
  roles: UserRole[];
  redirectTo?: string;
};

export function RequireRole({ roles, redirectTo = "/" }: Props) {
  const { user, isLoaded } = useUser();

  if (!isLoaded) return null;

  const role = user?.publicMetadata?.role as UserRole | undefined;

  if (!role || !roles.includes(role)) {
    return <Navigate to={redirectTo} replace />;
  }

  return <Outlet />;
}