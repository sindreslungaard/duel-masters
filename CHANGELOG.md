# Changelog

## [Unreleased]

- Updated codebase to support cards with multiple races (thanks @tydeFriz)

## 24/08/2023

- When creating a new duel, you will now stay in the lobby until the match starts. When someone joins, the host will have the option to either start the match or kick the guest. Kicking prevents the user from joining the duel again.

## 09/08/2023

- Fixed an issue where some new cards could not be saved to deck
- Updated the card displayed in the graveyard pile to be the last discarded card instead of the first
- Fixed an issue where the lobby chat would not scroll all the way down on page load if there were any pinned messages

## 08/08/2023
- New card "Gigaling Q"
- New card "Wisp Howler, Shadow of Tears"
- New card "Gigakail"
- New card "Steel-Turret Cluster"
- New card "Le Quist, the Oracle"
- Fixed an issue where cards would keep their stat changes between turns

## 05/08/2023
- Right clicking cards in "select card" popups during matches will now preview the card
- Updated the backside card image used in match action/popups to the high quality version
- Updated the color scheme of the site
- Moved "Copy deck list" button to the other deck-action icons
- Updated "Copy deck list" popup to always be centered on the screen
- Fixed an issue where long lobby chat messages would overflow and create a horizontal scrollbar
- Increased max decks per user from 50 to 200
- Refactored database models and decreased the file size of decks on disk
- New card "Cutthroat Skyterror" (thanks @Zayberex)
- New card "Cursed Pincher" (thanks @Zayberex)
- New card "Junkatz, Rabid Doll" (thanks @Zayberex)
- New card "Lupa, Poison-Tipped Doll" (thanks @Zayberex)
- New card "Pyrofighter Magnus" (thanks @Zayberex)
- New card "Bazagazeal Dragon" (thanks @Zayberex)
- New card "Valiant Warrior Exorious" (thanks @Zayberex)
- New card "Automated Weaponmaster Machai" (thanks @Zayberex)
- New card "Mighty Bandit, Ace of Thieves" (thanks @Zayberex)
- New card "Gigagriff" (thanks @Zayberex)
- New card "Carrier Shell" (thanks @Zayberex)
- New card "Slumber Shell" (thanks @Zayberex)

## 28/06/2023

- Fixed an issue where players sometimes could get more cards during a match than their deck contained
- Added functionality to view this changelog in the header menu
- Correct "Masked Pomegranate"'s mana cost from 4 to 5
- Fixed an issue where "Plasma Chaser"'s ability would trigger even if the attack was cancelled
- Removed shield trigger effect from "Recon Operation"

## 10/05/2023

- Added password recovery functionality
- Added a coin toss system to choose which player starts (thanks @fabianTMC)
- Matches are no longer closed when the non-host leaves before the match has started (thanks @fabianTMC)
- New card "Gigazoul" (thanks @fabianTMC)
- Added sorting to selected deck (thanks @fabianTMC)
- Added button to choose a random deck in the deck selection screen (thanks @fabianTMC)
- Added chat message to show who started the game and who joined (thanks @fabianTMC)
- New card "Crystal Jouster" (thanks @fabianTMC)
- New card "Ultra Mantis, Scourge of Fate" (thanks @fabianTMC)
- New card "Craze Valkyrie, the Drastic" (thanks @fabianTMC)
- Added copy deck list button to make tourney signups easier (thanks @fabianTMC)
- Updated the quality of the backside card image
- Fixed an incorrect spelling of "Toel, Vizier of Hope"'s name
- Fixed an issue causing Bone Piercer's chat message to display its name incorrectly
- Updated Bone Piercer to show the opponent a waiting modal while deciding which card to move
- Fixed an issue where "Abush Scorpion" and "Obsidian Scarab" could not attack on your next turn if moved from your manazone during your turn
- New card "Raptor Fish"

## 14/02/2023

- New Card "Sinister General Damudo" (thanks @fabianTMC)
- Added change password functionality to settings page
- Add button when checking a public deck to copy it to your decks (thanks @pablopenna)
- Fix for Suicide effect (e.g. Bloody Squito) not triggering when battling creatures with "return to... when would be destroyed"
- Fix for cards like "Mist Rias" and "Mongrel Man" losing their effect when another copy of those types of cards were destroyed
- Players can now set a custom playmat from the settings page
- New card "Neon Cluster"

## 31/01/2023

- Added shield numbers to all shields (thanks @fabianTMC)
- Fix for "Gregoria Princess Of War" not giving power and blocker to opponent's creatures as well (thanks @fabianTMC)
- Fix for "Three Eyed Dragonfly" not letting you decide to use its effect before selecting target to attack (thanks @fabianTMC)

## 07/01/2023

- New card "Pokolul" (thanks @jyotiskaghosh)
- Bugfixes for Ambush Scorpion, Obsidian Scarab, Mist Rias, Mongrel Man, Bone Spider, Mega Detonator (thanks @jyotiskaghosh)
- Improved image quality of all cards

## 04/01/2023

- New card "Solidskin Fish" (thanks @jyotiskaghosh)
- New card "Spikestrike Ichthys Q" (thanks @jyotiskaghosh)
- New card "Blazosaur Q" (thanks @jyotiskaghosh)
- New card "Gallia Zohl IronGuardian Q" (thanks @jyotiskaghosh)
- New card "Jewel Spider" (thanks @jyotiskaghosh)
- New card "Horned Mutant" (thanks @jyotiskaghosh)
- New card "Kip Chippotto" (thanks @jyotiskaghosh)
- New card "Balloonshroom Q" (thanks @jyotiskaghosh)
- Various bugfixes (thanks @jyotiskaghosh)

## 31/12/2022

- Fixed an issue where "Sieg Balicula, the Intense" would not properly remove the blocker effect it had given out. Also fixes cards like Scarlet Skyterror's effect not removing the blockers given by Sieg
- New card "King Tsunami" (thanks @jyotiskaghosh)
- New card "King Mazelan" (thanks @jyotiskaghosh)
- New card "Aqua Surfer" (thanks @jyotiskaghosh)
- New card "Vashuna, Sword Dancer" (thanks @jyotiskaghosh)
- New card "Slime Veil" (thanks @jyotiskaghosh)
- New card "Brutal Charge" (thanks @jyotiskaghosh)
- New card "Miracle Quest" (thanks @jyotiskaghosh)
- New card "Divine Riptide" (thanks @jyotiskaghosh)
- New card "Cataclysmic Eruption" (thanks @jyotiskaghosh)
- New card "Thunder Net" (thanks @jyotiskaghosh)
- New card "Bloodwing Mantis" (thanks @jyotiskaghosh)
- New card "Scissor Scarab" (thanks @jyotiskaghosh)
- New card "Ambush Scorpion" (thanks @jyotiskaghosh)
- New card "Obsidian Scarab" (thanks @jyotiskaghosh)

## 13/11/2022

- Fixed an issue where creatures could attack while tapped
- New card "Smash Horn Q"
- New card "Nocturnal Giant"
- New card "Moon Horn"
- New card "Scheming Hands"
- New card "Lurking Eel"
- New card "Crow Winger"
- New card "Snork La, Shrine Guardian"
- New card "Cyclone panic"
- New card "Syrius, Firmament Elemental" (thanks @jyotiskaghosh)
- New card "Syforce, Aurora Elemental" (thanks @jyotiskaghosh)
- New card "Kulus, Soulshine Enforcer" (thanks @jyotiskaghosh)
- New card "Calgo, Vizier of Rainclouds" (thanks @jyotiskaghosh)
- New card "La Byle, Seeker of the Winds" (thanks @jyotiskaghosh)
- New card "Glory Snow" (thanks @jyotiskaghosh)
- Add functionality for pinning lobby chat messages
- Add functionality for disabling match creation

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

---

Changes prior to 11/11/2021 was not documented using this format
