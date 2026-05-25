import {
  createFileRoute,
  Outlet,
  redirect,
} from "@tanstack/react-router";
import { getCurrentUser } from "../../lib/auth";
import { SidebarProvider } from "../../components/ui/sidebar";
import AdminSidebar from "../../components/AdminSidebar";
import AdminTopBar from "../../components/AdminTopBar";

export const Route = createFileRoute("/portal")({
  beforeLoad: async () => {
    try {
      const user = await getCurrentUser();
      if (user.role !== "admin") {
        throw redirect({ to: "/home" });
      }
      return { user };
    } catch (e) {
      // If it's already a redirect, re-throw it
      if (e && typeof e === "object" && "to" in e) throw e;
      throw redirect({ to: "/login" });
    }
  },
  component: PortalLayout,
});

function PortalLayout() {
  return (
    <SidebarProvider>
      <div className="flex min-h-screen w-full" style={{ backgroundColor: "var(--surface)" }}>
        <AdminSidebar />
        <main className="flex-1 flex flex-col min-w-0">
          <AdminTopBar />
          <div className="flex-1 px-8 md:px-12 py-8 md:py-10">
            <Outlet />
          </div>
        </main>
      </div>
    </SidebarProvider>
  );
}
