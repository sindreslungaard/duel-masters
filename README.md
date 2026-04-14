# duel-masters

A multiplayer browser based simulator for the [Duel Masters Trading Card Game](<https://duelmasters.fandom.com/wiki/Duel_Masters_(Card_Game)>)

> [!NOTE]
> This repository is not affiliated with the Duel Masters trademark or its rightful owner [Wizards of the Coast](https://wizards.com) and only strive to make it possible to play the original English version of the card game from 2004-2006 the same way you would be playing it with friends in real life.

#### What it is

- Aims to simulate how you would be playing the card game in real life
- Multiplayer & matchmaking
- Enforced rules for fair and competetive play
- Automation for all card effects and actions
- Deck building
- Currently 700+ cards implemented

#### What it is not

- A full fledged Duel Masters game with fancy animations and career modes similar to Heartstone, MTG Arena etc.

---

## Repository structure

- **`.dev`** - Local dev server for rapid development and testing
- **`packages`** - Shared packages, notably the duel interface component that can be injected into websites to make it possible to communicate with the simulator
- **`sim`** - The simulator server, also referred to as rule engine

## Running locally

### Prerequisites

- [Go 1.22 or higher](https://go.dev/dl/)
- [Node.js 18 or higher](https://nodejs.org/)

### Simulator server

```bash
cd sim
go run cmd/duel-masters/main.go
```

Alternatively use the predefined vscode debug launcher

### Dev server

```bash
cd .dev
npm install
npm run dev
```

The dev server will be available at `http://localhost:5173`

## Contributing

Please read through the [contribution guidelines](https://github.com/sindreslungaard/duel-masters/blob/master/CONTRIBUTING.md) before submitting a pull request

## Changelog

A changelog starting from 11/11/2021 can be found [here](https://github.com/sindreslungaard/duel-masters/blob/master/CHANGELOG.md)
