build:
	docker build -t sutram-go .

di:
	wire src/di/wire.go

d:
	air

b:
	@wire src/web/wire.go
	@echo "\n"
	@go build -o web src/web/wire_gen.go src/web/main.go
	@echo "\n"
	@./web

r:
	clear
	make b

run:
	docker run -it -p 8080:8080 sutram-go

t:
	go test ./...

c:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
