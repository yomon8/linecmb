BIN      := fdlinemerge
OSARCH   := "darwin/amd64 linux/amd64"

all: build

test: deps build
	./test_$(BIN).sh

deps:
	go get -d -v -t ./...
	go get github.com/golang/lint/golint
	go get github.com/mitchellh/gox

lint: deps
	go vet ./...
	golint -set_exit_status ./...

crossbuild:
	rm -fR ./pkg && mkdir ./pkg ;\
		gox \
		-osarch $(OSARCH) \
		-output "./pkg/{{.OS}}_{{.Arch}}/{{.Dir}}" \
		./cmd/...

build:
	go build -o $(BIN) ./cmd/...

linuxbuild:
	GOOS=linux GOARCH=amd64 make build

clean:
	go clean
