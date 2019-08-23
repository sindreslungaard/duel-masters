import events from "./events"
import { player } from "./player"

export interface match {
    player1: player,
    player2: player
}

export const createMatch = () => {

}

export const before = <K extends keyof events>(eventName: K, listener: (event: events[K], next: Function) => void) => {

}

export const after = <K extends keyof events>(eventName: K, listener: (event: events[K], next: Function) => void) => {

}