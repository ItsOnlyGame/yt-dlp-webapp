import { createRootRoute, Outlet } from "@tanstack/react-router";
import "../shadcn.css";
import "../index.css";

import { ThemeProvider } from "@/components/theme-provider";

export const Route = createRootRoute({
  component: () => (
    <>
      <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
        <div className="flex justify-center items-center w-screen h-screen">
          <div className="max-w-7xl w-full mx-4">
            <Outlet />
          </div>
        </div>
      </ThemeProvider>
    </>
  ),
});
