import { IPlayer } from "../player"
import { ICard } from "../cards/types"
import CardLocation from "../cards/types/card_location"

const findCardFx = (players: IPlayer[], virtualId: string, location?: CardLocation,): ICard => {
    
    let card

    for(let player of players) {

        const predicate = (x: ICard) => x.virtualId === virtualId
        
        switch(location) {

            case CardLocation.DECK: {
                card = player.deck.find(predicate)
                break
            }

            case CardLocation.HAND: {
                card = player.hand.find(predicate)
                break
            }

            case CardLocation.SHIELDZONE: {
                card = player.shieldzone.find(predicate)
                break
            }

            case CardLocation.MANAZONE: {
                card = player.manazone.find(predicate)
                break
            }

            case CardLocation.GRAVEYARD: {
                card = player.graveyard.find(predicate)
                break
            }

            case CardLocation.BATTLEZONE: {
                card = player.battlezone.find(predicate)
                break
            }

            case CardLocation.HIDDENZONE: {
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