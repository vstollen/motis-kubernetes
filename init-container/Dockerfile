FROM golang:1.18-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /init-script

FROM alpine

WORKDIR /

COPY --from=build /init-script /init-script
COPY schedules.txt test.txt

ENTRYPOINT ["/init-script"]
