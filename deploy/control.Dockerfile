FROM docker.io/golang:1.22-alpine as build

WORKDIR /gpuctl
RUN apk add make
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd ./cmd
COPY internal ./internal
COPY Makefile Makefile

RUN make control

FROM docker.io/alpine:latest

WORKDIR /gpuctl

COPY --from=build /gpuctl/control control
COPY ./deploy/control.prod.toml ./control.toml

ENTRYPOINT ["./control"]
