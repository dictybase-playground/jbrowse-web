import { defineConfig } from "vite"
import react from "@vitejs/plugin-react"
import { rolldownOptions } from "./rolldown.config.ts"

export default defineConfig({
  plugins: [react({ jsxRuntime: "classic" })],
  build: {
    lib: {
      entry: {
        "gene-info": "gene-info/index.tsx",
      },
    },
    emptyOutDir: false,
    rolldownOptions,
  },
})
