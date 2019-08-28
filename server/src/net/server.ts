import WebSocket from "ws"
import http from "http"
import noop from "../utils/noop"
import logger from "../utils/logger"
import { sendError, sendHello } from "./responses"
import User, { IUser } from "../models/user"
import { addPlayer, playerChooseDeck } from "../game/match"
import { IPlayer } from "../game/player"

interface IClientAttachment {
    isAlive: boolean,
    user?: IUser,
    player?: IPlayer
}

const clientRepository = new Map<WebSocket, IClientAttachment>()

export const getClientAttachments = (client: WebSocket) => {
    return clientRepository.get(client)
}

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

const parse = async (client: WebSocket, data: any) => {

    switch(data.header) {

        case "connect": {

            const token = data.token

            if(!token) {
                return sendError(client, "Missing authorization token", "login")
            }

            let user = await User.findOne({'sessions': { $elemMatch: { token: token }}})

            if(!user) {
                return sendError(client, "Unauthorized", "login")
            }

            let clientAttachments = clientRepository.get(client)

            if(!clientAttachments) {
                return client.terminate()
            }

            clientAttachments.user = user

            return sendHello(client)

        }

        case "join_match": {

            if(!data.matchUid) {
                return sendError(client, "Missing match id", "overview")
            }

            if(!data.inviteId) {
                return sendError(client, "Missing invite id", "overview")
            }

            let clientAttachments = clientRepository.get(client)

            if(!clientAttachments) {
                return client.terminate()
            }

            if(!clientAttachments.user) {
                return client.terminate()
            }

            return addPlayer(client, clientAttachments.user, data.matchUid, data.inviteId)

        }

        case "choose_deck": {

            if(!data.uid) {
                return sendError(client, "Missing deck id")
            }

            let clientAttachments = clientRepository.get(client)

            if(!clientAttachments) {
                return client.terminate()
            }

            if(!clientAttachments.player) {
                return client.terminate()
            }

            playerChooseDeck(clientAttachments.player, data.uid)
        }

        default: {
            logger.debug("Ws message received with no matching header", data)
        }

    }

}