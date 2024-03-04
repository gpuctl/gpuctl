# gpuctl - A GPU control room

## Deploying

The recommended way to deploy is using Docker and Docker Compose.

1. Clone the repo

   ```
   git clone https://github.com/gpuctl/gpuctl.git
   ```

2. Configure secrets

   ```
   cd deploy
   cp .env.example .env
   ```

   Then edit `.env` to contain the secrets. See the comments in `.env` for details.

3. Optionally edit `./deploy/control.prod.toml`

4. Create docker volumes for persistent data:

   ```
   docker volume create caddy_data
   docker volume create postgres_data
   ```

5. Build & Start Docker images:

   ```
   docker compose -f ./deploy/compose.yaml up --build --detach
   ```

To update the running version:

```
git pull
docker compose -f ./deploy/compose.yaml down
docker compose -f ./deploy/compose.yaml up --build --detach
```

(TODO: How to install/upgrade satellites)

## Developing

Building the `control` binary requires having build the `satellite` (so it can be embedded). Developers are encouraged to use `make control` and `make satellite` (over `go build`) for this purpose.

Additionally, `go run` won't work, due to loading the config files relative
the the path of the executable (and not the cwd).

### Configuration

To configure how the binary runs locally, copy `deploy/control.prod.toml` and
`deploy/satellite.example.toml` to `control.toml` and `satellite.toml` and
modify as needed. Some important options include:

- `postgres` and `inmemory` in `control.toml`: control which database interface is
  used. Set one to `true` and the other to `false`.
- username and password for onboarding new machines
- `API_URL` in `frontend/src/App.tsx`. Needs to match `WAPort` in `control.toml`
- `protocol` & `hostname` & `port` in `satellite.toml` need to match `GSPort`
  in `control.toml`

### Running go unit tests

Most of the tests have no service dependencies, and can be run with
`go test -short ./...`. (Note that `make all` is needed first to ensure it compilles)

Running all the tests with requires access to a postgres installation and role
with permission to create tables in a database. The tests are configured to
connect to a database called `postgres` owned by the role `postgres` at
`localhost` (ie. `postgres://postgres@localhost/postgres`). It will _ERASE_ the
contents of this database as part of test cleanup, so _DO NOT_ use it for
deployment. You can override this by setting the `TEST_URL` environment
variable.

When writing new tests, try not to use the postgres if at all possible.
If you do, you must add

```go
if testing.Short() {
    t.Skip("not connecting to postgres in short tests")
}
```

to the top of your test functions. CI enforces this.

### Frontend

See [frontend/README.md](frontend/README.md)
