# run: make image
FROM golang:1.20

WORKDIR /build
ADD . /build

RUN make init

ENV GOOS=linux
ENV GOARCH=386
RUN make build

FROM asciidoctor/docker-asciidoctor:1.68.0

RUN apk add --no-cache \
    git

COPY --from=0 build/monako /usr/bin/monako
RUN chmod +x /usr/bin/monako

WORKDIR /docs
