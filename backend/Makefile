.PHONY: brun

build:
	docker build -t sutram .

run:
	docker run -it -v ./backend:/app/backend -v ./db:/app/db -p 5678:5678 sutram

brun: build run
