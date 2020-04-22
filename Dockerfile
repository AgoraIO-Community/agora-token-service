# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine
RUN apk add git ca-certificates --update

RUN ls $GOPATH

RUN go get github.com/gin-gonic/gin
RUN ls /go/src/github.com/
# RUN GO111MODULE=on go get -v github.com/digitallysavvy/agora-token-server



ADD . /go/src/github.com/digitallysavvy/agora-token-server
RUN ls /go/src/github.com/digitallysavvy/agora-token-server

ENV APP_ID="4fdfd402ce0a45ea94d850f2124f0b36"
ENV APP_CERT="01e25afb21d14b7b832a7e2b0a43d5e1"

# Build the goingout command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
# RUN go build github.com/digitallysavvy/agora-token-server

# Run the goingout command by default when the container starts.
ENTRYPOINT /go/bin/agora-token-server

# Document that the service listens on port 8080.
EXPOSE 8080