FROM golang:1.20-bullseye as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o saladbowl


FROM debian:bullseye-slim

WORKDIR /app
COPY --from=build /app/saladbowl ./saladbowl

CMD [ "./saladbowl" ]