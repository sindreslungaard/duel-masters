FROM golang:1.22-alpine as build

WORKDIR /tmp/duel-masters
ENV NODE_OPTIONS=--openssl-legacy-provider
RUN apk add --update git nodejs npm
COPY . .
RUN cd ./webapp \
    && npm install \
    && npm run build \
    && cd ..
RUN go get -d -v ./... \
    && go build -o app ./cmd/duel-masters/main.go \
    && mkdir ./build \
    && mv ./app ./build/app \
    && mv ./DuelMastersCards.json ./build/DuelMastersCards.json \
    && mkdir ./build/webapp \
    && mv ./webapp/dist ./build/webapp


FROM golang:1.22-alpine

WORKDIR /go/src/duel-masters
COPY --from=build /tmp/duel-masters/build .
EXPOSE 80

CMD ["./app"]