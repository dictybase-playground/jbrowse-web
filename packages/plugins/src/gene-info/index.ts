import type PluginManager from "@jbrowse/core/PluginManager"
import { GeneInfoPanel } from "./components"
import { unknownHasTrackType, unknownHasType } from "./utils"
import { pipe } from "fp-ts/function"
import { MonoidAll as BMonoidAll, match as Bmatch } from "fp-ts/boolean"

export default class GeneInfoPlugin {
  name = "GeneInfo"
  version = "1.0.0"

  install(pluginManager: PluginManager) {
    pluginManager.addToExtensionPoint(
      "Core-extraFeaturePanel",
      (DefaultExtraFeature, { model, feature }) => {
        return pipe(
          BMonoidAll.concat(
            unknownHasTrackType(model),
            unknownHasType(feature),
          ),
          Bmatch(
            () => DefaultExtraFeature,
            () => ({ name: "Gene Info", Component: GeneInfoPanel }),
          ),
        )
      },
    )
  }

  configure() {}
}
