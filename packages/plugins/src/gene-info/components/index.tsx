import { useState, useEffect } from "react"
import { match, P } from "ts-pattern"

type Feature = {
  id: string
}

type GeneInfo = {
  id: string
  name_description: string[]
  gene_product: string[]
  synonyms: string[]
  description: string
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

const GRAPHQL_URL = "https://graphql.dictybase.dev/graphql"

const GeneInfoPanel = ({ feature }: { feature: Feature }) => {
  const [info, setInfo] = useState<GeneInfo | null>(null)
  const [loading, setLoading] = useState(true)
  const geneName = feature.id

  useEffect(() => {
    const fetchGeneInfo = async () => {
      const response = await fetch(GRAPHQL_URL, {
        method: "POST",
        body: JSON.stringify(geneInfoQuery(geneName)),
        headers: { "Content-Type": "application/json" },
      })
      const { data } = await response.json()
      setLoading(false)
      setInfo(data.geneGeneralInformation)
    }
    fetchGeneInfo()
  }, [geneName])

  return match({ loading, info })
    .with({ loading: true }, () => <div>Loading...</div>)
    .with({ info: P.select(P.nonNullable) }, (info) => {
      console.log(info)
        return (
            <div>
                <p>Description: {info.description}</p>
                <p>Gene Product: {info.gene_product}</p>
            </div>
        )
    })
    .otherwise(() => <></>)
}

export { GeneInfoPanel }
