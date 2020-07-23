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
docker build -t agora-token-service .
```
#3. Run the container 
```
docker run agora-token-service
```
> Note: for testing locally
```
docker run -p 8080:8080 agora-token-service
```

## Endpoints ##

### Ping ###
**endpoint structure**
```
/ping
```
response:
``` 
{"message":"pong"} 
```

### Token ###
The token endpoint requires a `channelName`, the user's `role` (subscriber/publisher), and the user's `uid` to generate both RTC and RTM tokens. 
`(optional)` Pass an integer to represent the token privilege lifetime in seconds.

**endpoint structure** 
```
/token/:channelName/:role/:uid/?expireTime
```

response:
``` 
{
  "rtcToken":" ",
  "rtmToken":" ",
} 
```