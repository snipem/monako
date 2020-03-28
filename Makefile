
#run: make run
SHELL := /bin/bash
.PHONY: compose test

clean:
	rm -r tmp/theme || true
	rm ./monako || true

deps:
	go mod download
	go get -u github.com/go-bindata/go-bindata/...

optional_deps:
	gem install asciidoctor asciidoctor-diagram

init: deps theme

build: clean
	go build -o ./monako github.com/snipem/monako/cmd/monako

theme: clean
	mkdir -p tmp/
	curl -o tmp/theme.zip --location https://github.com/snipem/monako-book/archive/master.zip
	${GOPATH}/bin/go-bindata -pkg theme -o internal/theme/bindata.go tmp/...

secrets:
	touch config/secrets.env && source config/secrets.env	

test:
	go test -covermode=count -coverprofile=coverage.out.tmp ./...
	cat coverage.out.tmp | grep -v "/bindata.go" > coverage.out
	rm coverage.out.tmp

coverage: test
	${GOPATH}/bin/goveralls -coverprofile=coverage.out -service=travis-ci -repotoken ${COVERALLS_TOKEN}

test_deps:
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls

run_prd: build secrets
		./monako -config ~/work/mopro/architecture/documentation/conf/config.prod.yaml \
			-menu-config ~/work/mopro/architecture/documentation/conf/menu.prod.md \
			-base-url http://localhost:8000

	$(MAKE) serve

run: build compose serve

run_local: build 
	# Runs locally, clons this git repo to use test data
	./monako -config test/config.local.yaml -menu-config test/config.menu.local.md
	$(MAKE) serve

compose:
	./monako -config configs/config.monako.yaml \
		-menu-config configs/config.menu.md

serve:
	echo "Serving under http://localhost:8000"
	/usr/bin/env python3 -m http.server 8000 --directory compose/public

image:
	docker build -t monako/monako:0.0.1 .

run_image:
	docker run -v ${PWD}:/docs monako/monako:0.0.1 monako

hooks:
	# setup git hooks
	git config --local core.hooksPath .githooks/
