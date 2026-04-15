import globals from '@jbrowse/core/ReExports/list'
import {RolldownOptions} from "vite"
import { join } from "path"

export default {
  input: ["plugins/src/index.ts"],
  external: globals,
  treeshake: { propertyReadSideEffects: false, moduleSideEffects: false },
  plugins: [],
  output: [
    {
      dir: "./dist",
      entryFileNames:
      format: 'esm',
      freeze: false,
      esModule: true,
      sourcemap: true,
      exports: 'named',
      inlineDynamicImports: true,
    },
  ],
  watch: { clearScreen: false},
}
