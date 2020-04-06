
#run: make run
SHELL := /bin/bash
.PHONY: compose test

clean:
	rm -r tmp/theme || true
	rm ./monako || true

deps:
	go mod download
	go get github.com/go-bindata/go-bindata/...

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
	# -coverpkg=./...  also calculates the whole coverage, for example code that was involded by the main test
	cat coverage.out.tmp | grep -v "/bindata.go" > coverage.out
	rm coverage.out.tmp

test_deps:
	go get golang.org/x/tools/cmd/cover

test_local_clone:
	git clone https://github.com/snipem/monako-test.git /tmp/testdata/monako-test

# Use this for local tests, uses the locally cloned test data from test_data step
test_local:
	MONAKO_TEST_REPO="/tmp/testdata/monako-test" $(MAKE) test

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
	./monako \
		-fail-on-error \
		-config configs/config.monako.yaml \
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
	brew install golangci/tap/golangci-lint
	git config --local core.hooksPath .githooks/
