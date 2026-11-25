import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// The dev server proxies API calls to the Go backend so the browser stays same-origin.
// Change this target if the backend does not run on :8080.
const target = "http://127.0.0.1:8080";

export default defineConfig({
  plugins: [react()],
  server: {
    host: "127.0.0.1",
    port: 5173,
    proxy: Object.fromEntries(
      ["/health", "/auth", "/me", "/uploads"].map((path) => [
        path,
        { target, changeOrigin: false },
      ]),
    ),
  },
});
