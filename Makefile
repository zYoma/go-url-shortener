
run:
	go run cmd/shortener/main.go -d postgresql://postgres:postgres@127.0.0.1:9432/shorther

test:
	go test ./...

mock:
	cd internal && mockery --all && cd -