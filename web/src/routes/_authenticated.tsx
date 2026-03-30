import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import AppNavbar from "../components/AppNavbar";
import { getCurrentUser } from "../lib/auth";

export const Route = createFileRoute("/_authenticated")({
  beforeLoad: async () => {
    try {
      const user = await getCurrentUser();
      return { user };
    } catch {
      throw redirect({ to: "/login" });
    }
  },
  component: AuthenticatedLayout,
});

function AuthenticatedLayout() {
  return (
    <div className="min-h-screen bg-[var(--surface)]">
      <AppNavbar />
      <Outlet />
    </div>
  );
}
