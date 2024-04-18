BINARY_NAME=shortener
VERSION=0.0.1

build:
	go build -o $(BINARY_NAME) -ldflags "-X 'main.buildVersion=$(VERSION)' -X 'main.buildDate=$$(date)' -X 'main.buildCommit=$$(git rev-parse HEAD)'" cmd/shortener/main.go 

run: build
	./$(BINARY_NAME)  -s

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