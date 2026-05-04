import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "node:path";
import { fileURLToPath } from "node:url";
var __dirname = path.dirname(fileURLToPath(import.meta.url));
export default defineConfig({
    plugins: [react()],
    server: {
        proxy: {
            "/api": {
                target: "http://127.0.0.1:8090",
                changeOrigin: true,
            },
        },
    },
    build: {
        chunkSizeWarningLimit: 1000,
        rollupOptions: {
            output: {
                manualChunks: {
                    react: ["react", "react-dom", "react-router-dom"],
                    query: ["@tanstack/react-query"],
                    flow: ["reactflow"],
                    ui: ["@radix-ui/react-dialog", "@radix-ui/react-separator", "@radix-ui/react-slot", "lucide-react", "sonner"],
                    utils: ["js-yaml", "class-variance-authority", "clsx", "tailwind-merge"],
                },
            },
        },
    },
    resolve: {
        dedupe: ["react", "react-dom"],
        alias: {
            "@": path.resolve(__dirname, "./src"),
        },
    },
});
