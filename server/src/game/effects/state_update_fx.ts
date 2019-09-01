import { IPlayer } from "../player"
import { IStateUpdate } from "../state"
import { denormalizeCardsFx, getOpponentFx } from "./index"
import { sendStateUpdate } from "../../net/responses"

const stateUpdateFx = (players: IPlayer[]) => {

    for(let player of players) {

        let opponent = getOpponentFx(player)

        let state: IStateUpdate = {
            myTurn: player.match.playerTurn === player,
            hasAddedManaThisRound: false, // TODO: decide whether or not to keep using this or refactor frontend to get rid of it serverside
            me: {
                playzone: denormalizeCardsFx(player.battlezone),
                shieldzone: player.shieldzone.map(x => x.virtualId),
                manazone: denormalizeCardsFx(player.manazone),
                hand: denormalizeCardsFx(player.hand),
                graveyard: denormalizeCardsFx(player.graveyard),
                deck: player.deck.length
            },
            opponent: {
                playzone: denormalizeCardsFx(opponent.battlezone),
                shieldzone: opponent.shieldzone.map(x => x.virtualId),
                manazone: denormalizeCardsFx(opponent.manazone),
                graveyard: denormalizeCardsFx(opponent.graveyard),
                deck: opponent.deck.length
            }
        }

        sendStateUpdate(player.client, state)

    }

}

export default stateUpdateFx