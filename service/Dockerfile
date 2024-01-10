FROM golang:1.21-bullseye as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o saladbowl-service ./cmd/app/app.go


FROM debian:bullseye-slim

WORKDIR /app
COPY --from=build /app/saladbowl-service ./saladbowl-service

EXPOSE 8080

CMD [ "./saladbowl-service" ]