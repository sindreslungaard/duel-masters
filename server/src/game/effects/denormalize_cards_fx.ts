import card, { IDenormalizedCard } from "../cards/types/card"

const denormalizeCardsFx = (cards: card[]): IDenormalizedCard[] => {

    let denormalizedCards: IDenormalizedCard[] = []

    for(let card of cards) {
        
        denormalizedCards.push({
            uid: card.id,
            virtualId: card.virtualId,
            name: card.name,
            civilization: card.civilization,
            tapped: card.tapped,
            canBePlayed: true
        })

    }

    return denormalizedCards

}

export default denormalizeCardsFx