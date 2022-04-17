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


# Setting up for local development (Frontend and Backend)

1. Fork the `duel-masters` repo on GitHub.
2. Clone your fork locally:
    
```
git clone https://github.com/sindreslungaard/duel-masters.git
```

3. Set up [MongoDB locally](https://www.mongodb.com/try/download/community) or use a [cloud provider](https://www.mongodb.com/atlas/database).
4. Set up environment variables from the `.env.default` file (if you use Vscode it will look a `.env` file and set the variables for you. You have to create this file yourself based on the `.env.default`)

Environment variables or `.env` file example:
```
port=80
mongo_uri=mongodb://127.0.0.1:27017
mongo_name='duel-masters'
restart_after=
```


5. Navigate to the `webapp` directory and run `npm install`. Then run either `npm run build` or `npm run watch` to build or watch the files.

6. Run the application. If you're using Vscode simply hit F5 or `Run -> Start Debugging`. To run manually use `go run cmd/duel-masters/main.go`


7. Go to `http://localhost` and create a user as well as a deck. To set the deck as a standard deck, find it in MongoDB and change the `standard` field to `true`.

# Setting up for local development (Frontend only)

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

# Changelog
A changelog starting from 11/11/2021 can be found [here](https://github.com/sindreslungaard/duel-masters/blob/master/CHANGELOG.md)