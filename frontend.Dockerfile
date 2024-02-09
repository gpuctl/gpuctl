FROM node:21.6 AS builder

WORKDIR /frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build

# Build our own Caddy with a DNS provider module
FROM caddy:2.7.6-builder AS caddy-modules

RUN xcaddy build --with github.com/caddy-dns/hetzner

FROM caddy:2.7.6

COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /frontend/build /srv
COPY --from=caddy-modules /usr/bin/caddy /usr/bin/caddy
