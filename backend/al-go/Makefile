### prod

build:
	go mod download
	wire src/web/wire.go
	go build -o s-web src/web/wire_gen.go src/web/main.go

### local

dbuild:
	docker build -t sutram-go .

drun:
	docker run -it -p 9999:9999 sutram-go

### testing

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
