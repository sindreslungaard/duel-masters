import { civilization, family } from "./index"

export default interface card {
    uid: string,
    name: string,
    civilization: civilization,
    family: family,
    manaCost: number,
    manaRequirement: Array<civilization>,

    setup: Function
}