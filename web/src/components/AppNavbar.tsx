import { Link, useNavigate } from "@tanstack/react-router";
import { Bell, LayoutGrid, LogOut, User } from "lucide-react";
import { useState } from "react";
import { logout } from "../lib/auth";
import ThemeToggle from "./ThemeToggle";

export default function AppNavbar() {
  const navigate = useNavigate();
  const [showMenu, setShowMenu] = useState(false);

  const handleLogout = async () => {
    await logout();
    navigate({ to: "/login" });
  };

  return (
    <nav className="sticky top-0 z-50 bg-[#131313]/80 backdrop-blur-xl">
      <div className="flex justify-between items-center px-6 md:px-8 py-3 max-w-7xl mx-auto">
        {/* Left: Logo + Nav */}
        <div className="flex items-center gap-8 md:gap-10">
          <Link to="/home" className="no-underline">
            <span className="text-xl font-black tracking-tighter text-[var(--primary-container)] font-headline uppercase">
              Mobo
            </span>
          </Link>
          <div className="hidden md:flex gap-6 items-center">
            {[
              { label: "Cinemas", href: "/home" },
              { label: "Offers", href: "/home" },
              { label: "Magazine", href: "/home" },
            ].map((link) => (
              <Link
                key={link.label}
                to={link.href}
                className="font-label text-sm font-medium text-[#fff9ef]/70 hover:text-[#fff9ef] transition-colors no-underline"
              >
                {link.label}
              </Link>
            ))}
          </div>
        </div>

        {/* Right: Icons */}
        <div className="flex items-center gap-3 md:gap-4">
          <button
            className="text-[#fff9ef]/60 hover:text-[#fff9ef] transition-colors p-2"
            aria-label="Notifications"
          >
            <Bell size={20} />
          </button>
          <button
            className="hidden sm:flex text-[#fff9ef]/60 hover:text-[#fff9ef] transition-colors p-2"
            aria-label="Dashboard"
          >
            <LayoutGrid size={20} />
          </button>
          <ThemeToggle />

          {/* User Avatar / Menu */}
          <div className="relative">
            <button
              onClick={() => setShowMenu(!showMenu)}
              className="w-9 h-9 rounded-full bg-[var(--primary-container)]/30 border-2 border-[var(--primary-container)]/50 flex items-center justify-center text-[var(--primary-container)] hover:bg-[var(--primary-container)]/50 transition-colors"
              aria-label="User menu"
            >
              <User size={16} />
            </button>

            {showMenu && (
              <div className="absolute right-0 mt-2 w-48 bg-[var(--surface-container-high)] rounded-xl shadow-2xl overflow-hidden fade-in z-50">
                <Link
                  to="/home"
                  className="flex items-center gap-3 px-4 py-3 text-sm text-[var(--on-surface)] hover:bg-[var(--surface-container-highest)] transition-colors no-underline"
                  onClick={() => setShowMenu(false)}
                >
                  <User size={16} />
                  Profile
                </Link>
                <button
                  onClick={handleLogout}
                  className="flex items-center gap-3 px-4 py-3 text-sm text-[var(--error)] hover:bg-[var(--surface-container-highest)] transition-colors w-full text-left"
                >
                  <LogOut size={16} />
                  Sign Out
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
