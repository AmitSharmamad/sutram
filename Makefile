.PHONY: build && run

build:
	docker-compose build

run:
	docker-compose up -d

deploy:
	docker-compose up --build -d
