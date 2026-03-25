import { join } from "path"
import { pipe } from "fp-ts/function"

const base = process.argv[2] ?? join(import.meta.dir, "..")
const appDir = join(base, "jbrowse2")
const dataDir = join(base, "test_data")

const server = Bun.serve({
  routes: {
    "/test_data/*": (req) => {
      const { pathname } = new URL(req.url)
      return new Response(Bun.file(join(dataDir, pathname.slice("/test_data".length))))
    },
  },
  fetch(req) {
    const { pathname } = new URL(req.url)
    const filePath = pathname === "/" ? "index.html" : pathname
    return new Response(Bun.file(join(appDir, filePath)))
  },
  error(err: NodeJS.ErrnoException) {
    if (err.code === "ENOENT") return new Response("Not Found", { status: 404 })
    return new Response(err.message, { status: 500 })
  },
})

console.log(`App:  ${appDir}`)
console.log(`Data: ${dataDir}`)
console.log(`Listening on http://${server.hostname}:${server.port}`)
