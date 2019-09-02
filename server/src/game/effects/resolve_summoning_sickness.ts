import { ICreature } from "../cards/types"

const resolveSummoningSicknessFx = (creature: ICreature) => {
    creature.summoningSickness = false
}

export default resolveSummoningSicknessFx