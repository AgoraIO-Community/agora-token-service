# Agora Token Webservice
![Go](https://github.com/digitallysavvy/agora-token-server/workflows/Go/badge.svg?branch=master) ![Docker Image CI](https://github.com/digitallysavvy/agora-token-server/workflows/Docker%20Image%20CI/badge.svg?branch=master)  

Written in Golang, using [Gin framework](https://github.com/gin-gonic/gin) to create a RESTful webservice for generating user tokens for use with the [Agora.io](https://www.agora.io) platform. 

Agora.io Advanced Guide: [Token Management](https://docs.agora.io/en/Video/token_server)

## Deploy to Heroku ##
[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://www.heroku.com/deploy/?template=https://github.com/AgoraIO-Community/agora-token-service)

## How to Run ##
Set the APP_ID and APP_CERTIFICATE env variables.
```bash
cp .env.example .env
```

```bash
go run main.go
```

---

The pre-compiled binaries are also available in [releases](https://github.com/maxxfrazer/agora-token-service/releases).

## Docker ##

#1. Open the `Dokerfile` and update the values for `APP_ID` and `APP_CERT`

```bash
ENV APP_ID=""
ENV APP_CERTIFICATE=""
```
#2. To build the container: 

```bash
docker build -t agora-token-service .
```
#3. Run the container 

```bash
docker run agora-token-service
```
> Note: for testing locally
```bash
docker run -p 8080:8080 agora-token-service
```

## Endpoints ##

### Ping ###
**endpoint structure**
```bash
/ping
```
response:
``` json
{"message":"pong"} 
```

### RTC Token ###
The `rtc` token endpoint requires a `tokentype` (uid || userAccount), `channelName`, and the user's `uid` (type varies based on `tokentype`). 
`(optional)` Pass an integer to represent the token lifetime in seconds.

**endpoint structure** 
```
/rtc/:channelName/:role/:tokentype/:uid/?expireTime
```

response:
``` json
{"rtcToken":" "} 
```

## RTM Token ##
The `rtm` token endpoint requires the user's `uid`. 
`(optional)` Pass an integer to represent the privelege lifetime in seconds.
**endpoint structure** 
```
/rtm/:uid/?expireTime
```

response:
``` json
{"rtmToken":" "} 
```

### Both Tokens ###
The `rte` token endpoint generates both the `rtc` and `rtm` tokens with a single request. This endpoint requires a `tokentype` (uid || userAccount), `channelName`, and the user's `uid` (type varies `String/Int` based on `tokentype`). 
`(optional)` Pass an integer to represent the token lifetime in seconds.

**endpoint structure** 
```
/rte/:channelName/:role/:tokentype/:uid/?expireTime
```

response:
``` json
{
  "rtcToken":" ",
  "rtmToken":" " 
} 
```
