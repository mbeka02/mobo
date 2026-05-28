import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Building2, Users, ArrowRight, ArrowLeft, Info } from "lucide-react";
import { toast } from "sonner";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../../../components/ui/form";
import { Input } from "../../../components/ui/input";
import { api, isAPIError } from "../../../services/api";
import type { CreateVenueAPIResponse } from "../../../components/venues/types";

export const Route = createFileRoute("/portal/venues/new")({
  component: VenueAddPage,
});

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

function VenueAddPage() {
  const navigate = useNavigate();
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
      navigate({ to: "/portal/venues" });
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

  const inputClassName =
    "h-auto bg-[var(--surface-container-high)] text-[var(--on-surface)] rounded-xl px-4 py-4 text-sm font-body border-none focus-visible:ring-2 focus-visible:ring-[var(--primary)] placeholder:text-[var(--on-surface-variant)]/50";

  return (
    <div className="rise-in flex flex-col min-h-[calc(100vh-8rem)]">
      {/* ── Back Button ── */}
      <Link
        to="/portal/venues"
        className="self-start flex items-center gap-2 mb-6 text-sm font-body font-medium cursor-pointer transition-colors hover:opacity-70 no-underline"
        style={{ color: "var(--on-surface-variant)" }}
      >
        <ArrowLeft className="w-4 h-4" />
        Back to Venues
      </Link>

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
              <Link
                to="/portal/venues"
                className="px-6 py-3 rounded-xl font-body font-bold text-sm cursor-pointer transition-colors hover:opacity-80 no-underline"
                style={{
                  color: "var(--primary)",
                  backgroundColor: "transparent",
                }}
              >
                Cancel
              </Link>
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
