import { civilization, family } from "./index"

export default interface card {
    name: string,
    civilization: civilization,
    family: family,
    manaCost: number,
    manaRequirement: Array<civilization>,

    setup: Function
}