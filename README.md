# gpuctl - A GPU control room

## Deploying

Needs a Postgres database. URL is passed in in the `control.toml` file. For
_reasons_, don't use `postgres://postgres@localhost/postgres`.

## Running tests

Running tests with `go test ./...` requires access to a postgres installation
and role with permission to create tables in a database.

The tests are configured to connect to a database called `postgres` owned by the
role `postgres` at `localhost` (ie. `postgres://postgres@localhost/postgres`).
It will *ERASE* the contents of this database as part of test cleanup, so *DO
NOT* use it for deployment.

You can override this by setting the `TEST_URL` environment variable.

## Deploying

Deploys using Docker and Docker Compose. You will need external volumes for
storing persistant database contents and certificates:

```
docker volume create caddy_data
docker volume create postgres_data

docker compose up --build -d
```
