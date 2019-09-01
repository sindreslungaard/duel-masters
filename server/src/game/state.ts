import { IDenormalizedCard } from "./cards/types/card";

export interface IPlayerState {
    playzone: IDenormalizedCard[],
    shieldzone: string[],
    manazone: IDenormalizedCard[],
    graveyard: IDenormalizedCard[],
    deck: number
}

export interface IMePlayerState extends IPlayerState {
    hand: IDenormalizedCard[]
}

export interface IStateUpdate {
    myTurn: boolean,
    hasAddedManaThisRound: boolean,
    me: IMePlayerState,
    opponent: IPlayerState
}