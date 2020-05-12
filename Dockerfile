# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine
RUN apk add git ca-certificates --update

RUN ls $GOPATH

RUN go get github.com/gin-gonic/gin
ADD . /go/src/github.com/digitallysavvy/agora-token-server

ENV APP_ID=""
ENV APP_CERT=""

# move to the working directory
WORKDIR $GOPATH/src/github.com/digitallysavvy/agora-token-server
# Build the token server command inside the container.
RUN go build
# RUN go run main.go
# Run the token server by default when the container starts.
ENTRYPOINT ./agora-token-server

# Document that the service listens on port 8080.
EXPOSE 8080