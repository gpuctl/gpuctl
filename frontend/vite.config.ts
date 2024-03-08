import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    // Keep this in sync with the `internal/webapi/server.go` CORS header.
    port: 5173,
  },
});
