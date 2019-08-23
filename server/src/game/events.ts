interface TurnStartEvent {

}

interface TurnEndEvent {

}

export default interface Events {
    "turn-start": TurnStartEvent,
    "turn-end": TurnEndEvent
}