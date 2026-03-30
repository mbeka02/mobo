import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Eye, EyeOff, ArrowRight } from "lucide-react";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../components/ui/form";
import { Input } from "../components/ui/input";
import { login, redirectToGoogleOAuth, isAuthError } from "../lib/auth";

export const Route = createFileRoute("/login")({ component: LoginPage });

/* ─── Schema ─── */
const loginSchema = z.object({
  email: z.string().email("Please enter a valid email address"),
  password: z.string().min(1, "Password is required"),
});

type LoginFormValues = z.infer<typeof loginSchema>;

/* ─── Page ─── */
function LoginPage() {
  const navigate = useNavigate();
  const [showPassword, setShowPassword] = useState(false);
  const [serverError, setServerError] = useState("");

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  const onSubmit = async (values: LoginFormValues) => {
    setServerError("");
    try {
      await login(values.email, values.password);
      navigate({ to: "/home" });
    } catch (err) {
      if (isAuthError(err)) {
        setServerError(err.message);
      } else {
        setServerError("Something went wrong. Please try again.");
      }
    }
  };

  return (
    <div className="min-h-screen flex">
      {/* ─── Left Panel: Cinematic ─── */}
      <div className="hidden lg:flex lg:w-[55%] relative bg-[#131313] flex-col justify-between p-10 overflow-hidden">
        <div className="absolute inset-0 z-0">
          <img
            src="/images/auth-projector.png"
            alt="Vintage film projector"
            className="w-full h-full object-cover opacity-60"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-[#131313] via-[#131313]/40 to-transparent" />
        </div>

        <Link to="/" className="relative z-10 no-underline">
          <span className="text-2xl font-bold tracking-tighter text-[var(--primary-container)] font-headline">
            Mobo
          </span>
        </Link>

        <div className="relative z-10 space-y-6 mb-16">
          <h1 className="text-5xl xl:text-6xl font-headline font-bold text-[#fff9ef] leading-[0.95] tracking-tight">
            Experience
            <br />
            <span className="text-[var(--primary-container)] italic">
              Cinema
            </span>
            <br />
            Redefined.
          </h1>
          <p className="text-[#fff9ef]/70 text-lg max-w-sm font-body leading-relaxed">
            The Digital Hearth for movie lovers. Curated experiences, seamless
            ticketing, and the warmth of a shared story.
          </p>
        </div>

        <div className="relative z-10 flex items-center gap-3">
          <div className="flex -space-x-2">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className="w-9 h-9 rounded-full bg-[var(--primary-container)]/30 border-2 border-[#131313] flex items-center justify-center"
              >
                <span className="text-[#fff9ef]/60 text-xs font-bold">
                  {String.fromCharCode(64 + i)}
                </span>
              </div>
            ))}
          </div>
          <span className="text-[#fff9ef]/60 font-label text-xs tracking-[0.1em] uppercase">
            Joined by 12K+ Cinephiles
          </span>
        </div>
      </div>

      {/* ─── Right Panel: Form ─── */}
      <div className="flex-1 flex items-center justify-center p-6 sm:p-10 bg-[var(--surface)]">
        <div className="w-full max-w-md space-y-8">
          <Link to="/" className="lg:hidden block mb-8 no-underline">
            <span className="text-2xl font-bold tracking-tighter text-[var(--primary-container)] font-headline">
              Mobo
            </span>
          </Link>

          <div>
            <h2 className="text-3xl font-headline font-bold text-[var(--on-surface)]">
              Welcome back
            </h2>
            <p className="text-[var(--on-surface-variant)] mt-2">
              Please enter your details to continue your journey.
            </p>
          </div>

          {/* Google OAuth */}
          <button
            type="button"
            onClick={redirectToGoogleOAuth}
            className="flex items-center justify-center gap-3 w-full py-4 rounded-xl bg-[var(--surface-container-high)] text-[var(--on-surface)] font-bold text-sm hover:bg-[var(--surface-container-highest)] transition-colors cursor-pointer"
          >
            <svg viewBox="0 0 24 24" width="20" height="20">
              <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 01-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4" />
              <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
              <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" />
              <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
            </svg>
            CONTINUE WITH GOOGLE
          </button>

          {/* Divider */}
          <div className="flex items-center gap-4">
            <div className="flex-1 h-px bg-[var(--outline-variant)]" />
            <span className="text-[var(--on-surface-variant)] font-label text-xs tracking-[0.1em] uppercase">
              or email
            </span>
            <div className="flex-1 h-px bg-[var(--outline-variant)]" />
          </div>

          {/* Server Error */}
          {serverError && (
            <div className="bg-[var(--error-container)] text-[var(--on-error-container)] px-4 py-3 rounded-xl text-sm font-medium">
              {serverError}
            </div>
          )}

          {/* Form */}
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-[var(--on-surface-variant)] font-label text-xs font-bold tracking-[0.08em] uppercase">
                      Email Address
                    </FormLabel>
                    <FormControl>
                      <Input
                        type="email"
                        placeholder="name@example.com"
                        className="h-auto bg-[var(--surface-container-high)] text-[var(--on-surface)] rounded-xl px-4 py-4 text-sm font-body border-none focus-visible:ring-2 focus-visible:ring-[var(--primary)] placeholder:text-[var(--on-surface-variant)]/50"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <div className="flex justify-between items-center">
                      <FormLabel className="text-[var(--on-surface-variant)] font-label text-xs font-bold tracking-[0.08em] uppercase">
                        Password
                      </FormLabel>
                      <a
                        href="#"
                        className="text-[var(--primary)] font-label text-xs font-bold no-underline hover:underline"
                      >
                        Forgot?
                      </a>
                    </div>
                    <FormControl>
                      <div className="relative">
                        <Input
                          type={showPassword ? "text" : "password"}
                          className="h-auto bg-[var(--surface-container-high)] text-[var(--on-surface)] rounded-xl px-4 py-4 text-sm font-body border-none focus-visible:ring-2 focus-visible:ring-[var(--primary)] pr-12"
                          {...field}
                        />
                        <button
                          type="button"
                          onClick={() => setShowPassword(!showPassword)}
                          className="absolute right-4 top-1/2 -translate-y-1/2 text-[var(--on-surface-variant)] hover:text-[var(--on-surface)]"
                        >
                          {showPassword ? (
                            <EyeOff size={20} />
                          ) : (
                            <Eye size={20} />
                          )}
                        </button>
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <button
                type="submit"
                disabled={form.formState.isSubmitting}
                className="w-full py-4 bg-gradient-to-r from-[var(--primary)] to-[var(--primary-container)] text-white rounded-xl font-bold text-base flex items-center justify-center gap-2 hover:shadow-[0_12px_32px_rgba(171,54,0,0.3)] transition-all disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer"
              >
                {form.formState.isSubmitting ? (
                  "Signing in..."
                ) : (
                  <>
                    Sign In <ArrowRight size={18} />
                  </>
                )}
              </button>
            </form>
          </Form>

          <p className="text-center text-[var(--on-surface-variant)] text-sm">
            New to Mobo?{" "}
            <Link
              to="/signup"
              className="text-[var(--primary)] font-bold no-underline hover:underline"
            >
              Create an account
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
