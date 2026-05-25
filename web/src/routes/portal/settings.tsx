import { createFileRoute } from "@tanstack/react-router";
import { Settings } from "lucide-react";

export const Route = createFileRoute("/portal/settings")({
  component: SettingsPage,
});

function SettingsPage() {
  return (
    <div className="space-y-6 rise-in">
      <header className="space-y-2">
        <h2
          className="font-headline text-4xl font-semibold tracking-tight"
          style={{ color: "var(--on-surface)" }}
        >
          Settings
        </h2>
        <p style={{ color: "var(--on-surface-variant)" }}>
          Configure your admin preferences.
        </p>
      </header>
      <div
        className="rounded-xl p-16 flex flex-col items-center justify-center gap-4"
        style={{ backgroundColor: "var(--surface-container-low)" }}
      >
        <Settings className="w-12 h-12 opacity-30" style={{ color: "var(--on-surface-variant)" }} />
        <p className="text-lg font-medium" style={{ color: "var(--on-surface-variant)" }}>
          Settings panel coming soon
        </p>
      </div>
    </div>
  );
}
