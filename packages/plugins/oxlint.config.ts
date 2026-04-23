import { defineConfig } from "oxlint"

export default defineConfig({
  plugins: [
    "typescript",
    "unicorn",
    "react",
    "nextjs",
    "oxc",
    "import",
    "jsdoc",
    "jsx-a11y",
    "promise",
    "vitest",
  ],
  categories: {
    correctness: "error",
    perf: "error",
    suspicious: "error",
  },
  rules: {
    "react-in-jsx-scope": "allow",
    "no-unassigned-import": "allow",
    "no-autofocus": "allow",
    "no-standalone-expect": "allow",
  },
  settings: {
    react: {
      version: "17",
    },
  },
})
