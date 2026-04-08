import { defineConfig } from "vite"
import react from "@vitejs/plugin-react"

export default defineConfig({
  plugins: [react({ jsxRuntime: "classic" })],
  build: {
    lib: {
      entry: {
        "hello-world": "hello-world/index.js",
        "gene-info": "gene-info/index.tsx",
      },
      formats: ["es"],
    },
    outDir: "../jbrowse2/plugins",
    emptyOutDir: false,
  },
})
