
#run: make run
SHELL := /bin/bash

clean:
	rm ./monako || true

build: clean
	go build .

run: build
	./monako -config config.yaml -menu-config index.md -hugo-config config.toml
	hugo --source compose serve

run_prd: build
	touch secrets.env && source secrets.env && ./monako -config config.prod.yaml -menu-config index.prod.md -hugo-config config.prod.toml
	hugo --source compose serve