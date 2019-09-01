import WebSocket from "ws"
import { IMatch } from "./match"
import User, { IUser } from "../models/user"
import { IDeck } from "../models/deck"
import { card, creature } from "./cards/types"
import { shuffleDeckFx, drawCardsFx } from "./effects"

export interface IPlayer {
    client: WebSocket,
    user: IUser,
    match: IMatch,
    decks: IDeck[],
    deck?: card[],
    hand: card[],
    shieldzone: card[],
    manazone: card[],
    graveyard: card[],
    battlezone: creature[]
}

export const setupPlayer = (player: IPlayer) => {

    shuffleDeckFx(player)

    player.shieldzone = drawCardsFx(player, 5)
    player.hand = drawCardsFx(player, 5)

}