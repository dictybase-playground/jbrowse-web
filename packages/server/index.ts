import { Command } from "commander"
import { existsSync, statSync } from "fs"
import { pipe, flow } from "fp-ts/function"
import {
  map as Emap,
  fold as Efold,
  fromPredicate as EfromPredicate,
  Applicative as EApplicative,
  filterOrElse as EfilterOrElse,
} from "fp-ts/Either"
import { traverse as Atraverse } from "fp-ts/Array"
import { fromEntries, toEntries } from "fp-ts/Record"
import { startServer } from "./main"

const program = new Command()

const isDir = (directory: string) => statSync(directory).isDirectory()

const dirDoesNotExistError = ([name, dir]: [string, string]) => {
  return new Error(`$ Argument passed to ${name} path does not exist: ${dir}`)
}

const isNotDirectoryError = ([name, dir]: [string, string]) => {
  return new Error(`$ Argument passed to ${name} path is not a directory: ${dir}`)
}

const eitherDirExists = (directories: Record<string, string>) => {
  return pipe(
    directories,
    toEntries<string, string>,
    Atraverse(EApplicative)(
      flow(
        EfromPredicate(([, dir]) => existsSync(dir), dirDoesNotExistError),
        EfilterOrElse(([, dir]) => isDir(dir), isNotDirectoryError)
      )
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
          startServer(root, assets, { port })
        },
      ),
    )
  })

program.parse()
