import WebSocket from "ws"

export interface IPlayer {
    name: string,
    client: WebSocket
}