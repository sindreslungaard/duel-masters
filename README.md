# duel-masters

duel-masters is a simulator for the [Duel Masters Trading Card Game](<https://duelmasters.fandom.com/wiki/Duel_Masters_(Card_Game)>) for play in the browser.

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

## Contribute

Thanks for your interest in contributing! There are many ways to contribute to this project. [Get started here](https://github.com/sindreslungaard/duel-masters/blob/master/CONTRIBUTING.md)
