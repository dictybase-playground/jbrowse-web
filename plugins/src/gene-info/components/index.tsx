import type PluginManager from "@jbrowse/core/PluginManager"
import React, { useState, useEffect } from "react"
type Feature = {
  get: (key: string) => string
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

const TestComponent = () => {
  const [info] = useState<GeneInfo | null>(null)
  return <>{info}</>
}

const Component = () => {
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
  }

export { Component, TestComponent }
