import { Hono } from "hono"
import { logger } from "hono/logger"
import { serveStatic } from "hono/bun"
import { basename } from "path"

type ServerOptions = {
  port: number
}

export function startServer(
  root: string,
  assetsPath: string,
  options: ServerOptions,
) {
  const app = new Hono()

  app.use("*", logger())
  // Serve assets under /{assetsBaseDir}/*
  const assetsBaseDir = basename(assetsPath)
  app.use(
    `/${assetsBaseDir}/*`,
    serveStatic({
      root: assetsPath,
      rewriteRequestPath: (path) => path.replace(`/${assetsBaseDir}`, ""),
    }),
  )

  // Serve JBrowse app for everything else
  app.use("*", serveStatic({ root }))

  const server = Bun.serve({ port: options.port, fetch: app.fetch })
  console.log(`App:    ${root}`)
  console.log(`Assets: ${assetsPath}`)
  console.log(`Listening on http://${server.hostname}:${server.port}`)
}
