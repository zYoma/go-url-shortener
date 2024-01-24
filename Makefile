
run:
	go run cmd/shortener/main.go 

test:
	go test ./...

mock:
	cd internal && mockery --all && cd -