clean:
	rm bin/*

server:
	@go build -o bin/server main.go
	@./bin/server