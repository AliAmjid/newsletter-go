# Implementation

This project is a Go service that exposes a REST API for managing newsletters, posts and subscribers. The codebase follows a clean architecture style with clear separation between domain entities, use cases, repositories and HTTP delivery.

## Libraries

The service relies on a small set of external libraries:

- **[go-chi/chi](https://github.com/go-chi/chi)** – lightweight HTTP router used to build the REST endpoints.
- **[joho/godotenv](https://github.com/joho/godotenv)** – loads configuration from `.env` files during local development.
- **[go-playground/validator](https://github.com/go-playground/validator)** – request payload validation.
- **[google/uuid](https://github.com/google/uuid)** – generation of unique identifiers.
- **[lib/pq](https://github.com/lib/pq)** – PostgreSQL driver.
- **[mailgun/mailgun-go](https://github.com/mailgun/mailgun-go)** – sending transactional emails.
- **[permitio/permit-golang](https://github.com/permitio/permit-golang)** – authorisation checks.
- **[firebase.google.com/go](https://firebase.google.com/docs/reference/admin/go)** – Firebase authentication integration.
- **[pressly/goose](https://github.com/pressly/goose)** – database migrations (run separately).

Additional Google Cloud and OpenTelemetry packages are included for metrics and Firebase client support.

## System design

The code is organised into the following layers:

1. **Domain (`/domain`)** – defines core entities and repository interfaces.
2. **Use cases (`/internal/usecase`)** – business logic grouped by feature (posts, newsletters, authentication, subscribers and users).
3. **Repository implementations (`/internal/repository/postgres`)** – PostgreSQL based storage for domain interfaces.
4. **HTTP delivery (`/internal/delivery/http`)** – handlers, middleware and router setup built with chi.
5. **Dependency injection (`/internal/di`)** – wires services, repositories and external clients (database, Mailgun, Firebase, Permit).
6. **Mailer (`/internal/mailer`)** – renders HTML templates and sends emails via Mailgun.
7. **Database migrations (`/db/migrations`)** – SQL scripts managed by goose.

The `app/main.go` file creates the dependency container, sets up the router and starts the HTTP server. Environment variables are loaded with `godotenv` when running locally.
