# duel-masters

duel-masters is a multiplayer simulator for the [Duel Masters Trading Card Game](<https://duelmasters.fandom.com/wiki/Duel_Masters_(Card_Game)>) for play in the browser.

It aims to simulate how you would be playing the card game in real life, but with enforced rules and automations for the effects of each individual card.

## Run with docker

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

# Setting up for development

Run the server, webapp and db locally (recommended)

<details>

1. Fork the `duel-masters` repo on GitHub.
2. Clone your fork locally:

```
git clone https://github.com/sindreslungaard/duel-masters.git
```

3. Set up [MongoDB locally](https://www.mongodb.com/try/download/community) or use a [cloud provider](https://www.mongodb.com/atlas/database).
4. Create a `.env` file based on the `.env.default` file in the repository. Example:

```
port=80
mongo_uri=mongodb://127.0.0.1:27017
mongo_name='duel-masters'
restart_after=
```

5. Navigate to the `webapp` directory and run `npm install`. Then run either `npm run build` or `npm run watch` to build or watch the files.

6. Run the application. If you're using Vscode simply hit F5 or `Run -> Start Debugging`. To run manually use `go run cmd/duel-masters/main.go`. Note that debugging with vscode will automatically pick up environment variables from the `.env` file. If you're not using vscode's debug mode or if you use a different editor you will have to set up the environment variables yourself.

7. Go to `http://localhost` and create a user as well as a deck. To set the deck as a standard deck, find it in MongoDB and change the `standard` field to `true`.

</details>

<br><br>

Run only the frontend locally

<details>

1. git clone https://github.com/sindreslungaard/duel-masters.git

2. Navigate to the `webapp` directory and run `npm install` and `npm run serve`

3. Override your host config

```
localStorage.setItem(
    "config",
    JSON.stringify({
        host: "shobu.io",
        ws_protocol: "wss://",
        api: "https://shobu.io/api"
    })
)
```

</details>

<br><br>

Run the server, frontend and db locally with docker compose

<details>

1. Fork the `duel-masters` repo on GitHub.
2. Clone your fork locally:

```
git clone https://github.com/sindreslungaard/duel-masters.git
```

3. go into your cloned directory
4. Start the containers via:

```bash
docker compose up --build
```

5. Go to [http://localhost:9080](http://localhost:9080) and create a user and a deck

</details>

<br><br>

### Development tips

For testing cards locally there is a neat chat command that spawns in any card to the hand of the persons turn it is even if it's not in your deck `/add {cardUid}` (e.g `/add 4459f97b-8927-4231-a11d-f4f88175b71c`). The card uid can be found in `repository.go`.

To get access to the command you need to add the role `admin` to the permissions array on your user in the database.

# Changelog

A changelog starting from 11/11/2021 can be found [here](https://github.com/sindreslungaard/duel-masters/blob/master/CHANGELOG.md)
