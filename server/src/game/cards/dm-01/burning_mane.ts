import { ICreature, Civilization, Family } from "../types"
import { before, after } from "../../match"

const burningMane: ICreature = {

    id: "1d72eb3e-5185-449a-a16f-391bd2338343",
    name: "Burning Mane",
    civilization: Civilization.NATURE,
    family: Family.BEAST_FOLK,
    manaCost: 2,
    manaRequirement: [Civilization.NATURE],
    summoningSickness: true,

    setup(match, player) {
        
        before("turn-start", (event, next) => {
            // runs before turn start
        })

        after("turn-start", (event, next) => {
            // runs after turn start
        })

    }

}

export default burningMane