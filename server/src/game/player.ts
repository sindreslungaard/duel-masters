import WebSocket from "ws"

export interface player {
    name: string,
    client: WebSocket
}