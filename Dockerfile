FROM golang:1.22-alpine AS build

WORKDIR /tmp/duel-masters
COPY ./sim .

RUN go mod download \
    && go build -o app ./cmd/duel-masters/main.go \
    && mkdir ./build \
    && mv ./app ./build/app \
    && mv ./DuelMastersCards.json ./build/DuelMastersCards.json

FROM alpine:3

WORKDIR /go/src/duel-masters
COPY --from=build /tmp/duel-masters/build .
EXPOSE 80

CMD ["./app"]