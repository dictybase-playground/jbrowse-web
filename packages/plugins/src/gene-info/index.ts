import type PluginManager from "@jbrowse/core/PluginManager"
import { GeneInfoPanel } from "./components"

export default class GeneInfoPlugin {
  name = "GeneInfo"
  version = "1.0.0"

  install(pluginManager: PluginManager) {
    pluginManager.addToExtensionPoint(
      "Core-extraFeaturePanel",
      (DefaultExtraFeature, { model, feature }) => {
        if (model.trackType === "FeatureTrack" && feature.type === "gene") {
          return { name: "Gene Info", Component: GeneInfoPanel }
        }
        return DefaultExtraFeature
      },
    )
  }

  configure(_pluginManager: PluginManager) {}
}
