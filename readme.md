# Newsletter application

## Table of Contents
- [Diagram](#diagram)
- [Project Structure](#project-structure)
  - [Environment Variables](#environment-variables)
  - [Environment Management with dotenvx](#environment-management-with-dotenvx)
- [Database](#database)
  - [Tables Diagram](#tables-diagram)
  - [Tables](#tables)
  - [Relationships](#relationships)
  - [Database Migrations](#database-migrations)
- [Used GO packages](#used-go-packages)

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

## Database
This project uses **PostgreSQL** as the primary relational database, and [`goose`](https://github.com/pressly/goose) migration tool for managing database schema changes and migrations.

Below you’ll find the list of migration tables and instructions on how to use goose with this project.


### Tables diagram
![Database schema](./docs/assets/database-schema.svg)

### Tables

#### user
| Name | Type | Settings | References |
| - | - | - | - |
| **id** | UUID | 🔑 PK, null | |
| **email** | TEXT | not null, unique | |
| **created_at** | TIMESTAMP | not null, default: now() | |
| **firebase_uid** | TEXT | null, unique | |

#### password_reset_tokens
| Name | Type | Settings | References |
| - | - | - | - |
| **token** | TEXT | 🔑 PK, null | |
| **user_id** | UUID | null | fk_password_reset_tokens_user_id_user | 
| **expires_at** | TIMESTAMP | not null | |
| **created_at** | TIMESTAMP | not null, default: now() | |

#### newsletter
| Name | Type | Settings | References |
| - | - | - | - |
| **id** | UUID | 🔑 PK, null | |
| **name** | TEXT | not null | |
| **description** | TEXT | null | |
| **owner_id** | UUID | not null | fk_newsletter_owner_id_user |
| **created_at** | TIMESTAMP | not null, default: now() | |

#### subscription
| Name | Type | Settings | References |
| - | - | - | - |
| **id** | UUID | 🔑 PK, null | |
| **newsletter_id** | UUID | not null | fk_subscription_newsletter_id_newsletter |
| **email** | TEXT | not null | |
| **token** | TEXT | not null, unique | |
| **confirmed_at** | TIMESTAMP | null | |
| **created_at** | TIMESTAMP | not null, default: now() | |

#### post
| Name | Type | Settings | References |
| - | - | - | - | 
| **id** | UUID | 🔑 PK, null | |
| **newsletter_id** | UUID | not null | fk_post_newsletter_id_newsletter |
| **title** | TEXT | not null | |
| **content** | TEXT | not null | |
| **published_at** | TIMESTAMP | null | |

#### post_delivery
| Name | Type | Settings | References |
| - | - | - | - |
| **id** | UUID | 🔑 PK, null | |
| **post_id** | UUID | not null | fk_post_delivery_post_id_post |
| **subscription_id** | UUID | not null | fk_post_delivery_subscription_id_subscription |
| **opened** | BOOLEAN | not null, default: false | |


### Relationships
- **newsletter to user**: many_to_one
- **password_reset_tokens to user**: many_to_one
- **post to newsletter**: many_to_one
- **subscription to newsletter**: many_to_one
- **post_delivery to post**: many_to_one
- **post_delivery to subscription**: many_to_one

### Database migrations
This project uses [`goose`](https://github.com/pressly/goose) for managing database migrations.

Before applying migrations, make sure to have goose tool instaled
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

DB migration files are stored in `db/migrations` folder. To apply migrations call the following command from root folder:

```bash
goose -dir db/migrations postgres POSTGRES_CONNECTION_STRING up
```

To create a new migration:

```bash
goose create add_new_table sql
```
This will generate a new pair of .sql files (up and down) in the db/migrations directory.
For more information, refer to the goose documentation.


## Used GO packages
- [joho/godotenv](https://github.com/joho/godotenv)  
  Used to load environment variables from a `.env` file into the application at runtime.  
  This helps manage configuration in a clean and secure way during local development.

- [go-chi/chi](https://github.com/go-chi/chi)  
  A lightweight, idiomatic, and composable router for building HTTP services in Go.  
  We use it to define our API routes and middleware in a modular and efficient way.

- [lib/pq](https://github.com/lib/pq) and `database/sql`  
  `lib/pq` is a pure Go Postgres driver used together with the standard `database/sql` package  
  to handle database connections, queries, and transactions with PostgreSQL.