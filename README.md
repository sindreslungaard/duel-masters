# duel-masters

A multiplayer browser based simulator for the [Duel Masters Trading Card Game](<https://duelmasters.fandom.com/wiki/Duel_Masters_(Card_Game)>)

> [!WARNING]
> We are not affiliated with the Duel Masters trademark or its rightful owner [Wizards of the Coast](https://wizards.com) and only strive to make it possible to play the original English version of the card game from 2004-2006 the same way you would be playing it with friends in real life.

#### What it is
- Aims to simulate how you would be playing the card game in real life
- Multiplayer & matchmaking
- Enforced rules for fair and competetive play
- Automation for all card effects and actions
- Deck building
- Currently 400+ cards implemented

#### What it is not
- A full fledged Duel Masters game with fancy animations and career modes similar to Heartstone, MTG Arena etc.

---


# Run with docker

```bash
# Login to the github package registry
docker login ghcr.io

# Run the container
docker run -d \
    --name duel-masters \
    --restart unless-stopped \
    -p 80:80 \
    -e port=80 \
    -e mongo_name=<mongodb_name> \
    -e mongo_uri=<mongodb_connection_string> \
    ghcr.io/sindreslungaard/duel-masters/production:latest
```

# Development

1. Clone the repository

2. Set up [MongoDB locally](https://www.mongodb.com/try/download/community) or use a [cloud provider](https://www.mongodb.com/atlas/database)

3. Create an `.env` file in the root directory based on the `.env.default` file in the repository
```
port=80
mongo_uri=mongodb://127.0.0.1:27017
mongo_name='duel-masters'
restart_after=
```

5. Build the webapp<br>
`cd webapp`<br>
`npm install`<br>
`npm run build`<br>
or<br>
`npm run watch` to rebuild every time a change is made

6. Run the application. If you're using vscode simply press `F5` or `Run -> Start Debugging`. To run manually use `go run cmd/duel-masters/main.go`. Note that debugging with vscode will automatically pick up environment variables from the `.env` file. If you're not using vscode's debug mode or if you are using a different editor you will have to make sure the necessary environment variables are set.

7. Navigate to `http://localhost` and create a user as well as a deck. To set the deck as a standard deck, find it in the database and change the `standard` field to `true`.

# Contributing
Please read through the [contribution guidelines](https://github.com/sindreslungaard/duel-masters/blob/master/CONTRIBUTIONS.md) before submitting a pull request

# Changelog

A changelog starting from 11/11/2021 can be found [here](https://github.com/sindreslungaard/duel-masters/blob/master/CHANGELOG.md)
