import { Command } from "commander"
import { existsSync } from "fs"
import { pipe } from "fp-ts/function"
import { fold as Efold, fromPredicate as EfromPredicate, Applicative as EApplicative } from "fp-ts/Either"
import { map as RNEAmap, traverse } from "fp-ts/ReadonlyNonEmptyArray"
import { startServer } from "../server/index"

const program = new Command()

const eitherDirExists = (dirTuple: [string, string]) => { 
  return pipe(
    dirTuple,
    EfromPredicate(
      ([, dir]) => existsSync(dir),
      ([name, dir]) => new Error(`Error: ${name} directory does not exist: ${dir}`)
    )
  )
}

program
  .name("jbrowse-launch")
  .description("Serve a JBrowse2 application with local assets")
  .argument("<root>", "directory where the JBrowse2 application entrypoint is located")
  .argument("<assets>", "folder containing local assets to serve")
  .action((root: string, assets: string) => {
    pipe(
      [["root", root], ["asset", assets]],
      traverse(EApplicative)(eitherDirExists),
      Efold(
        ({ message }) => {
          console.log(message)
          process.exit(1)
        },
        () => {
          startServer(root, assets)
        }
      )
    )
  })

program.parse()
