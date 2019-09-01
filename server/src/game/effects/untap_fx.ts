import { card } from "../cards/types"

const untapFx = (card: card) => {
    card.tapped = false
}

export default untapFx