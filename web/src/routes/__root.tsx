import { Outlet, createRootRoute } from "@tanstack/react-router";
import { Toaster } from "sonner";

import "../styles.css";

export const Route = createRootRoute({
  component: RootComponent,
});

function RootComponent() {
  return (
    <>
      <Outlet />
      <Toaster
        position="top-right"
        toastOptions={{
          style: {
            background: "var(--surface-container-high)",
            color: "var(--on-surface)",
            border: "1px solid var(--outline-variant)",
            fontFamily: "'Manrope', sans-serif",
            borderRadius: "0.75rem",
          },
        }}
      />
    </>
  );
}
