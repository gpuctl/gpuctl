FROM docker.io/node:21.6 AS builder

WORKDIR /frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm install

COPY ./frontend/ ./

# increase heap space to prevent oom crashes
ENV NODE_OPTIONS="--max_old_space_size=1024"
RUN npm run build

# Build our own Caddy with a DNS provider module
FROM docker.io/caddy:2.7.6-builder AS caddy-modules

RUN xcaddy build --with github.com/caddy-dns/hetzner

FROM docker.io/caddy:2.7.6

COPY deploy/Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /frontend/build /srv
COPY --from=caddy-modules /usr/bin/caddy /usr/bin/caddy
