import { creature } from "../cards/types"

const resolveSummoningSicknessFx = (creature: creature) => {
    creature.summoningSickness = false
}

export default resolveSummoningSicknessFx