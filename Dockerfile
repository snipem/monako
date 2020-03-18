# run: make image
FROM golang:1.13

WORKDIR /build
ADD . /build

RUN make init build

FROM asciidoctor/docker-asciidoctor:1.1.0

RUN apk add --no-cache \
    git

COPY --from=0 build/monako /usr/bin/monako
RUN chmod +x /usr/bin/monako

WORKDIR /docs