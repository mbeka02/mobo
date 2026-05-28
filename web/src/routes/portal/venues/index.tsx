import { createFileRoute, Link } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { Building2, Users, MapPin, Plus } from "lucide-react";
import { api } from "../../../services/api";
import type { VenuesAPIResponse } from "../../../components/venues/types";
import {
  FeaturedVenueCard,
  VenueCard,
  OnboardCard,
  StatCard,
} from "../../../components/venues/VenueCards";

export const Route = createFileRoute("/portal/venues/")({
  component: VenueListPage,
});

function VenueListPage() {
  const { data, isLoading } = useQuery({
    queryKey: ["venues"],
    queryFn: async () => {
      const res = await api.get<VenuesAPIResponse>("/venues");
      return res.data ?? [];
    },
  });

  const venues = data ?? [];
  const totalCapacity = venues.reduce((sum, v) => sum + v.total_seats, 0);

  if (isLoading) {
    return <VenueListSkeleton />;
  }

  return (
    <div className="space-y-10 rise-in">
      {/* ── Header ── */}
      <header className="flex flex-col md:flex-row md:items-end md:justify-between gap-6">
        <div className="space-y-3 max-w-2xl">
          <span
            className="font-label text-xs font-bold tracking-[0.16em] uppercase"
            style={{ color: "var(--primary)" }}
          >
            Infrastructure
          </span>
          <h2
            className="font-headline text-4xl md:text-5xl font-semibold tracking-tight leading-[1.05]"
            style={{ color: "var(--on-surface)" }}
          >
            Digital Hearth
            <br />
            Locations.
          </h2>
          <p
            className="text-lg leading-relaxed"
            style={{ color: "var(--on-surface-variant)" }}
          >
            Manage cinematic spaces, coordinate screen allocations, and oversee
            the operational heartbeat of the Mobo network.
          </p>
        </div>

        <Link
          to="/portal/venues/new"
          className="self-start md:self-auto flex items-center gap-2 px-6 py-3.5 rounded-xl font-body font-bold text-sm transition-all active:scale-95 cursor-pointer hover:shadow-lg no-underline"
          style={{
            background:
              "linear-gradient(135deg, var(--primary), var(--primary-container))",
            color: "var(--on-primary)",
          }}
        >
          <Building2 className="w-4 h-4" />
          Add New Venue
        </Link>
      </header>

      {/* ── Stats Row ── */}
      <section className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <StatCard
          icon={<Building2 className="w-4 h-4" />}
          label="Total Venues"
          value={String(venues.length)}
          sub={venues.length > 0 ? "Across all regions" : "No venues yet"}
        />
        <StatCard
          icon={<Users className="w-4 h-4" />}
          label="Total Capacity"
          value={totalCapacity.toLocaleString()}
          sub="Combined seating"
        />
      </section>

      {/* ── Venue Grid ── */}
      {venues.length === 0 ? (
        <EmptyState />
      ) : (
        <section className="space-y-6">
          <FeaturedVenueCard venue={venues[0]} />

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {venues.slice(1).map((venue, i) => (
              <VenueCard key={venue.id} venue={venue} index={i + 1} />
            ))}
            <OnboardCard />
          </div>
        </section>
      )}
    </div>
  );
}

/* ── Empty State ── */
function EmptyState() {
  return (
    <div
      className="rounded-xl p-16 flex flex-col items-center justify-center gap-6 text-center"
      style={{ backgroundColor: "var(--surface-container-low)" }}
    >
      <div
        className="w-20 h-20 rounded-full flex items-center justify-center"
        style={{ backgroundColor: "var(--surface-container)" }}
      >
        <MapPin
          className="w-10 h-10 opacity-40"
          style={{ color: "var(--on-surface-variant)" }}
        />
      </div>
      <div className="space-y-2">
        <h3
          className="font-headline text-2xl font-bold"
          style={{ color: "var(--on-surface)" }}
        >
          No venues registered yet
        </h3>
        <p
          className="text-base max-w-md"
          style={{ color: "var(--on-surface-variant)" }}
        >
          Start building your Digital Hearth network by adding your first
          cinematic venue.
        </p>
      </div>
      <Link
        to="/portal/venues/new"
        className="flex items-center gap-2 px-6 py-3 rounded-xl font-body font-bold text-sm transition-all active:scale-95 cursor-pointer no-underline"
        style={{
          background:
            "linear-gradient(135deg, var(--primary), var(--primary-container))",
          color: "var(--on-primary)",
        }}
      >
        <Plus className="w-4 h-4" />
        Add Your First Venue
      </Link>
    </div>
  );
}

/* ── Skeleton ── */
function VenueListSkeleton() {
  return (
    <div className="space-y-10 animate-pulse">
      <div className="space-y-3">
        <div
          className="h-4 w-28 rounded"
          style={{ backgroundColor: "var(--surface-container-high)" }}
        />
        <div
          className="h-14 w-80 rounded-lg"
          style={{ backgroundColor: "var(--surface-container-high)" }}
        />
        <div
          className="h-5 w-[28rem] rounded"
          style={{ backgroundColor: "var(--surface-container-high)" }}
        />
      </div>
      <div className="grid grid-cols-3 gap-6">
        {[1, 2, 3].map((i) => (
          <div
            key={i}
            className="h-32 rounded-xl"
            style={{ backgroundColor: "var(--surface-container-low)" }}
          />
        ))}
      </div>
      <div
        className="h-56 rounded-xl"
        style={{ backgroundColor: "var(--surface-container-low)" }}
      />
    </div>
  );
}
