import { card } from "./types"
import { readdirSync, statSync } from "fs"
import { join } from "path"

const cardDatabase: {[key: string]: card} = {}

export const load = () => {

    let files = readdirSync(__dirname).filter(f => statSync(join(__dirname, f)).isDirectory())

}