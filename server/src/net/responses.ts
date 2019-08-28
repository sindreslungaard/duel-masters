import WebSocket from "ws"

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