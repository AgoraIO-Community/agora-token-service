FROM golang:alpine
RUN apk add git ca-certificates --update

COPY .env .
ADD . /go/src/github.com/AgoraIO-Community/agora-token-service

ENV GIN_MODE=release

# move to the working directory
WORKDIR $GOPATH/src/github.com/AgoraIO-Community/agora-token-service
# Build the token server command inside the container.
RUN go build -o agora-token-service -v cmd/main.go