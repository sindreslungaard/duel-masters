import WebSocket from "ws"
import { IDeck } from "../models/deck"

const send = (client: WebSocket, header: string, data?: any) => {
    data.header = header
    client.send(JSON.stringify(data))
}

export const sendError = (client: WebSocket, message: string, action?: string) => {
    send(client, "error", { message, action})
}

export const sendHello = (client: WebSocket) => {
    send(client, "hello")
}

export const sendChooseDeck = (client: WebSocket, decks: IDeck[]) => {
    send(client, "choose_deck", { decks })
}