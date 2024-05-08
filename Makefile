BINARY_NAME=shortener
VERSION=0.0.1

build:
	go build -o $(BINARY_NAME) -ldflags "-X 'main.buildVersion=$(VERSION)' -X 'main.buildDate=$$(date)' -X 'main.buildCommit=$$(git rev-parse HEAD)'" cmd/shortener/main.go 

run: build
	./$(BINARY_NAME) -d postgresql://postgres:postgres@127.0.0.1:9432/shorther

test:
	go test ./...

bench:
	go test -bench=. -memprofile=profiles/result.pprof ./internal/handlers

pprof:
	go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof

mock:
	cd internal && mockery --all && cd -

lint:
	go run cmd/staticlint/main.go ./... ; golangci-lint run ./...

proto:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=. \
	--grpc-gateway_opt=paths=source_relative \
	proto/shortener.proto
