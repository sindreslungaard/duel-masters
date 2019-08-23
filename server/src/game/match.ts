import events from "./events"
import { IPlayer } from "./player"
import shortid from "shortid"
import Phase from "./phase"

export interface IMatch {
    id: string,
    inviteId: string,
    name: string,
    description: string,
    players: Array<IPlayer>,
    phase: Phase,
    playerTurn?: IPlayer
}

let matches: Array<IMatch> = []

export const createMatch = (name: string, description: string): IMatch => {

    let match = {
        id: shortid.generate(),
        inviteId: shortid.generate(), 
        name,
        description,
        players: new Array<IPlayer>(),
        phase: Phase.IDLE
    }

    return match

}

export const before = <K extends keyof events>(eventName: K, listener: (event: events[K], next: Function) => void) => {

}

export const after = <K extends keyof events>(eventName: K, listener: (event: events[K], next: Function) => void) => {

}