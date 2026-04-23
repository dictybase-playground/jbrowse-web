import { defineConfig } from "vite"
import { rolldownOptions } from "./rolldown.config.ts"

export default defineConfig({
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
