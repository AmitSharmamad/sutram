.PHONY: build start stop restart logs quick-run

build:
	docker-compose -f docker-compose.yml build

start:
	docker-compose -f docker-compose.yml up -d

stop:
	docker-compose -f docker-compose.yml down

restart: stop start

logs:
	docker-compose -f docker-compose.yml logs -f

quick-run:
	./run-gateway.sh
