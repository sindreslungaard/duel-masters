import { civilization, family } from "./index"
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

export default interface card {
    match?: IMatch,
    id: string,
    virtualId?: string,
    name: string,
    civilization: civilization,
    family: family,
    manaCost: number,
    manaRequirement: Array<civilization>,
    tapped?: boolean

    setup: (match: IMatch, player: IPlayer) => void
}