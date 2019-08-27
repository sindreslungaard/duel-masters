import WebSocket from "ws"
import http from "http"
import noop from "../utils/noop"
import logger from "../utils/logger"

interface IClientAttachment {
    isAlive: boolean
}

const clientRepository = new Map<WebSocket, IClientAttachment>()

export const disposeClient = (client: WebSocket) => {
    clientRepository.delete(client)
    client.terminate()
}

export const connect = (web: http.Server) => {

    const wss = new WebSocket.Server({ server: web })

    wss.on("connection", client => {

        clientRepository.set(client, {
            isAlive: true
        })

        client.on("pong", () => clientRepository.get(client).isAlive = true)

        client.on("message", (data: WebSocket.Data) => {
            
            let request: any = {}

            try {
                request = JSON.parse(data.toString())
            } catch(err) {
                logger.debug("Unable to deserialize ws message", err)
                return
            }

            if(!request || !request.header) {
                logger.debug("Missing request header")
                return
            }

            try {
                parse(client, request)
            }
            catch(err) {
                logger.debug("Unable to parse message", err)
                return
            }

        })

        client.on("close", () => {
            disposeClient(client)
        })

    })

    setInterval(() => {
        wss.clients.forEach(client => {

            let clientAttachments = clientRepository.get(client)

            if(!clientAttachments) {
                return client.terminate()
            }

            if(!clientAttachments.isAlive) {

                clientAttachments.isAlive = false

                client.ping(noop)

            }

        })
    }, 30000)

}

const parse = (client: WebSocket, data: any) => {

    // TODO: this

}