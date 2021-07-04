# Contributing to duel-masters
First off – __thank you__ for considering contributing to `duel-masters`!

Please take a moment to review this document in order to make the contribution process easy and effective for everyone involved.

# Your first contribution

## Setting up the code for local development
This is how to set up `duel-masters` for local development:

1. Fork the `duel-masters` repo on GitHub.
2. Clone your fork locally:
    
```
git clone https://github.com/sindreslungaard/duel-masters.git
```

3. Set up a MongoDB server instance. One way is to [run MongoDB on Docker](https://hub.docker.com/_/mongo)
4. Set up the runtime environment. Assuming you have VisualStudio Code, this is how you set up to run `duel-masters` locally:

```
cd duel-masters

cp .env.default .env
nano .env

# You should see:
# port=80
# mongo_uri=
# mongo_name=

# Set:
mongo_uri="<mongodb connection string uri>"
mongo_name="<mongodb database name>"

# Example:
mongo_uri="mongodb://username:password@192.168.99.100:27017/duel-masters?authSource=admin"
mongo_name="duel-masters"
```

Read more on MongoDB connection strings [here](https://docs.mongodb.com/guides/server/drivers/)

If you do not use vsCode, you will have to set `mongo_uri` and `mongo_name` as environment variables manually.

5. Build the webapp:

```
# From /duel-masters
cd webapp
npm install
npm run build
```

6. Launch `duel-masters` through the vsCode launch configuration (`/.vscode/launch.json`)

If you do not use vsCode, run:
```
go get -d -v ./...
go install -v ./...
```

7. Access on `localhost:80` and create a user to demo your changes:

Create and save a deck, then go into the database through the mongo shell, find the deck and update attribute `"standard"` to `true`. 

Standard decks are available to all users. Without a standard deck you must create a new deck for each dummy user you wish to demo.