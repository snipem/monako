SHELL := /bin/bash

run:
	touch secrets.env && source secrets.env && go run main.go config.go
	hugo --source compose serve