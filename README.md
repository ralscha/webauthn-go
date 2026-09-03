# WebAuthn demo with Go

A passkey registration and authentication example with a Go API, PostgreSQL,
and an Angular client. It accompanies the
[original blog post](https://blog.rasc.ch/2024/09/webauthn_go.html).

## Run locally

Prerequisites: Go, Node.js/npm, Docker, and Docker Compose.

1. Start PostgreSQL:

   ```shell
   cd server
   docker compose up -d
   ```

2. Apply the database migration and start the API:

   ```shell
   go run ./cmd/migrate up
   go run ./cmd/api
   ```

3. In another terminal, install and start the client:

   ```shell
   cd client
   npm install
   npm start
   ```

Open <http://localhost:4200>. The development server proxies `/api/v1` to the
API at `127.0.0.1:8080`.

## Configuration

The API reads `server/app.yml`. Values can be overridden with environment
variables using the `WEBAUTHN_` prefix and underscores for nesting. For
example, `WEBAUTHN_DB_HOST=db:5432` overrides `db.host`.

For production, serve the client and proxy `/api/v1` to the Go API on the same
origin, then update `webauthn.rpId` and `webauthn.rpOrigins` to match that
origin. Secure session cookies are enabled by default outside the checked-in
development configuration.

## Checks

```shell
cd server
go test ./...
go vet ./...

cd ../client
npm run lint
npm run build
npm audit
```
