
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

run_prd: build
	touch secrets.env && source secrets.env && \
		./monako -config ~/work/mopro/architecture/documentation/conf/config.prod.yaml \
			-menu-config ~/work/mopro/architecture/documentation/conf/index.prod.md \
			-hugo-config ~/work/mopro/architecture/documentation/conf/config.prod.toml

run_prd: build
	./monako
	hugo --source compose serve