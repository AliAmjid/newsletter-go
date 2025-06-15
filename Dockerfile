FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o newsletter ./app
COPY ./internal/mailer/templates /app/internal/mailer/templates

FROM alpine
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app/newsletter .
COPY --from=build /app/internal/mailer/templates /app/internal/mailer/templates
EXPOSE 3000

CMD ["./newsletter"]
