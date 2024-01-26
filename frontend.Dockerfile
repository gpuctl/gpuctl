FROM node:21.6 AS builder

WORKDIR /frontend

COPY frontend/package.json frontend/package-lock.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build

FROM caddy:2.7.6

COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /frontend/build /srv
