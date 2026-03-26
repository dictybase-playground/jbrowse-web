import { Command } from "commander"
import { existsSync, PathLike } from "fs"
import { pipe } from "fp-ts/function"
import { map as Emap, fold as Efold, fromPredicate as EfromPredicate, Applicative as EApplicative } from "fp-ts/Either"
import { traverse } from "fp-ts/ReadonlyNonEmptyArray"
import { traverse as Atraverse } from "fp-ts/Array"
import { fromEntries, toEntries, traverseWithIndex as RtraverseWithIndex } from "fp-ts/Record"
import { startServer } from "../server/index"

const program = new Command()

const eitherDirExists = (directories: Record<string, string>) => { 
  return pipe(
    directories,
    toEntries<string, string>,
    Atraverse(EApplicative)
    (EfromPredicate(
      ([,dir]) => existsSync(dir),
      ([name, dir]) => new Error(`Error: ${name} directory does not exist: ${dir}`)
    )),
    Emap(fromEntries)
  )
}

program
  .name("jbrowse-launch")
  .description("Serve a JBrowse2 application with local assets")
  .argument("<root>", "directory where the JBrowse2 application entrypoint is located")
  .argument("<assets>", "folder containing local assets to serve")
  .action((root: string, assets: string) => {
    pipe(
      { root, assets },
      eitherDirExists,
      Efold(
        ({ message }) => {
          console.error(message)
          process.exit(1)
        },
        ({ root, assets }) => {
          startServer(root, assets)
        }
      )
    )
  })

program.parse()
