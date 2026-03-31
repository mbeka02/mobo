import { Link } from "@tanstack/react-router";
import { Search, Menu, X } from "lucide-react";
import { useState } from "react";

const navLinks = [
  { label: "Movies", href: "#" },
  { label: "Cinemas", href: "#" },
  { label: "Offers", href: "#" },
  { label: "Support", href: "#" },
];

export default function LandingNavbar() {
  const [mobileOpen, setMobileOpen] = useState(false);

  return (
    <nav className="fixed top-0 w-full z-50 bg-[#131313]/70 backdrop-blur-xl">
      <div className="flex justify-between items-center px-6 md:px-8 py-4 max-w-7xl mx-auto">
        {/* Left: Logo + Nav Links */}
        <div className="flex items-center gap-10">
          <Link to="/" className="no-underline">
            <span className="text-2xl font-bold tracking-tighter text-[var(--primary-container)] font-headline">
              Mobo
            </span>
          </Link>
          <div className="hidden md:flex gap-8 items-center">
            {navLinks.map((link) => (
              <a
                key={link.label}
                href={link.href}
                className="font-label text-sm font-medium text-[#fff9ef]/70 hover:text-[#fff9ef] transition-colors no-underline"
              >
                {link.label}
              </a>
            ))}
          </div>
        </div>

        {/* Right: Search + CTA */}
        <div className="flex items-center gap-4">
          <button
            className="text-[#fff9ef]/70 hover:text-[#ffb59c] transition-all p-2"
            aria-label="Search"
          >
            <Search size={20} />
          </button>
          <Link
            to="/login"
            className="hidden sm:inline-flex bg-[var(--primary)] text-white px-6 py-2.5 rounded-full font-label font-bold text-sm hover:scale-95 transition-transform no-underline"
          >
            Get Tickets
          </Link>

          {/* Mobile Menu Toggle */}
          <button
            className="md:hidden text-[#fff9ef]/70 p-2"
            onClick={() => setMobileOpen(!mobileOpen)}
            aria-label="Toggle menu"
          >
            {mobileOpen ? <X size={24} /> : <Menu size={24} />}
          </button>
        </div>
      </div>

      {/* Mobile Menu */}
      {mobileOpen && (
        <div className="md:hidden bg-[#131313]/95 backdrop-blur-xl px-6 pb-6 fade-in">
          <div className="flex flex-col gap-4">
            {navLinks.map((link) => (
              <a
                key={link.label}
                href={link.href}
                className="font-label text-base font-medium text-[#fff9ef]/80 hover:text-[#fff9ef] transition-colors no-underline py-2"
                onClick={() => setMobileOpen(false)}
              >
                {link.label}
              </a>
            ))}
            <Link
              to="/login"
              className="bg-[var(--primary)] text-white px-6 py-3 rounded-full font-label font-bold text-sm text-center hover:scale-95 transition-transform no-underline mt-2"
              onClick={() => setMobileOpen(false)}
            >
              Get Tickets
            </Link>
          </div>
        </div>
      )}
    </nav>
  );
}
