
run:
	go run cmd/shortener/main.go -d postgresql://postgres:postgres@127.0.0.1:9432/shorther

test:
	go test ./...

bench:
	go test -bench=. -memprofile=profiles/result.pprof ./internal/handlers

pprof:
	go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof

mock:
	cd internal && mockery --all && cd -

lint:
	go run cmd/staticlint/main.go ./...