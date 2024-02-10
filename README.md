# gpuctl - A GPU control room

## Development

### Control

TODO: an easier way to dev locally

#### Building & Testing

Run the `Makefile` in the root directory of the repository. This will build
binaries in the root dir. See configuration section.

Then, you can run tests with `go test ./...`. This requires access to a postgres
installation and role with permission to create tables in a database.

The tests are configured to connect to a database called `postgres` owned by the
role `postgres` at `localhost` (ie. `postgres://postgres@localhost/postgres`).
It will *ERASE* the contents of this database as part of test cleanup, so *DO
NOT* use it for deployment.

You can override this by setting the `TEST_URL` environment variable.

Alternatively, exclude the Postgres tests with `go test -short ./...`

#### Configuration

To configure how the binary runs locally, copy `deploy/control.prod.toml` and
`deploy/satellite.example.toml` to `control.toml` and `satellite.toml` and
modify as needed. Some important options include:

- `postgres` and `inmemory` in control.toml: control which database interface is
  used. Set one to `true` and the other to `false`.
- username and password for onboarding new machines
- `API_URL` in `frontend/src/App.tsx`. Needs to match `WAPort` in `control.toml`
- `protocol` & `hostname` & `port` in `satellite.toml` need to match `GSPort`
  in `control.toml`

### Frontend

This project was bootstrapped with [Create React App](https://github.com/facebook/create-react-app).

#### Available Scripts

In the project directory, you can run:

* `npm start`

Runs the app in the development mode.\
Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.\
You will also see any lint errors in the console.

* `npm test`

Launches the test runner in the interactive watch mode.\
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

* `npm run build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

* `npm run eject`

**Note: this is a one-way operation. Once you `eject`, you can’t go back!**

If you aren’t satisfied with the build tool and configuration choices, you can `eject` at any time. This command will remove the single build dependency from your project.

Instead, it will copy all the configuration files and the transitive dependencies (webpack, Babel, ESLint, etc) right into your project so you have full control over them. All of the commands except `eject` will still work, but they will point to the copied scripts so you can tweak them. At this point you’re on your own.

You don’t have to ever use `eject`. The curated feature set is suitable for small and middle deployments, and you shouldn’t feel obligated to use this feature. However we understand that this tool wouldn’t be useful if you couldn’t customize it when you are ready for it.

#### Learn More

You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).

To learn React, check out the [React documentation](https://reactjs.org/).

## Deploying gpuctl

### Using Docker Compose

#### Variables & Secrets

Copy `deploy/....` to `deploy/...` and modify it's contents with:

- `HETZNER_DNS_API_TOKEN` - your dns.hetzner.com API key (for HTTPS certificates)
- `GPU_SSH_KEY` - a base-64 encoded SSH key for logging into monitored machines with

#### Volumes

Note: most Docker commands need to be run as root/with sudo.

This method uses two volumes, one persisting the database and the other storing HTTPS certificates managed by Caddy. To create them, run:

```
docker volume create caddy_data
docker volume create postgres_data
```

#### Firewalls

No filewall changes should be needed, as Docker makes them itself. It will open access to port 80 for TCP, and port 443 for TCP & UDP.

#### Running/Stopping

From the `deploy/` directory, run:

```
docker compose build
docker compose up -d
```

Then, you can bring down a deployment with:

```
docker compose down
```

#### Postgres Debugging

It can be useful to get access to the database to see rows & make changes. From the `deploy/` directory, run:

```
docker compose exec -it postgres psql -U postgres postgres
```

## Building Docker images freestanding

Due to how Docker context works, you need to do this from the top level directory.

### Docker

```
docker build -f ./deploy/control.Dockerfile .
docker build -f ./deploy/frontend.Dockerfile .
```

### Podman

```console
alona@Ashtabula:~/dev/gpuctl$ podman build -f ./deploy/control.Dockerfile .
alona@Ashtabula:~/dev/gpuctl$ podman build -f ./deploy/frontend.Dockerfile .
```

## npm ERR! EMFILE: too many open files

You may see an error like:

```
[1/3] STEP 4/6: RUN npm install
npm ERR! code EMFILE
npm ERR! syscall open
npm ERR! path /root/.npm/_cacache/index-v5/7d/8e/9676576fe239de89dec5e769bcbad0def29c3e4fd33b2caebf5c78716e7e
npm ERR! errno -24
npm ERR! EMFILE: too many open files, open '/root/.npm/_cacache/index-v5/7d/8e/9676576fe239de89dec5e769bcbad0def29c3e4fd33b2caebf5c78716e7e'

npm ERR! A complete log of this run can be found in: /root/.npm/_logs/2024-02-08T11_54_42_543Z-debug-0.log
Error: building at STEP "RUN npm install": while running runtime: exit status 232
```

This can be solved with the `--ulimit` flag:

```console
alona@Ashtabula:~/dev/gpuctl$ podman build --ulimit=4096:4096 -f ./deploy/frontend.Dockerfile .
```
