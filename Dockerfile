FROM golang:1.18-alpine

WORKDIR /go/src/duel-masters
COPY . .

RUN apk add --update git nodejs npm
RUN cd ./webapp && npm install && npm run build
RUN cd ..

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 80

CMD ["duel-masters"]