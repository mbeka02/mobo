import { createFileRoute } from "@tanstack/react-router";
import { Film } from "lucide-react";

export const Route = createFileRoute("/portal/movies")({
  component: MoviesPage,
});

function MoviesPage() {
  return (
    <div className="space-y-6 rise-in">
      <header className="space-y-2">
        <h2
          className="font-headline text-4xl font-semibold tracking-tight"
          style={{ color: "var(--on-surface)" }}
        >
          Movies
        </h2>
        <p style={{ color: "var(--on-surface-variant)" }}>
          Manage your movie catalog.
        </p>
      </header>
      <div
        className="rounded-xl p-16 flex flex-col items-center justify-center gap-4"
        style={{ backgroundColor: "var(--surface-container-low)" }}
      >
        <Film className="w-12 h-12 opacity-30" style={{ color: "var(--on-surface-variant)" }} />
        <p className="text-lg font-medium" style={{ color: "var(--on-surface-variant)" }}>
          Movie management coming soon
        </p>
      </div>
    </div>
  );
}
