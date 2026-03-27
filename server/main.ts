import { basename } from "path"
import { pipe } from "fp-ts/function"
import { match as Bmatch } from "fp-ts/boolean"
import {
  getOrElse as OgetOrElse,
  fromPredicate as OfromPredicate,
} from "fp-ts/Option"
import { not } from "fp-ts/Predicate"
import {
  uniq as RNEAuniq,
  intercalate as RNEAintercalate,
} from "fp-ts/ReadonlyNonEmptyArray"
import {
  Eq as SEq,
  Monoid as SMonoid,
  split as Ssplit,
  Semigroup as SSemigroup,
} from "fp-ts/string"

/**
 * Starts a static server for serving the Jnrowse application.
 * `root` is the directory that contains the `index.html`.
 * `assetsPath` contains the jbrowse data assets such as FASTA files, GFF3 files, BAM files, et cetera.
 */
export function startServer(root: string, assetsPath: string) {
  // Get the base directory of data assets path.
  const assetsBaseDir = basename(assetsPath)
  // Use assetsBaseDir to construct wildcard route key for jbrowse data assets.
  const assetsRoute = `/${assetsBaseDir}/*`

  const server = Bun.serve({
    routes: {
      [assetsRoute]: ({ url }: Request) => {
        // Get pathname of request, e.g: `/test_data/fasta/dicty.fa`
        const { pathname: requestPath } = new URL(url)
        // requestPath:        `/test_data/fasta/dicty.fa`
        // assetsPath:      `assets/test_data`
        // assetsBaseDir:   `test_data`

        // file:            `assets/test_data/fasta/dicty.fa`
        // assetsPath + requestPath = `assets/test_data/test_data/fasta/dicty.fa`
        // assetsPath + (requestPath - /assetsBaseDir) = `assets/test_data/fasta/dicty.fa`
        const assetFile = pipe(
          SMonoid.concat(assetsPath, requestPath),
          Ssplit("/"),
          RNEAuniq(SEq),
          RNEAintercalate(SSemigroup)("/"),
        )
        return new Response(Bun.file(assetFile))
      },
    },
    fetch({ url }) {
      const { pathname } = new URL(url)

      const filePath = pipe(
        pathname,
        OfromPredicate(not((a) => SEq.equals(a, "/"))),
        OgetOrElse(() => "index.html"),
      )

      return new Response(Bun.file(`${root}/${filePath}`))
    },
    error(err: NodeJS.ErrnoException) {
      return pipe(
        err,
        ({ code }) => code === "ENOENT",
        Bmatch(
          () =>
            new Response("Not Found", {
              status: 404,
            }),

          () =>
            new Response(err.message, {
              status: 500,
            }),
        ),
      )
    },
  })

  console.log(`App:    ${root}`)
  console.log(`Assets: ${assetsPath}`)
  console.log(`Listening on http://${server.hostname}:${server.port}`)
}
