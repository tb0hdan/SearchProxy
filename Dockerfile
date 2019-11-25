# build stage
FROM golang:alpine AS build-env
ENV GOPATH /tmp
WORKDIR /searchproxy
ADD ./ssl /etc/ssl
ADD ./Makefile /searchproxy
ADD ./mirrors.yml /searchproxy
ADD ./httputil /searchproxy
ADD ./memcache /searchproxy
ADD ./mirrorsort /searchproxy
ADD ./workerpool /searchproxy
ADD ./go.mod /searchproxy
ADD ./go.sum /searchproxy
RUN apk update
RUN apk add git
RUN apk add make
RUN apk add gcc
RUN apk add libc-dev
RUN cd /searchproxy && make

# final stage
FROM alpine
WORKDIR /
COPY --from=build-env /etc/ssl /etc/ssl
COPY --from=build-env /src/searchproxy /
COPY --from=build-env /src/mirrors.yml /
EXPOSE 8000
ENTRYPOINT /searchproxy
