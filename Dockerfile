FROM golang:1.13-alpine

RUN apk add --no-cache --virtual .build-deps \
            gcc \
            libc-dev \
            git

RUN go get -u golang.org/x/lint/golint

CMD ["/usr/local/go/bin/go","version"]
