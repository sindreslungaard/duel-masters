/* 0ffdcae3-9db2-401b-8a82-dfad707b83cd */
/* dm01 */

self.set_name("Bronze-Arm Tribe")
self.set_power(1000)
self.set_civ("nature")
self.set_race("beast_folk")
self.set_mana_cost(3)
self.set_mana_req("nature")

self.use_trait("creature")

/* When you put this creature into the battlezone, put the top
   card of your deck into your manazone */    
match.on("card_moved", (event) => {

    if(event.card_id !== self.card_id || event.to !== "battlezone") {
        return;
    }

    const deck = match.get_deck(self.player)

    if(deck.length < 1) {
        return;
    }

    const card = deck[deck.length - 1]
    match.move_card(card, "manazone")
    match.send_chat("Server", `${card.name} was added to ${self.player.username}'s manazone from the top of their deck`)
    
})
