FROM golang:1.21

WORKDIR /gpuctl
COPY cmd ./cmd
COPY internal ./internal
COPY go.mod go.sum ./

RUN go build ./cmd/control

ENTRYPOINT ["./control", "-postgres"]
