FROM golang:1.24-alpine AS build
LABEL authors="Melio"

WORKDIR "/app"

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -trimpath -tags netgo -ldflags="-s -w -extldflags '-static'" -o /app/app .

FROM alpine AS final

WORKDIR /app

COPY --from=build /app/app /app/app

COPY .env .

COPY ./templates ./templates

RUN chmod +x /app/app

ENTRYPOINT ["/app/app"]

