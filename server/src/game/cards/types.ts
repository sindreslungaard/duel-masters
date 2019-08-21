export enum civilization {
    FIRE = "fire",
    WATER = "water",
    NATURE = "nature",
    LIGHT = "light",
    DARKNESS = "darkness"
}

export interface card {
    manaCost: number,
    manaRequirement: Array<civilization>,

    setup: Function
}

export interface creature extends card {

}

export interface spell extends card {

}