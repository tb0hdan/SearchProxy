# build stage
FROM golang:alpine AS build-env
WORKDIR /searchproxy
ADD ./ssl /etc/ssl
ADD . /searchproxy/
RUN apk update
RUN apk add git make gcc libc-dev
RUN make

# final stage
FROM alpine
WORKDIR /
COPY --from=build-env /etc/ssl /etc/ssl
COPY --from=build-env /searchproxy/searchproxy /
COPY --from=build-env /searchproxy/mirrors.yml /
EXPOSE 8000
ENTRYPOINT /searchproxy
