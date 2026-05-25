import {
  LayoutDashboard,
  Film,
  MapPin,
  CalendarDays,
  Settings,
  LogOut,
  Plus,
} from "lucide-react";
import { Link, useLocation } from "@tanstack/react-router";
import { logout } from "../lib/auth";
import { toast } from "sonner";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "./ui/sidebar";

const navItems = [
  { title: "Dashboard", icon: LayoutDashboard, href: "/portal/dashboard" },
  { title: "Movies", icon: Film, href: "/portal/movies" },
  { title: "Venues", icon: MapPin, href: "/portal/venues" },
  { title: "Scheduling", icon: CalendarDays, href: "/portal/scheduling" },
];

const footerItems = [
  { title: "Settings", icon: Settings, href: "/portal/settings" },
];

export default function AdminSidebar() {
  const location = useLocation();

  const handleLogout = async () => {
    try {
      await logout();
      window.location.href = "/login";
    } catch {
      toast.error("Failed to log out. Please try again.");
    }
  };

  return (
    <Sidebar
      className="border-r-0"
      style={{
        backgroundColor: "var(--surface-container-low)",
        color: "var(--on-surface)",
      }}
    >
      <SidebarHeader className="p-4 pb-0">
        {/* Brand */}
        <div className="flex items-center gap-3 px-2 mb-6">
          <div
            className="w-10 h-10 rounded-full flex items-center justify-center font-headline font-bold text-sm"
            style={{
              backgroundColor: "var(--surface-container-lowest)",
              color: "var(--primary)",
            }}
          >
            MA
          </div>
          <div>
            <h1
              className="font-headline text-lg font-bold tracking-tight"
              style={{ color: "var(--primary)" }}
            >
              Mobo Admin
            </h1>
            <p
              className="text-xs"
              style={{ color: "var(--on-surface-variant)" }}
            >
              Digital Hearth Portal
            </p>
          </div>
        </div>

        {/* New Event CTA */}
        <button
          className="w-full py-3 px-4 rounded-xl font-body font-bold text-sm flex items-center justify-center gap-2 transition-transform active:scale-95 mb-4"
          style={{
            background: "linear-gradient(135deg, var(--primary), var(--primary-container))",
            color: "var(--on-primary)",
          }}
        >
          <Plus className="w-4 h-4" />
          New Event
        </button>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              {navItems.map((item) => {
                const isActive = location.pathname === item.href;
                return (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton
                      asChild
                      isActive={isActive}
                      className="h-11 rounded-lg transition-all duration-300"
                      style={
                        isActive
                          ? {
                              backgroundColor: "var(--primary-container)",
                              color: "var(--on-primary-container)",
                            }
                          : {
                              color: "var(--on-surface-variant)",
                            }
                      }
                    >
                      <Link to={item.href}>
                        <item.icon className="w-5 h-5" />
                        <span className="font-body text-base">
                          {item.title}
                        </span>
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                );
              })}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter className="p-4 pt-0">
        <SidebarMenu>
          {footerItems.map((item) => {
            const isActive = location.pathname === item.href;
            return (
              <SidebarMenuItem key={item.title}>
                <SidebarMenuButton
                  asChild
                  isActive={isActive}
                  className="h-11 rounded-lg transition-all duration-300"
                  style={{ color: "var(--on-surface-variant)" }}
                >
                  <Link to={item.href}>
                    <item.icon className="w-5 h-5" />
                    <span className="font-body text-base">{item.title}</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            );
          })}
          <SidebarMenuItem>
            <SidebarMenuButton
              onClick={handleLogout}
              className="h-11 rounded-lg transition-all duration-300 cursor-pointer"
              style={{ color: "var(--on-surface-variant)" }}
            >
              <LogOut className="w-5 h-5" />
              <span className="font-body text-base">Logout</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  );
}
