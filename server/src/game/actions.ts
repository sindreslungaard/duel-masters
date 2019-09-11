import { IPlayer } from './player'
import { sendError, sendWarning } from '../net/responses'
import Phase from './phase'
import { findCardFx, moveCardToFx, tapFx, stateUpdateFx, getOpponentFx } from './effects'
import { Container } from './cards/types'

export const tryAddCardToManazone = (player: IPlayer, virtualId: string) => {

    if(player.match.playerTurn !== player) {
        return sendWarning(
            player.client,
            "Wait for your turn"
        )
    }

    if(player.chargedMana) {
        return sendWarning(
            player.client,
            "You have already charged mana this turn"
        )
    }

    if(!player.match) {
        return sendWarning(
            player.client,
            "Match no longer available"
        )
    }

    if(player.match.phase !== Phase.CHARGE_STEP) {
        return sendWarning(
            player.client,
            `Mana can not be charged whilst in the ${player.match.phase} phase`
        )
    }

    let card = findCardFx([player], virtualId, Container.HAND)

    if(!card) {
        return sendWarning(
            player.client,
            "Card is not in your hand"
        )
    }

    moveCardToFx(card, Container.MANA_ZONE)
    tapFx(card)

    player.chargedMana = true

    stateUpdateFx([player, getOpponentFx(player)])

}