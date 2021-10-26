FROM golang:1.15-alpine3.11 as build
#FROM golang:latest as build

WORKDIR /app

COPY . .
ARG netrc
ARG release
ENV CGO_ENABLED=0 RELEASE=$release NETRC=$netrc GOSUMDB=off


RUN echo $NETRC | base64 -d > ~/.netrc &&  \
    apk update && apk upgrade && \
    apk add --no-cache bash git openssh && \
    go clean --modcache && \
    go build -o app cmd/app/main.go


FROM alpine:3
COPY --from=build /app/app /usr/local/bin/app

ENV SPLIT_APP_PORT=3001

EXPOSE $SPLIT_APP_PORT

ENTRYPOINT [ "/usr/local/bin/app" ]
