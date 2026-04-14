import type PluginManager from "@jbrowse/core/PluginManager"
import { TestComponent } from "./components"


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
      () => {
          return { name: "Gene Info", Component: TestComponent }
      },
    )
  }

  configure(_pluginManager: PluginManager) {}
}
