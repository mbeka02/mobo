import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_authenticated/home")({
  component: HomePage,
});

function HomePage() {
  return (
    <main className="max-w-7xl mx-auto px-6 md:px-8 py-12">
      <div className="space-y-8">
        <div className="space-y-4">
          <h1 className="text-4xl md:text-5xl font-headline font-bold text-[var(--on-surface)] tracking-tight">
            Welcome to{" "}
            <span className="text-[var(--primary-container)]">Mobo</span>
          </h1>
          <p className="text-lg text-[var(--on-surface-variant)] max-w-xl">
            Your cinematic journey starts here. Browse movies, book seats, and
            discover curated experiences across Kenya.
          </p>
        </div>

        {/* Placeholder Cards */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 mt-12">
          {[
            {
              title: "Now Showing",
              desc: "Explore movies currently in theaters near you.",
              emoji: "🎬",
            },
            {
              title: "Coming Soon",
              desc: "Get early access to upcoming premieres and exclusives.",
              emoji: "🌟",
            },
            {
              title: "My Tickets",
              desc: "View your bookings, QR codes, and screening times.",
              emoji: "🎟️",
            },
          ].map((card) => (
            <div
              key={card.title}
              className="bg-[var(--surface-container-low)] rounded-2xl p-8 hover:bg-[var(--surface-container)] transition-colors"
            >
              <span className="text-4xl mb-4 block">{card.emoji}</span>
              <h3 className="text-xl font-headline font-bold mb-2">
                {card.title}
              </h3>
              <p className="text-[var(--on-surface-variant)] text-sm">
                {card.desc}
              </p>
            </div>
          ))}
        </div>
      </div>
    </main>
  );
}
