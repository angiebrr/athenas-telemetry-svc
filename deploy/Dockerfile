FROM golang:1.26.4-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o svr cmd/athenas/*.go

CMD [ "./svr" ]