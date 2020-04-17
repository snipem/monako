
# run: make theme && head internal/theme/bindata.go     
SHELL := /bin/bash
.PHONY: compose test

clean:
	rm ./monako || true

deps:
	go mod download

	# For theme
	git submodule update
	go get -u github.com/go-bindata/go-bindata/...

optional_deps:
	gem install asciidoctor asciidoctor-diagram

init: deps theme

build: clean
	go build -o ./monako github.com/snipem/monako/cmd/monako

theme: clean
	go generate cmd/monako/main.go

secrets:
	touch configs/secrets.env && source configs/secrets.env

test: clean_test
	go test -covermode=count -coverprofile=coverage.out.tmp ./...
	cat coverage.out.tmp | grep -v "/bindata.go" > coverage.out
	rm coverage.out.tmp

clean_test:
	rm -r tmp/testdata/ || true

test_deps:
	go get golang.org/x/tools/cmd/cover

clones_for_local_testing:
	git clone https://github.com/snipem/monako-test.git ${HOME}/temp/monako-testrepos/monako-test
	git clone https://github.com/gohugoio/hugo.git ${HOME}/temp/monako-testrepos/hugo

# Use this for local tests, uses the locally cloned test data from test_data step
test_local:
	MONAKO_TEST_REPO="${HOME}/temp/monako-testrepos/monako-test" $(MAKE) test

benchmark:
	go test -v ./pkg/compose/ -run=BenchmarkHugeRepositories -bench=Benchmark. -benchtime=10s

run_prd: build secrets
		env | grep USER
		./monako -config ~/work/mopro/architecture/documentation/conf/config.prod.yaml \
			-menu-config ~/work/mopro/architecture/documentation/conf/menu.prod.md \
			-base-url http://localhost:8000

	$(MAKE) serve

run: build compose serve

run_local: build 
	# Runs locally, clons this git repo to use test data
	./monako -config test/config.local.yaml -menu-config test/config.menu.local.md
	$(MAKE) serve

trace:
	go test -trace=tmp/trace.out ./cmd/monako
	go tool trace ./cmd/monako/ tmp/trace.out

compose:
	./monako \
		-fail-on-error \
		-config configs/config.monako.yaml \
		-menu-config configs/config.menu.md

serve:
	echo "Serving under http://localhost:8000"
	/usr/bin/env python3 -m http.server 8000 --directory compose/public

# setup git hooks
hooks:
	which golangci-lint || brew install golangci/tap/golangci-lint
	git config --local core.hooksPath githooks/
