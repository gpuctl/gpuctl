FROM docker.io/golang:1.21

WORKDIR /gpuctl
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd ./cmd
COPY internal ./internal
COPY Makefile Makefile

RUN make control

COPY ./deploy/control.prod.toml ./control.toml

ENTRYPOINT ["./control"]
