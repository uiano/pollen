# Build
FROM golang:alpine AS build-env

RUN apk --no-cache add build-base git mercurial gcc

ADD . /src

RUN cd /src && go build cmd/server.go

# Deploy
FROM alpine

WORKDIR /app

COPY --from=build-env /src/server /app/server
COPY --from=build-env /src/misc /app/misc/

EXPOSE 80
EXPOSE 443

ENTRYPOINT ./server --serve