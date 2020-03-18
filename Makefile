
#run: make run
SHELL := /bin/bash

clean:
	rm -r tmp/theme || true
	rm ./monako || true

deps:
	go mod download

init: deps theme

build: clean
	go build .

build_linux: clean
	mkdir -p builds/linux
	GOOS=linux GOARCH=386 go build -o builds/linux/monako .

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

image: build_linux
	docker build -t monako/monako:0.0.1 .

run_image:
	docker run -v ${PWD}:/docs monako/monako:0.0.1 monako