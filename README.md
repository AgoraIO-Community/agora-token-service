# Agora Token Webservice
![Go](https://github.com/digitallysavvy/agora-token-server/workflows/Go/badge.svg?branch=master) ![Docker Image CI](https://github.com/digitallysavvy/agora-token-server/workflows/Docker%20Image%20CI/badge.svg?branch=master)   
Written in Golang, using [Gin framework](https://github.com/gin-gonic/gin) to create a RESTful webservice for generating user tokens for use with the [Agora.io](https://www.agora.io) platform. 

Agora.io Advanced Guide: [Token Management](https://docs.agora.io/en/Video/token_server_cpp?platform=CPP)

## How to Run ##
Set the APP_ID and APP_CERT env variables.
```go
go run main.go
```

## Docker ##
#1. Open the `Dokerfile` and update the values for `APP_ID` and `APP_CERT`
```
ENV APP_ID=""
ENV APP_CERT=""
```
#2. To build the container: 
```
docker build -t agora-token-server .
```
#3. Run the container 
```
docker run agora-token-server
```
> Note: for testing locally
```
docker run -p 8080:8080 agora-token-server
```

## Endpoints ##

### Ping ###
**endpoint structure**
```/ping```
response:
``` {"message":"pong"} ```

### Token ###
Token endpoint requires a `tokentype` (uid || userAccount), `channelName`, and the user's `uid` (type varies based on `tokentype`). 
`(optional)` Pass an integer to represent the token lifetime in seconds.

**endpoint structure**
```/token/:tokentype/:channelName/:uid/?expireTime```

response:
``` {"token":" "} ```