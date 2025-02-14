# .PHONY: build && run
#
# build:
# 	docker-compose build
#
# run:
# 	docker-compose up -d
#
# down:
# 	docker-compose down
#
# deploy:
# 	docker-compose up --build -d


### prod

build:
	go mod download
	wire cmd/web/wire.go
	go build -o s-web cmd/web/wire_gen.go cmd/web/main.go

run:
	wire cmd/web/wire.go
	go build -o s-web cmd/web/wire_gen.go cmd/web/main.go
	./s-web

debug:
	lsof -ti tcp:9090 | xargs kill -9
	# wire src/web/wire.go
	go build -gcflags="all=-N -l" -o s-web src/web/wire_gen.go src/web/main.go 
	air

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
