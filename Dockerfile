# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine
RUN apk add git ca-certificates --update

RUN ls $GOPATH

RUN go get github.com/gin-gonic/gin
RUN ls /go/src/github.com/
# RUN GO111MODULE=on go get -v github.com/digitallysavvy/agora-token-server

ADD . /go/src/github.com/digitallysavvy/agora-token-server
# ADD /RtcTokenBuilder /go/src/github.com/digitallysavvy/RtcTokenBuilder
# ADD /AccessToken /go/src/github.com/digitallysavvy/AccessToken
RUN ls -l /go/src/github.com/digitallysavvy/agora-token-server

ENV APP_ID=""
ENV APP_CERT=""

# Build the goingout command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
WORKDIR $GOPATH/src/github.com/digitallysavvy/agora-token-server
# RUN go build github.com/digitallysavvy/agora-token-server
# RUN GO111MODULE=on go build github.com/digitallysavvy/agora-token-server
RUN go build
# RUN go run main.go
# Run the goingout command by default when the container starts.
ENTRYPOINT /go/bin/agora-token-server

# Document that the service listens on port 8080.
EXPOSE 8080