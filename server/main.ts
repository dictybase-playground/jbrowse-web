import { basename } from "path"

export function startServer(root: string, assets: string) {
  const assetsName = basename(assets)
  const assetsRoute = `/${assetsName}/*` as const

  const server = Bun.serve({
    routes: {
      [assetsRoute]: (req: Request) => {
        const { pathname } = new URL(req.url)
        return new Response(Bun.file(`${assets}${pathname.slice(assetsName.length + 1)}`))
      },
    },
    fetch(req) {
      const { pathname } = new URL(req.url)
      const filePath = pathname === "/" ? "index.html" : pathname
      return new Response(Bun.file(`${root}/${filePath}`))
    },
    error(err: NodeJS.ErrnoException) {
      if (err.code === "ENOENT") return new Response("Not Found", { status: 404 })
      return new Response(err.message, { status: 500 })
    },
  })

  console.log(`App:    ${root}`)
  console.log(`Assets: ${assets}`)
  console.log(`Listening on http://${server.hostname}:${server.port}`)
}
