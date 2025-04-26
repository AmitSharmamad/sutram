ch:
	docker-compose up -d sutram-clickhouse
pg:
	docker-compose up -d sutram-postgres

all:
	docker-compose up --build

down:
	docker-compose down

clean:
	rm -rf db/ch
	rm -rf db/pg

start:
	docker-compose down
	rm -rf db/ch
	rm -rf db/pg
	docker-compose up --build
