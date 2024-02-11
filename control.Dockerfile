FROM golang:1.21

WORKDIR /gpuctl
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY cmd ./cmd
COPY internal ./internal
COPY Makefile Makefile

RUN make control

COPY control.toml.default ./control.toml

ENTRYPOINT ["./control"]
