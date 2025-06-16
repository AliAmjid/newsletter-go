# Newsletter application

### Diagram

![img.png](https://uml.planttext.com/plantuml/png/NP9FJm8n4CNl_HIJFS41lNCmWg1WK14m7iH3jpjP9fGMErqC6h-xxRB0BhVclT_apP_UifVE5nijF7cDy8pVhD5xT9q72psdLkHH-SAX459vwo0PcWpU65q2aik742uWqkIXFXaj5bCOeAOTkZrCvBu2Uzij_0g0ZrxXLp2I9jFeFgDmZcp8yo9bvIYzaAUM_LduQsI7PfHahALM2jUYA-aokYxNekjo6NqIGdGclKnZO1AjtE7yTR8qhInjb-63lX1AsoA3v9uSHd9fOWzeF2dfoPgIEvhHEXPCqMt8Nv4zL5X7F-U2Wpb-ES9FadHUcAKPRi8BvkXsfB14Aqk8U2ZeT6xAGtIXF9F3hmBPCxHcm_a2UjDnMqGxOgVDTO7CpguHqeB983DecvCI9o3IzH52nKvg2il1QwRSmEFxeaaV-m-aMg5QmpOAp64-RfA3Vc3kPcy3i84fDt11L0C6Z35yJ8mRu7y0)
The architecture follows go-clean-arch principles:
![diagram.png](https://raw.githubusercontent.com/bxcodec/go-clean-arch/master/clean-arch.png)
```
@startuml
title Newsletter-Go Architecture

actor "End User" as User

rectangle "HTTP Delivery Layer" as App {
  [Auth Handler]
  [Newsletter Handler]
  [Subscriber Handler]
  [Post Handler]
}

rectangle "Usecase Layer" as Usecases {
  [Auth Usecase]
  [Newsletter Usecase]
  [Subscriber Usecase]
  [Post Usecase]
}

database "PostgreSQL\n(db)" as DB

cloud "Firebase\nAuthentication" as FirebaseAuth
cloud "Permit.io\nAuthorization" as PermitIO
cloud "Mailgun\nEmail Service" as Mailgun

User --> App : HTTP requests (REST API)
App --> Usecases : invoke business logic
Usecases --> FirebaseAuth : validate/sign JWT
Usecases --> PermitIO : check permissions
Usecases --> DB : CRUD operations
Usecases --> Mailgun : send emails
```

## Project Structure

The source code follows a simplified version of the [go-clean-arch](https://github.com/bxcodec/go-clean-arch) layout:

```
app/            // main program entry
domain/         // business entities and repository interfaces
internal/
  db/           // database setup
  delivery/http // HTTP handlers and router
  repository/   // data access implementations
  usecase/      // business services
```

### Environment Variables

```
POSTGRES_CONNECTION_STRING=<db-url>
PERMIT_API_KEY=<permit-key>
FIREBASE_CREDENTIALS=<path-to-service-account>
FIREBASE_API_KEY=<web-api-key>
APP_DOMAIN=<app-domain>
MAILGUN_DOMAIN=<mg-domain>
MAILGUN_API_KEY=<mg-key>
MAILGUN_FROM_EMAIL=<from-email>
```

### Environment Management with dotenvx

This project uses [dotenvx](https://dotenvx.com/) for secure environment variable management. dotenvx encrypts our environment variables and allows us to safely store them in version control.

#### Key Features
- **Encrypted Environment Variables**: All sensitive data is encrypted in the `.env.vault` file
- **Environment Separation**: Support for different environments (development, production, etc.)
- **Secure Deployment**: Private keys are never committed to the repository

#### Usage

1. **Setup**: Create your `.env` file with your actual environment variables

2. **Encrypt**: Generate encrypted `.env.vault` and key files:
   ```bash
   npx dotenvx encrypt
   ```

3. **Run Locally**: Use dotenvx to run the application:
   ```bash
   npx dotenvx run -- go run app/main.go
   ```

4. **Deploy**: In production, set the `DOTENV_KEY` environment variable to the private key from `.env.keys`

5. **Rotate Keys**: Update encryption keys:
   ```bash
   npx dotenvx rotate
   ```

For more information, visit [dotenvx documentation](https://dotenvx.com/docs/).

## Database migrations
DB migration files are stored in `db/migrations` folder. To apply migrations call the following command:

```bash
goose postgres "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable" up

## Documentation of implementation
[Go to documentation](./docs/implementation/readme.md)
```

To create a new migration:

```bash
goose create add_new_table sql
```
