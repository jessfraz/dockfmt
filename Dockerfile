FROM golang:alpine as builder
MAINTAINER Jessica Frazelle <jess@linux.com>

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN	apk add --no-cache \
	bash \
	ca-certificates

COPY . /go/src/github.com/jessfraz/dockfmt

RUN set -x \
	&& apk add --no-cache --virtual .build-deps \
		git \
		gcc \
		libc-dev \
		libgcc \
		make \
	&& cd /go/src/github.com/jessfraz/dockfmt \
	&& make static \
	&& mv dockfmt /usr/bin/dockfmt \
	&& apk del .build-deps \
	&& rm -rf /go \
	&& echo "Build complete."

FROM alpine:latest

COPY --from=builder /usr/bin/dockfmt /usr/bin/dockfmt
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs

ENTRYPOINT [ "dockfmt" ]
CMD [ "--help" ]
