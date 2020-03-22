
#run: make run
SHELL := /bin/bash

clean:
	rm -r tmp/theme || true
	rm ./monako || true

deps:
	go mod download
	go get -u github.com/go-bindata/go-bindata/...

optional_deps:
	gem install asciidoctor-diagram

init: deps theme

build: clean
	go build -o ./monako github.com/snipem/monako/cmd/monako

theme: clean
	mkdir -p tmp/
	curl -o tmp/theme.zip --location https://github.com/snipem/monako-book/archive/v6s.zip
	${GOPATH}/bin/go-bindata -pkg theme -o internal/theme/bindata.go tmp/...

test:
	go test -v -coverprofile=c.out.tmp ./...
	cat c.out.tmp | grep -v "/bindata.go" > c.out
	rm c.out.tmp

run_prd: build
	touch secrets.env && source secrets.env && \
		./monako -config ~/work/mopro/architecture/documentation/conf/config.prod.yaml \
			-menu-config ~/work/mopro/architecture/documentation/conf/index.prod.md \
			-hugo-config ~/work/mopro/architecture/documentation/conf/config.prod.toml
	$(MAKE) serve

run: build
	./monako -config configs/config.monako.yaml \
		-menu-config configs/config.menu.md \
		-hugo-config configs/config.hugo.toml
		$(MAKE) serve

serve:
	echo "Serving under http://localhost:8000"
	/usr/bin/env python3 -m http.server 8000 --directory compose/public

image:
	docker build -t monako/monako:0.0.1 .

run_image:
	docker run -v ${PWD}:/docs monako/monako:0.0.1 monako