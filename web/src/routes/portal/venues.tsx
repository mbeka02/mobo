import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  MapPin,
  Plus,
  Building2,
  Users,
  CheckCircle2,
  ArrowRight,
  ArrowLeft,
  Info,
} from "lucide-react";
import { toast } from "sonner";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../../components/ui/form";
import { Input } from "../../components/ui/input";
import { api, isAPIError } from "../../services/api";

export const Route = createFileRoute("/portal/venues")({
  component: VenuesPage,
});

/* ─── Types ─── */
interface VenueResponse {
  id: number;
  name: string;
  address: string;
  city: string;
  total_seats: number;
  created_at: string;
  updated_at?: string;
}

interface VenuesAPIResponse {
  status: number;
  message: string;
  data: VenueResponse[] | null;
}

interface CreateVenueAPIResponse {
  status: number;
  message: string;
  data: VenueResponse;
}

/* ─── Validation ─── */
const createVenueSchema = z.object({
  name: z.string().min(2, "Venue name must be at least 2 characters"),
  city: z.string().min(2, "City is required"),
  address: z.string().min(2, "Address is required"),
  total_seats: z.coerce
    .number({ invalid_type_error: "Must be a number" })
    .int("Must be a whole number")
    .min(1, "Must have at least 1 seat"),
});

type CreateVenueValues = z.infer<typeof createVenueSchema>;

/* ─── Helpers ─── */
const VENUE_GRADIENTS = [
  "linear-gradient(135deg, #AB3600 0%, #FF5F1F 100%)",
  "linear-gradient(135deg, #5C1A00 0%, #9B4425 100%)",
  "linear-gradient(135deg, #004B70 0%, #006493 100%)",
  "linear-gradient(135deg, #7C2E0F 0%, #FE916B 100%)",
  "linear-gradient(135deg, #003350 0%, #8DCDFF 100%)",
];

function getVenueGradient(index: number) {
  return VENUE_GRADIENTS[index % VENUE_GRADIENTS.length];
}

function getInitials(name: string) {
  return name
    .split(" ")
    .map((w) => w[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);
}

/* ─── Page ─── */
function VenuesPage() {
  const [view, setView] = useState<"list" | "add">("list");

  return view === "list" ? (
    <VenueListView onAdd={() => setView("add")} />
  ) : (
    <VenueAddView onCancel={() => setView("list")} />
  );
}

/* ════════════════════════════════════════════
   LIST VIEW
   ════════════════════════════════════════════ */
function VenueListView({ onAdd }: { onAdd: () => void }) {
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

        <button
          onClick={onAdd}
          className="self-start md:self-auto flex items-center gap-2 px-6 py-3.5 rounded-xl font-body font-bold text-sm transition-all active:scale-95 cursor-pointer hover:shadow-lg"
          style={{
            background:
              "linear-gradient(135deg, var(--primary), var(--primary-container))",
            color: "var(--on-primary)",
          }}
        >
          <Building2 className="w-4 h-4" />
          Add New Venue
        </button>
      </header>

      {/* ── Stats Row ── */}
      <section className="grid grid-cols-1 md:grid-cols-3 gap-6">
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
        <div
          className="rounded-xl p-6 flex flex-col justify-center"
          style={{ backgroundColor: "var(--surface-container-lowest)" }}
        >
          <div className="flex items-center gap-2 mb-3">
            <CheckCircle2
              className="w-4 h-4"
              style={{ color: "var(--primary)" }}
            />
            <span
              className="font-label text-xs tracking-widest uppercase font-bold"
              style={{ color: "var(--primary)" }}
            >
              System Status
            </span>
          </div>
          <h3
            className="font-headline text-2xl font-bold tracking-tight"
            style={{ color: "var(--on-surface)" }}
          >
            All systems
            <br />
            operational.
          </h3>
        </div>
      </section>

      {/* ── Venue Grid ── */}
      {venues.length === 0 ? (
        <EmptyState onAdd={onAdd} />
      ) : (
        <section className="space-y-6">
          {/* Featured venue (first) */}
          <FeaturedVenueCard venue={venues[0]} />

          {/* Remaining venues + onboard CTA */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {venues.slice(1).map((venue, i) => (
              <VenueCard key={venue.id} venue={venue} index={i + 1} />
            ))}
            <OnboardCard onAdd={onAdd} />
          </div>
        </section>
      )}
    </div>
  );
}

/* ── Stat Card ── */
function StatCard({
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

/* ── Featured Venue Card (first venue, wide) ── */
function FeaturedVenueCard({ venue }: { venue: VenueResponse }) {
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
function VenueCard({
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
function OnboardCard({ onAdd }: { onAdd: () => void }) {
  return (
    <button
      onClick={onAdd}
      className="rounded-xl p-6 flex flex-col items-center justify-center gap-4 min-h-[16rem] cursor-pointer transition-all hover:opacity-80 active:scale-[0.98]"
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
    </button>
  );
}

/* ── Empty State ── */
function EmptyState({ onAdd }: { onAdd: () => void }) {
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
      <button
        onClick={onAdd}
        className="flex items-center gap-2 px-6 py-3 rounded-xl font-body font-bold text-sm transition-all active:scale-95 cursor-pointer"
        style={{
          background:
            "linear-gradient(135deg, var(--primary), var(--primary-container))",
          color: "var(--on-primary)",
        }}
      >
        <Plus className="w-4 h-4" />
        Add Your First Venue
      </button>
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

/* ════════════════════════════════════════════
   ADD VIEW
   ════════════════════════════════════════════ */
function VenueAddView({ onCancel }: { onCancel: () => void }) {
  const queryClient = useQueryClient();

  const form = useForm<CreateVenueValues>({
    resolver: zodResolver(createVenueSchema),
    defaultValues: {
      name: "",
      city: "",
      address: "",
      total_seats: 0,
    },
  });

  const mutation = useMutation({
    mutationFn: async (values: CreateVenueValues) => {
      const res = await api.post<CreateVenueAPIResponse>(
        "/admin/venues",
        values
      );
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["venues"] });
      toast.success("Venue registered successfully!");
      onCancel();
    },
    onError: (err) => {
      if (isAPIError(err)) {
        toast.error(err.detail || err.message);
      } else {
        toast.error("Failed to register venue. Please try again.");
      }
    },
  });

  const onSubmit = (values: CreateVenueValues) => {
    mutation.mutate(values);
  };

  /* Shared input classname matching design.md input spec */
  const inputClassName =
    "h-auto bg-[var(--surface-container-high)] text-[var(--on-surface)] rounded-xl px-4 py-4 text-sm font-body border-none focus-visible:ring-2 focus-visible:ring-[var(--primary)] placeholder:text-[var(--on-surface-variant)]/50";

  return (
    <div className="rise-in flex flex-col min-h-[calc(100vh-8rem)]">
      {/* ── Back Button ── */}
      <button
        onClick={onCancel}
        className="self-start flex items-center gap-2 mb-6 text-sm font-body font-medium cursor-pointer transition-colors hover:opacity-70"
        style={{ color: "var(--on-surface-variant)" }}
      >
        <ArrowLeft className="w-4 h-4" />
        Back to Venues
      </button>

      {/* ── Header ── */}
      <header className="space-y-3 mb-10 max-w-2xl">
        <h2
          className="font-headline text-4xl md:text-5xl font-bold tracking-tight leading-[1.05]"
          style={{ color: "var(--on-surface)" }}
        >
          Establish a new
          <br />
          <span
            className="italic"
            style={{ color: "var(--primary-container)" }}
          >
            cinematic anchor.
          </span>
        </h2>
        <p
          className="text-base leading-relaxed"
          style={{ color: "var(--on-surface-variant)" }}
        >
          Add a new theater location to the Digital Hearth network. Define its
          capacity and geographical footprint.
        </p>
      </header>

      {/* ── Form ── */}
      <Form {...form}>
        <form
          onSubmit={form.handleSubmit(onSubmit)}
          className="flex-1 flex flex-col"
        >
          <div className="flex-1 space-y-8 max-w-2xl">
            {/* Section: Core Identity */}
            <section
              className="rounded-xl p-6 md:p-8 space-y-6"
              style={{ backgroundColor: "var(--surface-container-lowest)" }}
            >
              <div className="flex items-center gap-3 mb-2">
                <Building2
                  className="w-5 h-5"
                  style={{ color: "var(--on-surface)" }}
                />
                <h3
                  className="font-headline text-lg font-bold"
                  style={{ color: "var(--on-surface)" }}
                >
                  Core Identity
                </h3>
              </div>

              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-[var(--on-surface-variant)] font-label text-xs font-bold tracking-[0.08em] uppercase">
                      Venue Name
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder="e.g. The Grand Majestic"
                        className={inputClassName}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <FormField
                  control={form.control}
                  name="city"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-[var(--on-surface-variant)] font-label text-xs font-bold tracking-[0.08em] uppercase">
                        Location City
                      </FormLabel>
                      <FormControl>
                        <Input
                          placeholder="e.g. Nairobi"
                          className={inputClassName}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="address"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-[var(--on-surface-variant)] font-label text-xs font-bold tracking-[0.08em] uppercase">
                        Specific Area / Address
                      </FormLabel>
                      <FormControl>
                        <Input
                          placeholder="e.g. Westgate Mall, Westlands"
                          className={inputClassName}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </section>

            {/* Section: Capacity */}
            <section
              className="rounded-xl p-6 md:p-8 space-y-6"
              style={{ backgroundColor: "var(--surface-container-lowest)" }}
            >
              <div className="flex items-center gap-3 mb-2">
                <Users
                  className="w-5 h-5"
                  style={{ color: "var(--on-surface)" }}
                />
                <h3
                  className="font-headline text-lg font-bold"
                  style={{ color: "var(--on-surface)" }}
                >
                  Capacity & Infrastructure
                </h3>
              </div>

              <FormField
                control={form.control}
                name="total_seats"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-[var(--on-surface-variant)] font-label text-xs font-bold tracking-[0.08em] uppercase">
                      Total Seating Capacity
                    </FormLabel>
                    <FormControl>
                      <div className="flex items-center gap-3">
                        <Input
                          type="number"
                          min={0}
                          placeholder="0"
                          className={`${inputClassName} max-w-[10rem]`}
                          {...field}
                        />
                        <span
                          className="text-sm font-body"
                          style={{ color: "var(--on-surface-variant)" }}
                        >
                          seats
                        </span>
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </section>
          </div>

          {/* ── Footer ── */}
          <footer
            className="mt-10 py-5 flex flex-col md:flex-row items-center justify-between gap-4 rounded-xl px-6"
            style={{
              backgroundColor: "var(--surface-container-lowest)",
            }}
          >
            <div
              className="flex items-center gap-3 text-sm"
              style={{ color: "var(--on-surface-variant)" }}
            >
              <Info className="w-4 h-4 shrink-0" />
              <span>
                By registering this venue, it will immediately appear in the
                venues dashboard.
              </span>
            </div>

            <div className="flex items-center gap-3 shrink-0">
              <button
                type="button"
                onClick={onCancel}
                className="px-6 py-3 rounded-xl font-body font-bold text-sm cursor-pointer transition-colors hover:opacity-80"
                style={{
                  color: "var(--primary)",
                  backgroundColor: "transparent",
                }}
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={mutation.isPending}
                className="flex items-center gap-2 px-6 py-3 rounded-xl font-body font-bold text-sm transition-all active:scale-95 cursor-pointer disabled:opacity-60 disabled:cursor-not-allowed"
                style={{
                  background:
                    "linear-gradient(135deg, var(--primary), var(--primary-container))",
                  color: "var(--on-primary)",
                }}
              >
                {mutation.isPending ? (
                  "Registering..."
                ) : (
                  <>
                    Register Venue <ArrowRight className="w-4 h-4" />
                  </>
                )}
              </button>
            </div>
          </footer>
        </form>
      </Form>
    </div>
  );
}
