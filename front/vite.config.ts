import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig((env) => ({
  plugins: [react()],
  build: {
    target: ["esnext"],
    outDir: "../out/front",
    emptyOutDir: true,
    sourcemap: env.mode === "development",
  },
}));
