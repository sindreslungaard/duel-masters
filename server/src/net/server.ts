import WebSocket from "ws"
import http from "http"

interface IClientAttachment {
    isAlive: boolean
}

let clientRepository: Map<WebSocket, IClientAttachment>

export const connect = (web: http.Server) => {

    const wss = new WebSocket.Server({ server: web })

    wss.on("connection", client => {

        clientRepository.set(client, {
            isAlive: true
        })

    })

}