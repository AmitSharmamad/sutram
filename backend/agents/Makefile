.PHONY: brun

venv:
	python3 -m venv .venv
	chmod +x ./.venv/bin/activate
	source ./.venv/bin/activate

install:
	pip3 install -r requirements.txt

build:
	docker build -t sutram .

run:
	docker run -it -v .:/app/ -p 5678:5678 sutram

brun: build run
