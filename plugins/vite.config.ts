import { defineConfig } from "vite"

export default defineConfig({
  build: {
    lib: {
      entry: "hello-world/index.js",
      name: "JBrowsePluginHelloWorld",
      fileName: (format) => `hello-world.${format}.js`,
      formats: ["umd"],
    },
    outDir: "../jbrowse2/plugins",
    emptyOutDir: false,
    rollupOptions: {
      output: {
        exports: "named",
      },
    },
  },
})
