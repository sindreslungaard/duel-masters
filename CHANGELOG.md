# Changelog

All notable changes to this project as of 11/11/2021 will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- Fixed missing "can't be blocked" effect for Sea Slug
- Fixed missing "can't attack players" effect for Cannoneer Bargon

### Added

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
