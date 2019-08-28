import WebSocket from "ws"
import { IMatch } from "./match"
import User, { IUser } from "../models/user"

export interface IPlayer {
    client: WebSocket,
    user: IUser,
    match: IMatch
}