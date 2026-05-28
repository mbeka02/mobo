import { MapPin, Plus } from "lucide-react";
import { Link } from "@tanstack/react-router";
import {
  type VenueResponse,
  getVenueGradient,
  getInitials,
} from "./types";

/* ── Featured Venue Card (first venue, wide) ── */
export function FeaturedVenueCard({ venue }: { venue: VenueResponse }) {
  return (
    <div
      className="rounded-xl overflow-hidden grid grid-cols-1 md:grid-cols-12 gap-0"
      style={{ backgroundColor: "var(--surface-container-lowest)" }}
    >
      {/* Image area */}
      <div
        className="md:col-span-5 h-48 md:h-auto relative flex items-center justify-center"
        style={{ background: getVenueGradient(0) }}
      >
        <span className="absolute top-4 left-4 px-3 py-1 rounded-full text-xs font-bold tracking-wider uppercase bg-black/30 text-white backdrop-blur-sm">
          ● Active
        </span>
        <span className="text-white/30 font-headline text-6xl font-bold select-none">
          {getInitials(venue.name)}
        </span>
      </div>

      {/* Details */}
      <div className="md:col-span-7 p-6 md:p-8 flex flex-col justify-between gap-6">
        <div>
          <h3
            className="font-headline text-2xl font-bold tracking-tight"
            style={{ color: "var(--on-surface)" }}
          >
            {venue.name}
          </h3>
          <div
            className="flex items-center gap-1.5 mt-2 text-sm"
            style={{ color: "var(--on-surface-variant)" }}
          >
            <MapPin className="w-3.5 h-3.5" />
            {venue.address}, {venue.city}
          </div>
        </div>

        {/* Stats */}
        <div className="flex gap-6">
          <div
            className="rounded-lg px-5 py-3"
            style={{ backgroundColor: "var(--surface-container)" }}
          >
            <span
              className="font-label text-[10px] tracking-widest uppercase block"
              style={{ color: "var(--on-surface-variant)" }}
            >
              Capacity
            </span>
            <span
              className="font-headline text-xl font-bold"
              style={{ color: "var(--on-surface)" }}
            >
              {venue.total_seats.toLocaleString()}
            </span>
          </div>
          <div
            className="rounded-lg px-5 py-3"
            style={{ backgroundColor: "var(--surface-container)" }}
          >
            <span
              className="font-label text-[10px] tracking-widest uppercase block"
              style={{ color: "var(--on-surface-variant)" }}
            >
              Status
            </span>
            <span
              className="font-headline text-xl font-bold"
              style={{ color: "var(--primary)" }}
            >
              Optimal
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}

/* ── Standard Venue Card ── */
export function VenueCard({
  venue,
  index,
}: {
  venue: VenueResponse;
  index: number;
}) {
  return (
    <div
      className="rounded-xl overflow-hidden flex flex-col transition-transform hover:-translate-y-1 duration-300"
      style={{ backgroundColor: "var(--surface-container-lowest)" }}
    >
      {/* Image area */}
      <div
        className="h-40 relative flex items-center justify-center"
        style={{ background: getVenueGradient(index) }}
      >
        <span className="absolute top-3 left-3 px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-black/30 text-white backdrop-blur-sm">
          ● Active
        </span>
        <span className="text-white/30 font-headline text-5xl font-bold select-none">
          {getInitials(venue.name)}
        </span>
      </div>

      {/* Details */}
      <div className="p-5 flex flex-col flex-1">
        <h4
          className="font-headline text-lg font-bold"
          style={{ color: "var(--on-surface)" }}
        >
          {venue.name}
        </h4>
        <div
          className="flex items-center gap-1 mt-1.5 text-sm"
          style={{ color: "var(--on-surface-variant)" }}
        >
          <MapPin className="w-3 h-3 shrink-0" />
          {venue.address}, {venue.city}
        </div>

        <div className="flex gap-4 mt-4 pt-4">
          <div>
            <span
              className="font-label text-[10px] tracking-widest uppercase block"
              style={{ color: "var(--on-surface-variant)" }}
            >
              Capacity
            </span>
            <span
              className="font-headline text-base font-bold"
              style={{ color: "var(--on-surface)" }}
            >
              {venue.total_seats.toLocaleString()}
            </span>
          </div>
          <div>
            <span
              className="font-label text-[10px] tracking-widest uppercase block"
              style={{ color: "var(--on-surface-variant)" }}
            >
              Status
            </span>
            <span
              className="font-headline text-base font-bold"
              style={{ color: "var(--primary)" }}
            >
              Optimal
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}

/* ── Onboard CTA Card ── */
export function OnboardCard() {
  return (
    <Link
      to="/portal/venues/new"
      className="rounded-xl p-6 flex flex-col items-center justify-center gap-4 min-h-[16rem] cursor-pointer transition-all hover:opacity-80 active:scale-[0.98] no-underline"
      style={{
        backgroundColor: "transparent",
        border:
          "2px dashed color-mix(in srgb, var(--outline-variant) 40%, transparent)",
      }}
    >
      <div
        className="w-12 h-12 rounded-full flex items-center justify-center"
        style={{
          backgroundColor: "var(--surface-container)",
        }}
      >
        <Plus
          className="w-5 h-5"
          style={{ color: "var(--on-surface-variant)" }}
        />
      </div>
      <div className="text-center">
        <h4
          className="font-headline font-bold text-base"
          style={{ color: "var(--on-surface)" }}
        >
          Onboard Venue
        </h4>
        <p
          className="text-sm mt-1 max-w-[12rem]"
          style={{ color: "var(--on-surface-variant)" }}
        >
          Integrate a new physical location into the Mobo network.
        </p>
      </div>
    </Link>
  );
}

/* ── Stat Card ── */
export function StatCard({
  icon,
  label,
  value,
  sub,
}: {
  icon: React.ReactNode;
  label: string;
  value: string;
  sub: string;
}) {
  return (
    <div
      className="rounded-xl p-6 flex flex-col justify-center"
      style={{ backgroundColor: "var(--surface-container-lowest)" }}
    >
      <div className="flex items-center gap-2 mb-3">
        <span style={{ color: "var(--on-surface-variant)" }}>{icon}</span>
        <span
          className="font-label text-xs tracking-widest uppercase"
          style={{ color: "var(--on-surface-variant)" }}
        >
          {label}
        </span>
      </div>
      <h3
        className="font-headline text-3xl font-bold tracking-tight"
        style={{ color: "var(--on-surface)" }}
      >
        {value}
      </h3>
      <p
        className="text-sm mt-2"
        style={{ color: "var(--on-surface-variant)" }}
      >
        {sub}
      </p>
    </div>
  );
}
