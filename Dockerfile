FROM golang:1.14-alpine

WORKDIR /go/src/duel-masters
COPY . .

RUN apk add --update nodejs npm
RUN cd ./webapp && npm install && npm run build && rm -rf ./node_modules && && rm -rf ./public
RUN cd ..

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 80

CMD ["duel-masters"]