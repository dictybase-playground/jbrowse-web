import { Command } from "commander"
import { existsSync, statSync } from "fs"
import { pipe } from "fp-ts/function"
import {
  map as Emap,
  fold as Efold,
  fromPredicate as EfromPredicate,
  Applicative as EApplicative,
} from "fp-ts/Either"
import { and } from "fp-ts/Predicate"
import { traverse as Atraverse } from "fp-ts/Array"
import { fromEntries, toEntries } from "fp-ts/Record"
import { startServer } from "./main"

const program = new Command()

const isDir = (directory: string) => statSync(directory).isDirectory()

const directoryCheck = pipe(existsSync, and(isDir))

const dirError = ([name, dir]: [string, string]) => {
  if (!existsSync(dir)) return new Error(`${name} directory not found: ${dir}`)
  return new Error(`${name} path is not a directory: ${dir}`)
}

const eitherDirExists = (directories: Record<string, string>) => {
  return pipe(
    directories,
    toEntries<string, string>,
    Atraverse(EApplicative)(
      EfromPredicate(([, dir]) => directoryCheck(dir), dirError),
    ),
    Emap(fromEntries),
  )
}

program
  .name("jbrowse-launch")
  .description("Serve a JBrowse2 application with local assets")
  .argument(
    "<root>",
    "directory where the JBrowse2 application entrypoint is located",
  )
  .argument("<assets>", "folder containing local assets to serve")
  .option("-p, --port <port>", "port to use for the server", "3000")
  .action((root: string, assets: string, { port }: { port: string }) => {
    pipe(
      { root, assets },
      eitherDirExists,
      Efold(
        ({ message }) => {
          console.error(message)
          process.exit(1)
        },
        ({ root, assets }) => {
          startServer(root, assets, { port: parseInt(port, 10) })
        },
      ),
    )
  })

program.parse()
