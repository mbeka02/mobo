import { createFileRoute } from "@tanstack/react-router";
import { CalendarDays } from "lucide-react";

export const Route = createFileRoute("/portal/scheduling")({
  component: SchedulingPage,
});

function SchedulingPage() {
  return (
    <div className="space-y-6 rise-in">
      <header className="space-y-2">
        <h2
          className="font-headline text-4xl font-semibold tracking-tight"
          style={{ color: "var(--on-surface)" }}
        >
          Scheduling
        </h2>
        <p style={{ color: "var(--on-surface-variant)" }}>
          Manage showtimes and screenings.
        </p>
      </header>
      <div
        className="rounded-xl p-16 flex flex-col items-center justify-center gap-4"
        style={{ backgroundColor: "var(--surface-container-low)" }}
      >
        <CalendarDays className="w-12 h-12 opacity-30" style={{ color: "var(--on-surface-variant)" }} />
        <p className="text-lg font-medium" style={{ color: "var(--on-surface-variant)" }}>
          Showtime scheduling coming soon
        </p>
      </div>
    </div>
  );
}
