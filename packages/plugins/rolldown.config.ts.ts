import {  ResolvedBuildEnvironmentOptions } from "vite"
import { isExternal } from "./isExternal"
import jbrowseGlobals from "@jbrowse/core/ReExports/list"
import externalGlobals from "rollup-plugin-external-globals"
import { pipe } from "fp-ts/function"
import { map as Amap } from "fp-ts/Array"
import { map as Rmap, fromEntries as RfromEntries } from "fp-ts/Record"

const globalsMap = pipe(
  jbrowseGlobals,
  Amap((id): [string, string] => [id, id]),
  RfromEntries,
  Rmap((global) => `JBrowseExports["${global}"]`),
)

const rolldownOptions: ResolvedBuildEnvironmentOptions["rolldownOptions"] = {
  input: ["plugins/src/gene-info/index.ts"],
  external: isExternal, // not sure if necessary.
  treeshake: { propertyReadSideEffects: false, moduleSideEffects: false },
  plugins: [externalGlobals(globalsMap)],
  output: [
    {
      dir: "./dist",
      entryFileNames: "testplugin.js",
      format: "esm",
      esModule: true,
      sourcemap: false,
      exports: "named",
      codeSplitting: false,
    },
  ],
  watch: { clearScreen: false },
}

export { rolldownOptions }
