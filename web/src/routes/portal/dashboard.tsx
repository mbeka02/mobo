import { createFileRoute, Link } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip as RechartsTooltip,
  ResponsiveContainer,
  Cell,
} from "recharts";
import {
  TrendingUp,
  Film,
  MapPin,
  CalendarDays,
  ArrowRight,
} from "lucide-react";
import { api } from "../../services/api";

export const Route = createFileRoute("/portal/dashboard")({
  component: DashboardPage,
});

interface DashboardStats {
  total_revenue: string;
  tickets_sold: number;
  active_venues: number;
  active_movies: number;
}

interface MonthlyRevenueEntry {
  month: string;
  revenue: string;
}

interface DashboardResponse {
  status: number;
  message: string;
  data: {
    stats: DashboardStats;
    monthly_revenue: MonthlyRevenueEntry[];
  };
}

function formatCurrency(value: string): string {
  const num = parseFloat(value);
  if (isNaN(num)) return "KES 0";
  return `KES ${num.toLocaleString("en-KE", { minimumFractionDigits: 0 })}`;
}

function formatMonth(dateStr: string): string {
  if (!dateStr) return "";
  const date = new Date(dateStr);
  return date.toLocaleDateString("en-KE", { month: "short" });
}

function DashboardPage() {
  const { data, isLoading: loading } = useQuery({
    queryKey: ["dashboardStats"],
    queryFn: async () => {
      const res = await api.get<DashboardResponse>("/admin/dashboard/stats");
      return res.data;
    },
  });

  const stats = data?.stats;
  const chartData = (data?.monthly_revenue || []).map((entry) => ({
    name: formatMonth(entry.month),
    revenue: parseFloat(entry.revenue) || 0,
  }));

  const quickActions = [
    {
      title: "Add Movie",
      desc: "Register a new film",
      icon: Film,
      href: "/portal/movies",
      accentVar: "--primary",
    },
    {
      title: "New Venue",
      desc: "Setup a location",
      icon: MapPin,
      href: "/portal/venues",
      accentVar: "--tertiary",
    },
    {
      title: "Schedule Showtime",
      desc: "Book a screening",
      icon: CalendarDays,
      href: "/portal/scheduling",
      accentVar: "--secondary",
    },
  ];

  if (loading) {
    return (
      <div className="space-y-8 animate-pulse">
        <div className="space-y-3">
          <div
            className="h-12 w-64 rounded-lg"
            style={{ backgroundColor: "var(--surface-container-high)" }}
          />
          <div
            className="h-5 w-96 rounded-md"
            style={{ backgroundColor: "var(--surface-container-high)" }}
          />
        </div>
        <div className="grid grid-cols-12 gap-8">
          <div
            className="col-span-8 h-64 rounded-xl"
            style={{ backgroundColor: "var(--surface-container-low)" }}
          />
          <div className="col-span-4 space-y-8">
            <div
              className="h-[7.5rem] rounded-xl"
              style={{ backgroundColor: "var(--surface-container-low)" }}
            />
            <div
              className="h-[7.5rem] rounded-xl"
              style={{ backgroundColor: "var(--surface-container-low)" }}
            />
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-10 rise-in">
      {/* Page Header */}
      <header className="space-y-2 max-w-3xl">
        <h2
          className="font-headline text-4xl md:text-5xl font-semibold tracking-tight"
          style={{ color: "var(--on-surface)" }}
        >
          Overview
        </h2>
        <p
          className="text-lg"
          style={{ color: "var(--on-surface-variant)" }}
        >
          Here is a summary of your venue and ticketing performance for the
          current season.
        </p>
      </header>

      {/* Metrics Bento Grid */}
      <section className="grid grid-cols-1 md:grid-cols-12 gap-8">
        {/* Revenue Hero Card */}
        <div
          className="md:col-span-8 rounded-xl p-8 flex flex-col justify-between relative overflow-hidden group"
          style={{ backgroundColor: "var(--surface-container-lowest)" }}
        >
          {/* Glow decoration */}
          <div
            className="absolute -right-20 -top-20 w-64 h-64 rounded-full blur-3xl transition-colors duration-700"
            style={{ backgroundColor: "color-mix(in srgb, var(--primary) 5%, transparent)" }}
          />

          <div className="relative z-10">
            <span
              className="font-label text-xs tracking-widest uppercase mb-2 block"
              style={{ color: "var(--on-surface-variant)" }}
            >
              Total Revenue
            </span>
            <div className="flex items-baseline gap-4 mt-2">
              <h3
                className="font-headline text-3xl md:text-4xl font-bold tracking-tight"
                style={{ color: "var(--on-surface)" }}
              >
                {stats ? formatCurrency(stats.total_revenue) : "KES 0"}
              </h3>
              <span
                className="text-sm flex items-center font-medium px-2 py-1 rounded-md"
                style={{
                  color: "var(--primary)",
                  backgroundColor: "color-mix(in srgb, var(--primary-fixed-dim) 20%, transparent)",
                }}
              >
                <TrendingUp className="w-3 h-3 mr-1" />
                14%
              </span>
            </div>
          </div>

          {/* Revenue Chart */}
          <div className="relative z-10 mt-8">
            {chartData.length > 0 ? (
              <ResponsiveContainer width="100%" height={120}>
                <BarChart data={chartData}>
                  <XAxis
                    dataKey="name"
                    axisLine={false}
                    tickLine={false}
                    tick={{ fontSize: 12, fill: "var(--on-surface-variant)" }}
                  />
                  <YAxis hide />
                  <RechartsTooltip
                    cursor={false}
                    contentStyle={{
                      backgroundColor: "var(--surface-container-high)",
                      border: "none",
                      borderRadius: "0.75rem",
                      color: "var(--on-surface)",
                      fontFamily: "'Manrope', sans-serif",
                      fontSize: "0.875rem",
                    }}
                    formatter={(value: any) => [
                      formatCurrency(String(value)),
                      "Revenue",
                    ]}
                  />
                  <Bar dataKey="revenue" radius={[4, 4, 0, 0]}>
                    {chartData.map((_entry, index) => (
                      <Cell
                        key={`cell-${index}`}
                        fill={
                          index === chartData.length - 1
                            ? "var(--primary-container)"
                            : "var(--surface-container-high)"
                        }
                        opacity={index === chartData.length - 1 ? 1 : 0.7}
                      />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div
                className="h-24 flex items-end gap-2 opacity-60"
              >
                {[25, 50, 33, 75, 66, 100].map((h, i) => (
                  <div
                    key={i}
                    className="flex-1 rounded-t-sm"
                    style={{
                      height: `${h}%`,
                      backgroundColor:
                        i === 5
                          ? "var(--primary-container)"
                          : "var(--surface-container-high)",
                      opacity: i === 5 ? 1 : 0.7,
                    }}
                  />
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Side stat cards */}
        <div className="md:col-span-4 flex flex-col gap-8">
          {/* Tickets Sold */}
          <div
            className="flex-1 rounded-xl p-6 flex flex-col justify-center transition-transform hover:-translate-y-1 duration-300"
            style={{ backgroundColor: "var(--surface-container-low)" }}
          >
            <span
              className="font-label text-xs tracking-widest uppercase mb-2 block"
              style={{ color: "var(--on-surface-variant)" }}
            >
              Tickets Sold
            </span>
            <h3
              className="font-headline text-3xl font-bold tracking-tight"
              style={{ color: "var(--on-surface)" }}
            >
              {stats ? stats.tickets_sold.toLocaleString() : "0"}
            </h3>
            <p className="text-sm mt-2" style={{ color: "var(--on-surface-variant)" }}>
              Across {stats?.active_venues || 0} active venues
            </p>
          </div>

          {/* Active Movies */}
          <div
            className="flex-1 rounded-xl p-6 flex flex-col justify-center transition-transform hover:-translate-y-1 duration-300"
            style={{ backgroundColor: "var(--surface-container-low)" }}
          >
            <div className="flex justify-between items-start mb-2">
              <span
                className="font-label text-xs tracking-widest uppercase block"
                style={{ color: "var(--on-surface-variant)" }}
              >
                Active Movies
              </span>
              <Film
                className="w-5 h-5 opacity-70"
                style={{ color: "var(--primary)" }}
              />
            </div>
            <h3
              className="font-headline text-3xl font-bold tracking-tight"
              style={{ color: "var(--on-surface)" }}
            >
              {stats?.active_movies || 0}
            </h3>
            <p className="text-sm mt-2" style={{ color: "var(--on-surface-variant)" }}>
              Currently showing
            </p>
          </div>
        </div>
      </section>

      {/* Quick Actions */}
      <section>
        <header className="mb-6">
          <h3
            className="font-headline text-xl font-semibold"
            style={{ color: "var(--on-surface)" }}
          >
            Quick Actions
          </h3>
        </header>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {quickActions.map((action) => (
            <Link
              key={action.title}
              to={action.href}
              className="group text-left p-6 rounded-xl flex flex-col h-40 justify-between relative overflow-hidden transition-all duration-300"
              style={{
                backgroundColor: "var(--surface-container-lowest)",
                border: "1px solid color-mix(in srgb, var(--outline-variant) 20%, transparent)",
              }}
            >
              <div
                className="w-12 h-12 rounded-lg flex items-center justify-center group-hover:scale-110 transition-transform duration-500"
                style={{
                  backgroundColor: "var(--surface-container)",
                  color: `var(${action.accentVar})`,
                }}
              >
                <action.icon className="w-5 h-5" />
              </div>
              <div>
                <h4
                  className="font-headline font-semibold text-lg transition-colors"
                  style={{ color: "var(--on-surface)" }}
                >
                  {action.title}
                </h4>
                <p
                  className="text-sm mt-1"
                  style={{ color: "var(--on-surface-variant)" }}
                >
                  {action.desc}
                </p>
              </div>
              <ArrowRight
                className="absolute right-6 bottom-6 w-5 h-5 opacity-0 translate-x-4 group-hover:opacity-100 group-hover:translate-x-0 transition-all duration-300"
                style={{ color: "var(--outline-variant)" }}
              />
            </Link>
          ))}
        </div>
      </section>
    </div>
  );
}
