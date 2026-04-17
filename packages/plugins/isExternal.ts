import jbrowseGlobals from '@jbrowse/core/ReExports/list'
import { isAbsolute } from "path"

const isExternal = (id: string) => {
    if (id.startsWith('regenerator-runtime')) {
        return false;
    }
    return !id.startsWith('.') && isAbsolute(id) && jbrowseGlobals.includes(id)
}

export { isExternal }
