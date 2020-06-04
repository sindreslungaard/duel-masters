FROM golang:1.14

WORKDIR /go/src/duel-masters
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 80

CMD ["duel-masters"]