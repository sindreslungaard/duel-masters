import { IPlayer } from "../player"
import shuffle from "../../utils/shuffle"

const shuffleDeckFx = (player: IPlayer): void => {
    if(!player.deck) {
        return
    }
    
    shuffle(player.deck)
}

export default shuffleDeckFx