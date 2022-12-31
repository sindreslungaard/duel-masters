FROM golang:1.18-alpine

WORKDIR /go/src/duel-masters
COPY . .

RUN apk add --update git nodejs=16.17.1-r0 npm
RUN cd ./webapp && npm install && npm run build
RUN cd ..

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 80

CMD ["duel-masters"]