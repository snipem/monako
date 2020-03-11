
#run: make run
SHELL := /bin/bash

clean:
	rm -r tmp/theme || true
	rm ./monako || true

deps:
	go get -u github.com/go-bindata/go-bindata/...

init: deps theme

build: clean
	go build .

theme: clean
	mkdir -p tmp/
	wget https://github.com/alex-shpak/hugo-book/archive/v6.zip -O tmp/theme.zip
	go-bindata tmp/...

test:
	go test -v

run: build
	./monako -trace -config config.yaml -menu-config index.md -hugo-config config.toml
	hugo --source compose serve

run_prd: build
	touch secrets.env && source secrets.env && ./monako -config config.prod.yaml -menu-config index.prod.md -hugo-config config.prod.toml
	hugo --source compose serve