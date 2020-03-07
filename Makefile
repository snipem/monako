SHELL := /bin/bash

run:
	source secrets.env && go run main.go config.go
	hugo --source compose serve