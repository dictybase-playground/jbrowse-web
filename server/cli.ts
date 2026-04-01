import { Command } from "commander"
import { existsSync } from "fs"
import { pipe } from "fp-ts/function"
import {
  map as Emap,
  fold as Efold,
  fromPredicate as EfromPredicate,
  Applicative as EApplicative,
} from "fp-ts/Either"
import { traverse as Atraverse } from "fp-ts/Array"
import { fromEntries, toEntries } from "fp-ts/Record"
import { startServer } from "./main"

const program = new Command()

const eitherDirExists = (directories: Record<string, string>) => {
  return pipe(
    directories,
    toEntries<string, string>,
    Atraverse(EApplicative)(
      EfromPredicate(
        ([, dir]) => existsSync(dir),
        ([name, dir]) =>
          new Error(`Error: ${name} directory does not exist: ${dir}`),
      ),
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
  .option("-p --port [port]", "port to use for the server", "3000")
  .action((root: string, assets: string, { port }: { port: number }) => {
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
