FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o newsletter ./app

FROM alpine
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app/newsletter .
COPY .env .env
EXPOSE 3000
CMD ["./newsletter"]