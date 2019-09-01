import { creature, civilization, family } from "../types"
import { before, after } from "../../match"

const burningMane: creature = {

    id: "1d72eb3e-5185-449a-a16f-391bd2338343",
    name: "Burning Mane",
    civilization: civilization.NATURE,
    family: family.BEAST_FOLK,
    manaCost: 2,
    manaRequirement: [civilization.NATURE],

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