import { Hono } from "hono"
import { logger } from "hono/logger"
import { serveStatic } from "hono/bun"
import { basename, dirname } from "path"

export function startServer(root: string, assetsPath: string) {
  const app = new Hono()
  app.use("*", logger())
  const assetsBaseDir = basename(assetsPath)

  // Serve assets under /{assetsBaseDir}/*
  app.use(`/${assetsBaseDir}/*`, serveStatic({ root: dirname(assetsPath) }))

  // Serve JBrowse app for everything else
  app.use("*", serveStatic({ root }))

  const server = Bun.serve({ port: 3000, fetch: app.fetch })
  console.log(`App:    ${root}`)
  console.log(`Assets: ${assetsPath}`)
  console.log(`Listening on http://${server.hostname}:${server.port}`)
}
