import { ICard } from "./types"
import { readdirSync, statSync } from "fs"
import { join } from "path"

const cardDatabase: {[key: string]: ICard} = {}

export const load = async () => {

    let dirs = readdirSync(__dirname).filter(f => statSync(join(__dirname, f)).isDirectory() && statSync(join(__dirname, f)))

    for(let set of dirs) {

        if(set === "types") {
            continue
        }

        let files = readdirSync(__dirname + "/" + set).filter(f => statSync(join(__dirname + "/" + set, f)).isFile())

        for(let file of files) {
            
            if(file.includes(".js.map")) {
                continue
            }

            let card = await import(__dirname + "/" + set + "/" + file)

            cardDatabase[card.default.id] = card.default
            
        } 

    }

}

export const getCard = (id: string): ICard => {
    return cardDatabase[id]
}