
#run: make run
SHELL := /bin/bash

clean:
	rm -r tmp/theme || true
	rm ./monako || true

deps:
	go mod download
	go get -u github.com/go-bindata/go-bindata/...

init: deps theme

build: clean
	go build .

theme: clean
	mkdir -p tmp/
	curl -o tmp/theme.zip --location https://github.com/alex-shpak/hugo-book/archive/v6.zip
	${GOPATH}/bin/go-bindata tmp/...

test:
	go test -v

run_prd: build
	touch secrets.env && source secrets.env && \
		./monako -config ~/work/mopro/architecture/documentation/conf/config.prod.yaml \
			-menu-config ~/work/mopro/architecture/documentation/conf/index.prod.md \
			-hugo-config ~/work/mopro/architecture/documentation/conf/config.prod.toml
	hugo --source compose serve

run: build
	./monako
	hugo --source compose serve

image:
	docker build -t monako/monako:0.0.1 .

run_image:
	docker run -v ${PWD}:/docs monako/monako:0.0.1 monako