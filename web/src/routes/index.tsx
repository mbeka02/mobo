import { createFileRoute, Link } from "@tanstack/react-router";
import LandingNavbar from "../components/LandingNavbar";
import {
  Star,
  ArrowRight,
  ChevronLeft,
  ChevronRight,
  QrCode,
  MapPin,
  Users,
  Send,
  Ticket,
} from "lucide-react";

export const Route = createFileRoute("/")({ component: LandingPage });

/* ─── Cinemas Data ─── */
const cinemas = [
  {
    name: "Anga Sky Cinema",
    location: "Panari Center, Mombasa Rd",
    image: "/images/cinema-anga.png",
    price: "800/-",
    badge: "PREMIUM",
    active: true,
  },
  {
    name: "Century Cinepax",
    location: "Sarit Centre, Westlands",
    image: "/images/cinema-century.png",
    active: false,
  },
  {
    name: "Fox Cineplex",
    location: "Sarit Center, Nairobi",
    image: "/images/cinema-fox.png",
    active: false,
  },
];

/* ─── Landing Page ─── */
function LandingPage() {
  return (
    <div className="min-h-screen bg-[var(--surface)]">
      <LandingNavbar />

      {/* ═══════ HERO SECTION ═══════ */}
      <section className="relative min-h-[90vh] flex items-center px-6 md:px-8 overflow-hidden pt-20">
        {/* Background */}
        <div className="absolute inset-0 z-0">
          <img
            src="/images/hero-cinema.png"
            alt="Luxury cinema interior"
            className="w-full h-full object-cover opacity-90"
          />
          <div className="absolute inset-0 bg-gradient-to-r from-[#131313] via-[#131313]/50 to-transparent" />
        </div>

        <div className="relative z-10 max-w-7xl mx-auto w-full grid md:grid-cols-2 gap-12 items-center">
          {/* Left: Copy */}
          <div className="space-y-8 rise-in">
            <div className="inline-block px-4 py-1.5 rounded-full bg-[var(--primary-container)]/20 border border-[var(--primary-container)]/30">
              <span className="text-[var(--primary-container)] font-label text-xs font-bold tracking-[0.16em] uppercase">
                The Digital Hearth
              </span>
            </div>
            <h1 className="text-5xl sm:text-6xl md:text-8xl font-headline font-bold text-[#fff9ef] leading-[0.9] tracking-tight">
              Experience Cinema{" "}
              <span className="text-[var(--primary-container)]">Redefined.</span>
            </h1>
            <p className="text-lg md:text-xl text-[#fff9ef]/80 max-w-md font-body leading-relaxed">
              Mobo is the Digital Hearth for movie lovers. Discover curated
              cinematic journeys and book seats in seconds.
            </p>
            <div className="flex flex-wrap gap-4 pt-2">
              <Link
                to="/login"
                className="inline-flex items-center gap-2 px-8 py-4 bg-gradient-to-br from-[var(--primary)] to-[var(--primary-container)] text-white rounded-xl font-bold shadow-xl hover:scale-95 transition-transform no-underline"
              >
                Explore Now <ArrowRight size={18} />
              </Link>
              <button className="px-8 py-4 bg-white/10 backdrop-blur-md border border-white/20 text-[#fff9ef] rounded-xl font-bold hover:bg-white/20 transition-all">
                Join the Club
              </button>
            </div>
          </div>

          {/* Right: Featured Movie Card */}
          <div
            className="hidden md:block relative rise-in"
            style={{ animationDelay: "200ms" }}
          >
            <div className="absolute -top-12 -right-12 w-64 h-64 bg-[var(--primary)]/20 blur-[100px] rounded-full" />
            <div className="bg-white/5 backdrop-blur-2xl p-6 rounded-3xl border border-white/10 shadow-2xl rotate-2 hover:rotate-0 transition-transform duration-500">
              <img
                src="/images/cinema-marquee.png"
                alt="Cinema marquee"
                className="rounded-2xl w-full aspect-[4/5] object-cover"
              />
              <div className="mt-6 flex justify-between items-end">
                <div>
                  <p className="text-xs text-[#fff9ef]/60 font-label uppercase tracking-[0.16em]">
                    Now Showing
                  </p>
                  <h3 className="text-2xl font-headline font-bold text-[#fff9ef]">
                    Dune: Part Two
                  </h3>
                </div>
                <div className="bg-[var(--primary-container)] p-3 rounded-xl text-white">
                  <Ticket size={24} />
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* ═══════ BENTO FEATURES ═══════ */}
      <section className="py-24 md:py-32 px-6 md:px-8 bg-[var(--surface)]">
        <div className="max-w-7xl mx-auto">
          <div className="mb-16 md:mb-20 space-y-4">
            <h2 className="text-3xl sm:text-4xl md:text-6xl font-headline font-bold tracking-tight text-[var(--on-surface)]">
              Designed for the{" "}
              <span className="italic font-light">Cinephile.</span>
            </h2>
            <p className="text-[var(--on-surface-variant)] max-w-xl text-base md:text-lg">
              We've built Mobo to feel as tactile as a physical ticket stub while
              being fast as light.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-12 gap-6 auto-rows-auto">
            {/* Feature 1: Seamless Ticketing */}
            <div className="md:col-span-8 bg-[var(--surface-container-low)] rounded-[2rem] p-8 md:p-12 flex flex-col justify-between group overflow-hidden relative min-h-[300px]">
              <div className="max-w-md relative z-10">
                <QrCode
                  size={48}
                  className="text-[var(--primary)] mb-6"
                  strokeWidth={1.5}
                />
                <h3 className="text-2xl md:text-4xl font-headline font-bold mb-4">
                  Seamless Ticketing
                </h3>
                <p className="text-[var(--on-surface-variant)] text-base md:text-lg leading-relaxed">
                  Focus on the film, not the queue. Direct M-Pesa integration
                  means instant confirmation and a secure QR code delivered right
                  to your device.
                </p>
              </div>
            </div>

            {/* Feature 2: Curated Experiences */}
            <div className="md:col-span-4 bg-[var(--primary)] text-white rounded-[2rem] p-8 md:p-12 flex flex-col items-start justify-center space-y-6 relative overflow-hidden">
              <div className="absolute inset-0 opacity-20 bg-[radial-gradient(circle_at_center,_white,_transparent_70%)]" />
              <h3 className="text-2xl md:text-3xl font-headline font-bold relative z-10">
                Curated Experiences
              </h3>
              <p className="text-white/80 font-body relative z-10">
                From the boutique charm of Nairobi to the coastal magic of
                Mombasa. We partner with only the finest premium cinemas.
              </p>
              <button className="bg-white text-[var(--primary)] px-8 py-3 rounded-full font-bold relative z-10 hover:scale-105 transition-transform">
                Explore Locations
              </button>
            </div>

            {/* Feature 3: Community */}
            <div className="md:col-span-5 bg-[var(--surface-container-high)] rounded-[2rem] p-8 md:p-12 flex flex-col justify-end relative overflow-hidden min-h-[350px]">
              <div className="absolute top-0 left-0 w-full h-1/2 bg-gradient-to-b from-[var(--surface-container-high)] to-transparent z-10" />
              <img
                src="/images/cinema-community.png"
                alt="Friends at cinema"
                className="absolute inset-0 w-full h-full object-cover grayscale brightness-75 hover:grayscale-0 transition-all duration-700"
              />
              <div className="relative z-20 text-white">
                <h3 className="text-2xl md:text-3xl font-headline font-bold mb-2">
                  The Cinephile Community
                </h3>
                <p className="text-white/80">
                  Join 12,000+ members and unlock loyalty rewards on every
                  booking.
                </p>
              </div>
            </div>

            {/* Feature 4: Loyalty */}
            <div className="md:col-span-7 bg-[var(--surface-container-lowest)] rounded-[2rem] p-8 md:p-12 flex flex-col sm:flex-row items-center gap-8">
              <div className="shrink-0">
                <div className="w-28 h-28 md:w-32 md:h-32 rounded-full bg-[var(--primary-container)]/10 flex items-center justify-center text-[var(--primary-container)]">
                  <Star size={56} strokeWidth={1.5} />
                </div>
              </div>
              <div>
                <h3 className="text-xl md:text-2xl font-headline font-bold mb-2 text-[var(--on-surface)]">
                  Loyalty That Counts
                </h3>
                <p className="text-[var(--on-surface-variant)]">
                  Collect stamps for every movie. Your 5th ticket is on the house.
                  Cinema as it was meant to be experienced—exclusive and rewarding.
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* ═══════ CINEMA CAROUSEL ═══════ */}
      <section className="py-24 md:py-32 bg-[var(--surface-container-low)] overflow-hidden">
        <div className="max-w-7xl mx-auto px-6 md:px-8 mb-12 flex flex-col sm:flex-row justify-between items-start sm:items-end gap-4">
          <div className="space-y-2">
            <span className="text-[var(--primary)] font-label text-sm font-bold tracking-[0.2em] uppercase">
              Premiering Soon
            </span>
            <h2 className="text-3xl sm:text-4xl md:text-5xl font-headline font-bold">
              The Nairobi Selection
            </h2>
          </div>
          <div className="flex gap-3">
            <button className="w-12 h-12 rounded-full border border-[var(--outline)] flex items-center justify-center hover:bg-[var(--surface-container-high)] transition-colors">
              <ChevronLeft size={20} />
            </button>
            <button className="w-12 h-12 rounded-full bg-[var(--primary)] text-white flex items-center justify-center">
              <ChevronRight size={20} />
            </button>
          </div>
        </div>

        <div className="flex gap-6 md:gap-8 px-6 md:px-8 overflow-x-auto no-scrollbar">
          {cinemas.map((cinema) => (
            <div
              key={cinema.name}
              className={`flex-shrink-0 w-[300px] md:w-[420px] rounded-3xl p-5 md:p-6 transition-all duration-300 ${
                cinema.active
                  ? "bg-[var(--surface-container-lowest)] shadow-xl"
                  : "bg-[var(--surface-dim)]/40 scale-95 opacity-60"
              }`}
            >
              <img
                src={cinema.image}
                alt={cinema.name}
                className="w-full aspect-[16/9] object-cover rounded-2xl mb-5 md:mb-6"
              />
              <div className="space-y-4">
                <div className="flex justify-between items-start">
                  <div>
                    <h4 className="text-xl md:text-2xl font-headline font-bold">
                      {cinema.name}
                    </h4>
                    <p className="text-[var(--on-surface-variant)] font-medium flex items-center gap-1 mt-1">
                      <MapPin size={14} /> {cinema.location}
                    </p>
                  </div>
                  {cinema.badge && (
                    <span className="bg-[var(--primary-container)]/10 text-[var(--primary-container)] px-3 py-1 rounded-full text-xs font-bold">
                      {cinema.badge}
                    </span>
                  )}
                </div>
                {cinema.active && cinema.price && (
                  <div className="pt-4 flex justify-between items-center" style={{ borderTop: "1px solid var(--outline-variant)" }}>
                    <span className="font-bold text-lg">
                      Tickets from {cinema.price}
                    </span>
                    <button className="text-[var(--primary)] font-bold flex items-center gap-1 hover:gap-2 transition-all">
                      View Times <ArrowRight size={16} />
                    </button>
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* ═══════ COMMUNITY CTA ═══════ */}
      <section className="py-32 md:py-40 bg-[var(--surface)]">
        <div className="max-w-4xl mx-auto px-6 md:px-8 text-center space-y-10 md:space-y-12">
          <div className="w-20 h-20 md:w-24 md:h-24 bg-[var(--primary)] rounded-full mx-auto flex items-center justify-center text-white">
            <Users size={40} />
          </div>
          <h2 className="text-4xl sm:text-5xl md:text-7xl font-headline font-bold tracking-tight">
            Be part of the{" "}
            <span className="text-[var(--primary)]">Inner Circle.</span>
          </h2>
          <p className="text-lg md:text-2xl text-[var(--on-surface-variant)] font-light leading-relaxed max-w-2xl mx-auto">
            Join over 12,000 cinephiles in Nairobi and Mombasa who get early
            access to premieres, exclusive discounts, and invitations to
            Mobo-only screenings.
          </p>
          <div className="pt-6 md:pt-8">
            <button className="px-10 md:px-12 py-5 md:py-6 bg-[var(--primary)] text-white rounded-2xl font-bold text-lg md:text-xl hover:shadow-[0_20px_50px_rgba(171,54,0,0.3)] transition-all">
              Join Mobo Club Today
            </button>
            <p className="mt-6 text-[var(--on-surface-variant)]/60 font-label text-sm tracking-[0.1em]">
              NO FEES. JUST PURE CINEMA.
            </p>
          </div>
        </div>
      </section>

      {/* ═══════ FOOTER ═══════ */}
      <footer className="bg-[#1c1b1b] w-full pt-16 pb-8">
        <div className="max-w-7xl mx-auto px-6 md:px-8 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-10 md:gap-12">
          <div>
            <div className="text-xl font-bold text-[var(--primary-container)] mb-4 font-headline">
              Mobo
            </div>
            <p className="text-[#fff9ef]/50 font-body text-sm leading-relaxed">
              Redefining the cinematic experience in East Africa through design,
              community, and technology.
            </p>
          </div>
          <div>
            <h5 className="text-[#fff9ef] font-headline font-bold mb-5">
              Quick Links
            </h5>
            <ul className="space-y-3 list-none p-0">
              {["About Us", "Partner with Us", "Help Center"].map((link) => (
                <li key={link}>
                  <a
                    href="#"
                    className="text-[#fff9ef]/50 hover:text-[#ffb59c] transition-colors font-label text-sm no-underline"
                  >
                    {link}
                  </a>
                </li>
              ))}
            </ul>
          </div>
          <div>
            <h5 className="text-[#fff9ef] font-headline font-bold mb-5">
              Legal
            </h5>
            <ul className="space-y-3 list-none p-0">
              {["Privacy Policy", "Terms of Service"].map((link) => (
                <li key={link}>
                  <a
                    href="#"
                    className="text-[#fff9ef]/50 hover:text-[#ffb59c] transition-colors font-label text-sm no-underline"
                  >
                    {link}
                  </a>
                </li>
              ))}
            </ul>
          </div>
          <div>
            <h5 className="text-[#fff9ef] font-headline font-bold mb-5">
              Newsletter
            </h5>
            <p className="text-[#fff9ef]/50 font-body text-xs mb-4">
              Get weekly showtimes and offers.
            </p>
            <div className="relative">
              <input
                type="email"
                placeholder="Email address"
                className="w-full bg-[#2a2828] border-none rounded-lg text-white font-label text-sm px-4 py-3 focus:outline-none focus:ring-1 focus:ring-[var(--primary)]"
              />
              <button
                className="absolute right-3 top-1/2 -translate-y-1/2 text-[var(--primary)]"
                aria-label="Subscribe"
              >
                <Send size={18} />
              </button>
            </div>
          </div>
        </div>
        <div className="max-w-7xl mx-auto px-6 md:px-8 mt-12 md:mt-16 pt-8 border-t border-white/5 flex flex-col md:flex-row justify-between items-center gap-4">
          <p className="text-[#fff9ef]/50 font-label text-xs">
            © 2024 Mobo Ticketing. Crafted with Cinematic Soul in Nairobi.
          </p>
        </div>
      </footer>
    </div>
  );
}
