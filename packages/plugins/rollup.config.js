import globals from '@jbrowse/core/ReExports/list'
import { createRollupConfig } from '@jbrowse/development-tools'
import { inspect } from "util"

const config = createRollupConfig(globals, {
  includeESMBundle: true
})

console.log(inspect(config[1]))

export default config
