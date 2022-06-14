# Changelog

All notable changes to this project as of 11/11/2021 will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
Fixed an issue with Sniper Mosquito allowing it to return any amount of cards from manazone to hand

### Added

### Changed
Changed max deck limit from 30 to 50

## v2.7 - 22/05/2022

### Fixed

- Fixed an issue where "Trox, General of Destruction"'s ability crashes the match

### Added

- Added alphabetical sorting to deck builder categories (thanks @vasyop)

## v2.6 - 12/05/2022

### Fixed

- Fixed an issue where the card selection popup would not re-appear after reconnecting to a duel
- Fixed an issue that made it possible to cause multiple player actions/events to be handled at the same time, causing weird behaviour in matches
- Fixed an issue where the choose deck event could be executed after a match had started
- Removed suicide effect from Skeleton Soldier, the Defiled (thanks @vasyop)

## v2.5 - 28/04/2022

### Fixed

- Fixed an issue where double breakers sometimes only broke 1 shield while Alcadeias, Lord of Spirits was in play
- Fixed an issue where Gregoria, Princess of War did not properly give and remove power bonus and blocker to demon commands
- Fixed an issue where Creeping Plague did not remove slayer at the end of the turn, as well as giving slayer to the opponent
- Fixed an issue where Vampire Silphy, Burst Shot and Searing Wave would destroy cards sequentially, checking the next cards power after every destroy (thanks @vasyop)
- Fixed an issue with drag and drop and selecting cards for firefox (thanks @vasyop)
- Correctly display the graveyard's owner when spectating (thanks @vasyop)
- Fixed an issue causing matches to not be shown in the match list for a certain period of time

### Added

- Added ban system
- Added mute system (thanks @vasyop)
- Added option for no upside down cards to settings (thanks @vasyop)
- Added help text that suggests select cards by dragging (thanks @vasyop)

### Changed

- Duels can now be created with an empty name, in that case a name will be randomly chosen from a preset of names based on terms from the anime

## v2.4 - 19/04/2022

### Fixed

- Fixed missing "can't be blocked" effect for Sea Slug
- Fixed missing "can't attack players" effect for Cannoneer Bargon

### Added

- New card "Ballus, Dogfight Enforcer Q"
- New card "Split-Head Hydroturtle Q"
- New card "Bladerush Skyterror Q"
- New card "Ruthless Skyterror"
- New card "Death Cruzer, the Annihilator"
- Dismiss large card(s) by clicking the overlay (thanks @vasyop)
- Select cards by dragging over them (thanks @vasyop)
- Attack shields and creatures by drag and drop (thanks @vasyop)

## [v2.3] - 09/04/2022

### Fixed

- Removed doublebreaker from "Supporting Tulip"
- Corrected Masked Pomegranate's mana cost from 5 to 4
- Corrected Muramasa, Duke of Blades' power from 3000 to 2000
- Fixed an issue where Silver Axe kept adding mana when clicking "attack creature" while opponent had no creatures in the battle zone
- Fixed an issue where Horrid Worm's effect did not apply when the opponent has a blocker but chose to not block
- Fixed a typo in "Dogarn, the Marauder"'s name

### Added

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

## [v2.2] - 21/01/2022

### Fixed

- Corrected "Rayla, Truth Enforcer"'s power from 5500 to 3000 (thanks @gaurav281296)
- Fixed an issue where "Horrid Worm"'s effect was not working when attacking a shield trigger or opponent had blockers but didn't block (thanks @gaurav281296)

## [v2.1] - 09/12/2021

### Added

- Added spectating

## [v2.0.32] - 11/11/2021

### Fixed

- Fixed an issue where Parasite Worm would discard the opponent's cards after breaking shields, causing it to sometimes discard the broken shield.
- Fixed an issue where Shadow Moon Cursed Shade didn't remove the buffs it gave other darkness creature after being removed from the battle zone.

### Added

- Added CHANGELOG.md
- Added new event hook for when an attack is confirmed to be happening
- Added timestamps to lobby chat messages

### Changed

- Updated README with steps to set up the project for local development
