SHELL := /bin/bash

clean:
	rm ./monako || true

build: clean
	go build .

run: build
	touch secrets.env && source secrets.env && ./monako -config config.prod.yaml -menu-config index.prod.md -hugo-config config.prod.toml
	hugo --source compose serve