# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine
RUN apk add git ca-certificates --update

# fetch dependancies from github (Gin and Agora Token Service)
RUN go get github.com/gin-gonic/gin
# RUN go get github.com/AgoraIO-Community/agora-token-service
ADD . /go/src/github.com/AgoraIO-Community/agora-token-service

ENV APP_ID=""
ENV APP_CERTIFICATE=""

# move to the working directory
WORKDIR $GOPATH/src/github.com/AgoraIO-Community/agora-token-service
# Build the token server command inside the container.
RUN go build
# RUN go run main.go
# Run the token server by default when the container starts.
ENTRYPOINT ./agora-token-service

# Document that the service listens on port 8080.
EXPOSE 8080