.PHONY: help

help:
	cat Makefile

dev:
	cp ./infra/dev/.env .env
	cp ./infra/dev/Dockerfile Dockerfile
	docker build -t sutram .
	make devr

devr:
	docker run -it -v ./backend:/app/backend -v ./db:/app/db -p 5678:5678 sutram 

prod:
	cp ./infra/prod/.env .env
	cp ./infra/prod/Dockerfile Dockerfile
	docker build -t sutram .
	docker run -it -p 5678:5678 sutram
