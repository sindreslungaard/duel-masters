import { ICard } from "../cards/types"

const untapFx = (card: ICard) => {
    card.tapped = false
}

export default untapFx