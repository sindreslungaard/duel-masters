import events from "./events"

export const before = <K extends keyof events>(eventName: K, listener: (event: events[K], next: Function) => void) => {

}

export const on = <K extends keyof events>(eventName: K, listener: (event: events[K], next: Function) => void) => {

}

export const after = <K extends keyof events>(eventName: K, listener: (event: events[K], next: Function) => void) => {

}