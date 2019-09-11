import WebSocket from "ws"
import { IMatch } from "./match"
import User, { IUser } from "../models/user"
import { IDeck } from "../models/deck"
import { ICard, ICreature } from "./cards/types"
import { shuffleDeckFx, drawCardsFx } from "./effects"

export interface IPlayer {
    client: WebSocket,
    user: IUser,
    match: IMatch,
    decks: IDeck[],
    deck?: ICard[],
    hand: ICard[],
    shieldzone: ICard[],
    manazone: ICard[],
    graveyard: ICard[],
    battlezone: ICreature[],
    hiddenzone: ICard[],
    chargedMana: boolean
}

export const setupPlayer = (player: IPlayer) => {

    shuffleDeckFx(player)

    player.shieldzone = drawCardsFx(player, 5)
    player.hand = drawCardsFx(player, 5)

}