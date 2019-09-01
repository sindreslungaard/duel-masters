import { IPlayer } from "../player"

const getOpponentFx = (player: IPlayer): IPlayer => {
    
    let p1 = player.match.player1
    let p2 = player.match.player2

    return p1 === player ? p1 : p2

}

export default getOpponentFx