import type PluginManager from "@jbrowse/core/PluginManager"

import { GeneInfoPanel } from "./components"


type Feature = {
  get: (key: string) => string
}

const geneInfoQuery = (gene: string) => ({
  operationName: "GeneGeneralInformationSummary",
  variables: { gene },
  query: `query GeneGeneralInformationSummary($gene: String!) {
    geneGeneralInformation(gene: $gene) {
      id
      name_description
      gene_product
      synonyms
      description
    }
  }`,
})

export default class GeneInfoPlugin {
  name = "GeneInfo"
  version = "1.0.0"

  install(pluginManager: PluginManager) {
    pluginManager.addToExtensionPoint(
      "Core-extraFeaturePanel",
      (DefaultExtraFeature, { model, feature}) => {
        console.log({model, feature})
        if (model.trackType === "FeatureTrack" && feature.type === "gene") {
          return { name: "Gene Info", Component: GeneInfoPanel }
        }
        return DefaultExtraFeature
      },
    )
  }
}
