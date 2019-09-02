import { ICard, Container, ICreature } from "../cards/types"

const moveCardToFx = (card: ICard, container: Container) => {
    
    const player = card.player

    let filter = (x: ICard) => x !== card

    player.deck = player.deck.filter(filter)
    player.hand = player.hand.filter(filter)
    player.shieldzone = player.shieldzone.filter(filter)
    player.manazone = player.manazone.filter(filter)
    player.graveyard = player.graveyard.filter(filter)
    player.battlezone = player.battlezone.filter(filter)
    player.hiddenzone = player.hiddenzone.filter(filter)

    switch(container) {

        case Container.DECK: {
            player.deck.push(card)
            break
        }
        
        case Container.HAND: {
            player.hand.push(card)
            break
        }

        case Container.SHIELD_ZONE: {
            player.shieldzone.push(card)
            break
        }

        case Container.MANA_ZONE: {
            player.manazone.push(card)
            break
        }

        case Container.GRAVEYARD: {
            player.graveyard.push(card)
            break
        }

        case Container.BATTLE_ZONE: {
            player.battlezone.push(<ICreature>card)
            break
        }

        case Container.HIDDEN_ZONE: {
            player.hiddenzone.push(card)
            break
        }
        
        default: {
            break
        }

    }

}

export default moveCardToFx