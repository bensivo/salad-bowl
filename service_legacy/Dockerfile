FROM golang:1.20-bullseye as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o saladbowl-service


FROM debian:bullseye-slim

WORKDIR /app
COPY --from=build /app/saladbowl-service ./saladbowl-service

CMD [ "./saladbowl-service" ]