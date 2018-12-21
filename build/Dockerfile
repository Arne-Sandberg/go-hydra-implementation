FROM golang:1.11-alpine

ENV GO111MODULE=on

RUN apk update && apk upgrade && apk add --no-cache build-base curl git

WORKDIR /go/src/github.com/jonaslimads/go-hydra-implementation

COPY . .

# Once CompileDaemon is a hot reload builder and is not part of the project code, it needs to be out of go.mod
RUN GO111MODULE=off go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -log-prefix=false -build="go install" -command="go-hydra-implementation"