import { IPlayer } from "../player"
import { ICard } from "../cards/types"

const drawCardsFx = (player: IPlayer, cardsToDraw: number): ICard[] => {

    let cards: ICard[] = []

    for(let i = 0; i < cardsToDraw; i++) {

        if(i > player.deck.length) {
            break
        }

        cards.push(player.deck.shift())

    }

    return cards

}

export default drawCardsFx 