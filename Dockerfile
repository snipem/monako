# run: make image
FROM asciidoctor/docker-asciidoctor

RUN apk add --no-cache \
    git

ADD builds/linux/monako /usr/bin/monako
RUN chmod +x /usr/bin/monako

WORKDIR /docs