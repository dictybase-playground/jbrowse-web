import { defineConfig } from "vite"
import react from "@vitejs/plugin-react"
import { rolldownOptions } from "./esmRolldownConfig"

export default defineConfig({
  plugins: [react({ jsxRuntime: "classic" })],
  build: {
    lib: {
      entry: {
        "gene-info": "gene-info/index.tsx",
      },
      formats: ["es"],
    },
    outDir: "../jbrowse2/plugins",
    emptyOutDir: false,
    rolldownOptions,
  },
})
