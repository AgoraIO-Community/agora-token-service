# agora-token-server

## How to Run ##
Set the APP_ID and APP_CERT env variables.
```go
go run main.go
```

## Docker ##
Currently Docker is failing to properly build the token server. I'm using Go Modules for linking of submodules. Running into issues with `docker build -t token-server .` 

```
docker build -t token-server .
Sending build context to Docker daemon  31.96MB
Step 1/13 : FROM golang:alpine
 ---> ccda0e5ccbfc
Step 2/13 : RUN apk add git ca-certificates --update
 ---> Using cache
 ---> 8ad6969b201d
Step 3/13 : RUN ls $GOPATH
 ---> Using cache
 ---> 6b6b59444163
Step 4/13 : RUN go get github.com/gin-gonic/gin
 ---> Using cache
 ---> 6bcb0fea804b
Step 5/13 : RUN ls /go/src/github.com/
 ---> Using cache
 ---> 31c2e7ed5a16
Step 6/13 : ADD . /go/src/github.com/digitallysavvy/agora-token-server
 ---> 965c6a11403a
Step 7/13 : RUN ls -l /go/src/github.com/digitallysavvy/agora-token-server
 ---> Running in 50979c187720
total 14732
drwxr-xr-x    2 root     root          4096 Apr 22 04:58 AccessToken
-rw-r--r--    1 root     root          1231 Apr 22 05:51 Dockerfile
-rw-r--r--    1 root     root          1063 Apr 22 03:36 LICENSE
-rw-r--r--    1 root     root           106 Apr 22 05:54 README.md
drwxr-xr-x    2 root     root          4096 Apr 22 04:58 RtcTokenBuilder
-rwxr-xr-x    1 root     root      15050020 Apr 22 05:34 agora-token-server
-rw-r--r--    1 root     root           102 Apr 22 05:20 go.mod
-rw-r--r--    1 root     root          3561 Apr 22 05:45 go.sum
-rw-r--r--    1 root     root          2689 Apr 22 05:31 main.go
Removing intermediate container 50979c187720
 ---> 625913e88cb7
Step 8/13 : ENV APP_ID="4fdfd402ce0a45ea94d850f2124f0b36"
 ---> Running in 1d1e3a8c1739
Removing intermediate container 1d1e3a8c1739
 ---> ee226edcaa04
Step 9/13 : ENV APP_CERT="01e25afb21d14b7b832a7e2b0a43d5e1"
 ---> Running in 4419d8c0eafa
Removing intermediate container 4419d8c0eafa
 ---> 60867263638d
Step 10/13 : WORKDIR $GOPATH/src/github.com/digitallysavvy/agora-token-server
 ---> Running in a1fb8f30b631
Removing intermediate container a1fb8f30b631
 ---> d79693a7c0df
Step 11/13 : RUN go build
 ---> Running in b9a2d9b27525
go: downloading github.com/gin-gonic/gin v1.6.2
go: finding module for package github.com/digitallysavvy/agora-token-server/rtctokenbuilder
go: downloading github.com/gin-contrib/sse v0.1.0
go: downloading gopkg.in/yaml.v2 v2.2.8
go: downloading github.com/ugorji/go v1.1.7
go: downloading github.com/golang/protobuf v1.3.3
go: downloading github.com/mattn/go-isatty v0.0.12
go: downloading github.com/go-playground/validator/v10 v10.2.0
go: downloading github.com/ugorji/go/codec v1.1.7
go: downloading golang.org/x/sys v0.0.0-20200116001909-b77594299b42
go: downloading github.com/leodido/go-urn v1.2.0
go: downloading github.com/go-playground/universal-translator v0.17.0
go: downloading github.com/go-playground/locales v0.13.0
main.go:9:2: no matching versions for query "latest"
The command '/bin/sh -c go build' returned a non-zero code: 1
```