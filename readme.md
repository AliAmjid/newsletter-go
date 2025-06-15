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
