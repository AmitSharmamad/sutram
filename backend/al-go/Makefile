build:
	docker build -t sutram-go .

web:
	go run ./src/web/main.go

run:
	docker run -it -p 8080:8080 sutram-go

t:
	go test ./...

c:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
