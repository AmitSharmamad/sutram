build:
	docker build -t sutram-go .

di:
	wire src/di/wire.go

web:
	clear
	make di
	go build -tags wireinject -o web src/web/main.go
	./web

run:
	docker run -it -p 8080:8080 sutram-go

t:
	go test ./...

c:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
