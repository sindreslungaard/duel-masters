import { ICard } from "../cards/types"

const tapFx = (card: ICard) => {
    card.tapped = true
}

export default tapFx