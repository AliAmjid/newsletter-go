FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
COPY . .
RUN go build -o newsletter ./app
COPY ./internal/mailer/templates /app/internal/mailer/templates

FROM alpine
RUN apk add --no-cache ca-certificates go
WORKDIR /app
COPY --from=build /app/newsletter .
COPY --from=build /app/internal/mailer/templates /app/internal/mailer/templates
COPY --from=build /go/bin/goose /usr/local/bin/goose
COPY db/migrations /app/db/migrations
COPY .env .env
EXPOSE 3000

CMD ["./newsletter"]
