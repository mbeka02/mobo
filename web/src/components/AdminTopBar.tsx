import { Search, Bell, HelpCircle } from "lucide-react";
import ThemeToggle from "./ThemeToggle";

export default function AdminTopBar() {
  return (
    <header
      className="sticky top-0 z-40 flex justify-between items-center px-8 h-16 transition-all"
      style={{
        backgroundColor: "color-mix(in srgb, var(--surface) 70%, transparent)",
        backdropFilter: "blur(24px)",
        WebkitBackdropFilter: "blur(24px)",
      }}
    >
      <div />

      <div className="flex items-center gap-3">
        {/* Search */}
        <div className="relative hidden md:block">
          <Search
            className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4"
            style={{ color: "var(--on-surface-variant)" }}
          />
          <input
            type="text"
            placeholder="Search..."
            className="border-0 text-sm rounded-full pl-10 pr-4 py-2 w-48 transition-all focus:w-64 focus:outline-none focus:ring-2"
            style={{
              backgroundColor: "var(--surface-container-high)",
              color: "var(--on-surface)",
              fontFamily: "'Manrope', sans-serif",
            }}
          />
        </div>

        {/* Actions */}
        <button
          className="w-10 h-10 rounded-full flex items-center justify-center transition-colors duration-300 hover:opacity-80"
          style={{ color: "var(--on-surface-variant)" }}
        >
          <Bell className="w-5 h-5" />
        </button>

        <button
          className="w-10 h-10 rounded-full flex items-center justify-center transition-colors duration-300 hover:opacity-80"
          style={{ color: "var(--on-surface-variant)" }}
        >
          <HelpCircle className="w-5 h-5" />
        </button>

        <ThemeToggle />
      </div>
    </header>
  );
}
