import { Civilization, Family } from "./index"
import { IMatch } from "../../match"
import { IPlayer } from './../../player'

export interface IDenormalizedCard {
    uid: string,
    virtualId: string,
    name: string,
    civilization: string,
    tapped: boolean,
    canBePlayed: boolean
}

export default interface ICard {
    match?: IMatch,
    id: string,
    virtualId?: string,
    name: string,
    civilization: Civilization,
    family: Family,
    manaCost: number,
    manaRequirement: Array<Civilization>,
    tapped?: boolean

    setup: (match: IMatch, player: IPlayer) => void
}