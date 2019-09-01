import { IMatch } from "../match"
import Phase from "../phase"

const setPhaseFx = (match: IMatch, phase: Phase) => {
    match.phase = phase
}

export default setPhaseFx