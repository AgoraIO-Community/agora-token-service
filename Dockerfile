# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine
RUN apk add git ca-certificates --update

# fetch dependancies github
RUN go get go.mod 

ENV APP_ID=""
ENV APP_CERTIFICATE=""

# move to the working directory
WORKDIR $GOPATH/src/github.com/AgoraIO-Community/agora-token-server
# Build the token server command inside the container.
RUN go build -o agora-token-server -v cmd/main.go
# RUN go run main.go
# Run the token server by default when the container starts.
ENTRYPOINT ./agora-token-server

# Document that the service listens on port 8080.
EXPOSE 8080