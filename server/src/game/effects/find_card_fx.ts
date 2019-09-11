import { IPlayer } from "../player"
import { ICard, Container } from "../cards/types"

const findCardFx = (players: IPlayer[], virtualId: string, location?: Container): ICard => {
    
    let card

    for(let player of players) {

        const predicate = (x: ICard) => x.virtualId === virtualId
        
        switch(location) {

            case Container.DECK: {
                card = player.deck.find(predicate)
                break
            }

            case Container.HAND: {
                card = player.hand.find(predicate)
                break
            }

            case Container.SHIELD_ZONE: {
                card = player.shieldzone.find(predicate)
                break
            }

            case Container.MANA_ZONE: {
                card = player.manazone.find(predicate)
                break
            }

            case Container.GRAVEYARD: {
                card = player.graveyard.find(predicate)
                break
            }

            case Container.BATTLE_ZONE: {
                card = player.battlezone.find(predicate)
                break
            }

            case Container.HIDDEN_ZONE: {
                card = player.hiddenzone.find(predicate)
                break
            }

            default: {
                card = player.deck.find(predicate)
                card = player.hand.find(predicate)
                card = player.shieldzone.find(predicate)
                card = player.manazone.find(predicate)
                card = player.graveyard.find(predicate)
                card = player.battlezone.find(predicate)
                card = player.hiddenzone.find(predicate)
                break
            }

        }

    }

    return card

}

export default findCardFx