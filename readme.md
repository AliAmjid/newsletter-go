# Newsletter application

### Diagram

![img.png](docs/img.png)
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
MAILGUN_DOMAIN=<mg-domain>
MAILGUN_API_KEY=<mg-key>
MAILGUN_FROM_EMAIL=<from-email>
```

## Database migrations
DB migration files are stored in `db/migrations` folder. To apply migrations call following command:

goose postgres "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable" up