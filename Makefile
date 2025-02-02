.PHONY: build && run

build:
	docker-compose build

run:
	docker-compose up -d

down:
	docker-compose down

deploy:
	docker-compose up --build -d
