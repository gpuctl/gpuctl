FROM docker.io/node:21.6 AS builder

WORKDIR /frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm install

COPY ./frontend/ ./

# This doens't run tsc, but that's run elsewhere in CI
RUN npm run build-nocheck

# Build our own Caddy with a DNS provider module
FROM docker.io/caddy:2.7.6-builder AS caddy-modules

RUN xcaddy build --with github.com/caddy-dns/hetzner

FROM docker.io/caddy:2.7.6

COPY deploy/Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /frontend/dist /srv
COPY --from=caddy-modules /usr/bin/caddy /usr/bin/caddy
