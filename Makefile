
#run: make run
SHELL := /bin/bash

clean:
	rm ./monako || true

build: clean
	go build .

test:
	go test -v

run: build
	touch secrets.env && source secrets.env && \
		./monako -config ~/work/mopro/architecture/documentation/conf/config.prod.yaml \
			-menu-config ~/work/mopro/architecture/documentation/conf/index.prod.md \
			-hugo-config ~/work/mopro/architecture/documentation/conf/config.prod.toml

run_prd: build
	touch secrets.env && source secrets.env && ./monako -config config.prod.yaml -menu-config index.prod.md -hugo-config config.prod.toml
	hugo --source compose serve