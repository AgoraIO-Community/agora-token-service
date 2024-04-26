# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:alpine
RUN apk add git ca-certificates --update

# fetch dependencies github
# RUN go get -u github.com/gin-gonic/gin

ADD . /go/src/github.com/AgoraIO-Community/agora-token-service

# # fetch dependencies from github (Gin and Agora Token Service)
# RUN go install github.com/gin-gonic/gin@latest
# # RUN go install github.com/AgoraIO-Community/agora-token-service
# ADD . /go/src/github.com/AgoraIO-Community/agora-token-service

ARG APP_ID
ARG APP_CERTIFICATE
ARG SERVER_PORT
ARG CORS_ALLOW_ORIGIN
ENV APP_ID $APP_ID
ENV APP_CERTIFICATE $APP_CERTIFICATE
ENV SERVER_PORT $SERVER_PORT
ENV CORS_ALLOW_ORIGIN $CORS_ALLOW_ORIGIN

# move to the working directory
WORKDIR $GOPATH/src/github.com/AgoraIO-Community/agora-token-service
# Build the token server command inside the container.
RUN go build -o agora-token-service -v cmd/main.go
# RUN go run main.go
# Run the token server by default when the container starts.
ENTRYPOINT ./agora-token-service

# Document that the service listens on port 8080.
EXPOSE $SERVER_PORT