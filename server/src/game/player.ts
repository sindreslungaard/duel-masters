import WebSocket from "ws"
import { IMatch } from "./match"
import User, { IUser } from "../models/user"
import { IDeck } from "../models/deck"
import { card } from "./cards/types"

export interface IPlayer {
    client: WebSocket,
    user: IUser,
    match: IMatch,
    decks: IDeck[],
    deck?: IDeck,
    shieldzone: card[],
    manazone: card[],
    graveyard: card[],
    battlezone: card[]
}