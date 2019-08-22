import { creature, civilization, family } from "../types"
import { before } from "../../match"

const burningMane: creature = {

    name: "Burning Mane",
    civilization: civilization.NATURE,
    family: family.BEAST_FOLK,
    manaCost: 2,
    manaRequirement: [civilization.NATURE],

    setup() {
        
        before('turn-start', (event, next) => {
            
        })

    }

}

export default burningMane