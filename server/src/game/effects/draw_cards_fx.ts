import { IPlayer } from "../player"
import { card } from "../cards/types"

const drawCardsFx = (player: IPlayer, cardsToDraw: number): card[] => {

    let cards: card[] = []

    for(let i = 0; i < cardsToDraw; i++) {

        if(i > player.deck.length) {
            break
        }

        cards.push(player.deck.shift())

    }

    return cards

}

export default drawCardsFx 