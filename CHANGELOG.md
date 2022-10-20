# Changelog

All notable changes to this project as of 11/11/2021 will be documented in this file.

## [Unreleased]

- Fixed an issue where creatures could attack while tapped
- New card "Smash Horn Q"
- New card "Nocturnal Giant"
- New card "Moon Horn"

## 03/09/2022

- Fix for match sometimes not ending when a player runs out of cards in their deck if the last card was moved by a card ability
- Updated Enchanted Soil to only let you retrieve creatures from your graveyard
- Fixed an issue where Elf-X could keep decreasing a cards mana cost by cancelling playing a card before selecting mana
- Fixed an issue where Sieg Balicula, the Intense would not give blocker to your other light creatures if they are tapped
- Fixed an issue where Pouch Shell did not keep the tapped state of the target card when its evolution was removed
- Fixed an issue where Whisking Whirlwind did not untap all creatures at the end of the turn
- Fixed an issue causing Jack Viper, Shadow of Doom's effect to apply to non-darkness creatures too
- Fixed an issue where Bolzard Dragon could use its effect infinite times by attacking creatures while opponent had no creatures
- Fixed an issue where Dark Reversal and similar shield triggers did not work properly during opponent's turn
- Fixed an issue where Three-Eyed Dragonfly could cancel the attack and still get the benefits from its effect
- Fixed an issue where Re Bil, Seeker of Archery's effect did not go away when it was removed from the battlezone
- Fixed an issue where Diamond Cutter would not apply to creatures summoned after the spell was cast
- Fixed missing effect for Gigamantis
- Fixed an issue with Artisan Picora's effect that caused matches to freeze

## 18/08/2022

- Fixed an issue with Sniper Mosquito allowing it to return any amount of cards from manazone to hand
- Fixed an issue with Whisking Whirlwind where it wouldn't untap your creatures at the end of the turn (thanks @Fightlapa)
- Added sorting of mana cost in deck builder (thanks @Fightlapa)
- New card "Skullsweeper Q" (thanks @Fightlapa)
- New card "Enchanted Soil" (thanks @Fightlapa)
- New card "Avalanche Giant" (thanks @Fightlapa)
- Changed max deck limit from 30 to 50
- Cards are now previewed while hovering the card name in the deck builder (thanks @Fightlapa)
- Updated Mongrel Man's ability to be optional (thanks @Fightlapa)

## 22/05/2022

- Fixed an issue where "Trox, General of Destruction"'s ability crashes the match
- Added alphabetical sorting to deck builder categories (thanks @vasyop)

## 12/05/2022

- Fixed an issue where the card selection popup would not re-appear after reconnecting to a duel
- Fixed an issue that made it possible to cause multiple player actions/events to be handled at the same time, causing weird behaviour in matches
- Fixed an issue where the choose deck event could be executed after a match had started
- Removed suicide effect from Skeleton Soldier, the Defiled (thanks @vasyop)

## 28/04/2022

- Fixed an issue where double breakers sometimes only broke 1 shield while Alcadeias, Lord of Spirits was in play
- Fixed an issue where Gregoria, Princess of War did not properly give and remove power bonus and blocker to demon commands
- Fixed an issue where Creeping Plague did not remove slayer at the end of the turn, as well as giving slayer to the opponent
- Fixed an issue where Vampire Silphy, Burst Shot and Searing Wave would destroy cards sequentially, checking the next cards power after every destroy (thanks @vasyop)
- Fixed an issue with drag and drop and selecting cards for firefox (thanks @vasyop)
- Correctly display the graveyard's owner when spectating (thanks @vasyop)
- Fixed an issue causing matches to not be shown in the match list for a certain period of time
- Added ban system
- Added mute system (thanks @vasyop)
- Added option for no upside down cards to settings (thanks @vasyop)
- Added help text that suggests select cards by dragging (thanks @vasyop)
- Duels can now be created with an empty name, in that case a name will be randomly chosen from a preset of names based on terms from the anime

## 19/04/2022

- Fixed missing "can't be blocked" effect for Sea Slug
- Fixed missing "can't attack players" effect for Cannoneer Bargon
- New card "Ballus, Dogfight Enforcer Q"
- New card "Split-Head Hydroturtle Q"
- New card "Bladerush Skyterror Q"
- New card "Ruthless Skyterror"
- New card "Death Cruzer, the Annihilator"
- Dismiss large card(s) by clicking the overlay (thanks @vasyop)
- Select cards by dragging over them (thanks @vasyop)
- Attack shields and creatures by drag and drop (thanks @vasyop)

## 09/04/2022

- Removed doublebreaker from "Supporting Tulip"
- Corrected Masked Pomegranate's mana cost from 5 to 4
- Corrected Muramasa, Duke of Blades' power from 3000 to 2000
- Fixed an issue where Silver Axe kept adding mana when clicking "attack creature" while opponent had no creatures in the battle zone
- Fixed an issue where Horrid Worm's effect did not apply when the opponent has a blocker but chose to not block
- Fixed a typo in "Dogarn, the Marauder"'s name
- Card preview and new duel dialog can now be closed by clicking outside the popup (thanks @AstroProjection)
- Added new filter "race" in deck builder
- Add button to stop spectating
- New card "Twin-Cannon Skyterror"
- New card "Bolgash Dragon"
- New card "La Guile, Seeker of Skyfire"
- New card "Sea Slug"
- New card "Rikabu, the Dismantler"
- New card "Cannoneer Bargon"
- New card "Bombat, General of Speed"
- New card "Billion-Degree Dragon"

## 21/01/2022

- Corrected "Rayla, Truth Enforcer"'s power from 5500 to 3000 (thanks @gaurav281296)
- Fixed an issue where "Horrid Worm"'s effect was not working when attacking a shield trigger or opponent had blockers but didn't block (thanks @gaurav281296)

## 09/12/2021

- Added spectating

## 11/11/2021

- Fixed an issue where Parasite Worm would discard the opponent's cards after breaking shields, causing it to sometimes discard the broken shield.
- Fixed an issue where Shadow Moon Cursed Shade didn't remove the buffs it gave other darkness creature after being removed from the battle zone.
- Added CHANGELOG.md
- Added new event hook for when an attack is confirmed to be happening
- Added timestamps to lobby chat messages
- Updated README with steps to set up the project for local development
