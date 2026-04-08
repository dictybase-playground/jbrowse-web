import type PluginManager from "@jbrowse/core/PluginManager"

const GRAPHQL_URL = "https://graphql.dictybase.dev/graphql"

type GeneInfo = {
  id: string
  name_description: string[]
  gene_product: string[]
  synonyms: string[]
  description: string
}

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
    const React = pluginManager.jbrequire("react") as typeof import("react")
    const { useState, useEffect } = React
    const TestComponent = () => {
      return <>hello</>
    }
    const GeneInfoPanel = ({ feature }: { feature: Feature }) => {
      const [info, setInfo] = useState<GeneInfo | null>(null)
      const geneName = feature.get("name")

      useEffect(() => {
        const fetchGeneInfo = async () => {
          const response = await fetch(GRAPHQL_URL, {
            method: "POST",
            body: JSON.stringify(geneInfoQuery(geneName)),
            headers: { "Content-Type": "application/json" },
          })
          const { data } = await response.json()
          setInfo(data.geneGeneralInformation)
        }
        fetchGeneInfo()
      }, [geneName])

      if (!info) return <div>Loading...</div>

      return (
        <div>
          <p>{info.id}</p>
          <p>{geneName}</p>
          <p>{info.description}</p>
        </div>
      )
    }

    pluginManager.addToExtensionPoint(
      "Core-extraFeaturePanel",
      () => {
          return { name: "Gene Info", Component: TestComponent }
      },
    )
  }

  configure(_pluginManager: PluginManager) {}
}
