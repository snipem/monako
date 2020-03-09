SHELL := /bin/bash

clean:
	rm ./monako || true

run: clean
	touch secrets.env && source secrets.env && go build && ./monako -config config.prod.yaml
	hugo --source compose serve