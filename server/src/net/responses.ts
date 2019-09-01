import WebSocket from "ws"
import { IDeck } from "../models/deck"
import { IStateUpdate } from "../game/state"

const send = (client: WebSocket, header: string, data?: any) => {
    if(!data) data = {}
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

export const sendWarning = (client: WebSocket, message: string) => {
    send(client, "warning", message)
}

export const sendStateUpdate = (client: WebSocket, state: IStateUpdate) => {
    send(client, "state_update", state)
}