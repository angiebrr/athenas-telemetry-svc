FROM golang:1.26.4-alpine

WORKDIR /app

COPY ./athenas-telemetry-svc/go.mod ./

RUN go mod download

COPY ./athenas-telemetry-svc .

RUN go build -o svr ./cmd/athenas/*.go

CMD [ "./svr" ]