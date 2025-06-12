ch:
	docker-compose up -d sutram-clickhouse
pg:
	docker-compose up -d sutram-postgres

all:
	make down
	make clean
	make up

up:
	docker-compose up -d --build

down:
	docker-compose down

clean:
	rm -rf db/ch
	rm -rf db/pg

start:
	docker-compose down -v
	docker-compose up --build -d
