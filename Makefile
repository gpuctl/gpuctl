all: control satellite

# Go uses it's own cache, so we're fine to run these redund
.PHONY: internal/assets/satellite-amd64-linux
internal/assets/satellite-amd64-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o $@ ./cmd/satellite/

.PHONY: internal/assets/satellite-amd64-linux
control: internal/assets/satellite-amd64-linux
	go build -v -o $@ ./cmd/control

.PHONY: satellite
satellite:
	go build -v -o $@ ./cmd/satellite