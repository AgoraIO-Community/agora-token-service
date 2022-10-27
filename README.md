# Agora Token Webservice
![Go](https://github.com/AgoraIO-Community/agora-token-service/workflows/Go/badge.svg?branch=master) ![Docker Image CI](https://github.com/AgoraIO-Community/agora-token-service/workflows/Docker%20Image%20CI/badge.svg?branch=master)  

Written in Golang, using [Gin framework](https://github.com/gin-gonic/gin) to create a RESTful webservice for generating user tokens for use with the [Agora.io](https://www.agora.io) platform. 

Agora Advanced Guide: [Token Management](https://docs.agora.io/en/video-calling/develop/authentication-workflow).

## Deploy to Railway.app ##

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template/NKYzQA?referralCode=waRWUT)

## Deploy to Render ##

<a href="https://render.com/deploy?repo=https://github.com/AgoraIO-Community/agora-token-service">
  <img src="https://render.com/images/deploy-to-render-button.svg" alt="Deploy to Render">
</a>

## Deploy to Heroku ##
[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://www.heroku.com/deploy/?template=https://github.com/AgoraIO-Community/agora-token-service)

## How to Run ##

Set the APP_ID and APP_CERTIFICATE env variables.

```bash
cp .env.example .env
```

```bash
go run cmd/main.go
```

Without using `.env`, you can also set the environment variables as such:

```bash
APP_ID=app_id APP_CERTIFICATE=app_cert go run cmd/main.go
```

---

The pre-compiled binaries are also available in [releases](https://github.com/AgoraIO-Community/agora-token-service/releases).

## Docker ##

#1. To build the container, with app id and certificate: 

```bash
docker build -t agora-token-service --build-arg APP_ID=$APP_ID APP_CERTIFICATE=$APP_CERTIFICATE .
```

#2. Run the container 

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
The `rtc` token endpoint requires a `tokenType` (uid || userAccount), `channelName`, and the user's `uid` (type varies based on `tokenType`). 
`expiry(optional)` Pass an integer to represent the token lifetime in seconds.

**endpoint structure** 
```
/rtc/:channelName/:role/:tokenType/:rtcuid/?expiry=3600
```

response:
``` json
{"rtcToken":" "} 
```

## RTM Token ##
The `rtm` token endpoint requires the user's `rtmuid`. 
`expiry(optional)` Pass an integer to represent the privelege lifetime in seconds.
**endpoint structure** 
```
/rtm/:rtmuid/?expiry=3600
```

response:
``` json
{"rtmToken":" "} 
```

### Both Tokens ###
The `rte` token endpoint generates both the `rtc` and `rtm` tokens with a single request. This endpoint requires a `tokenType` (uid || userAccount), `channelName`, the user's `rtcuid` (type varies `String/Int` based on `tokenType`) and `rtmuid` which is a `String`. Omitting `rtmuid` will assume it's the same as `rtcuid`.
`expiry(optional)` Pass an integer to represent the token lifetime in seconds.

**endpoint structure** 
```
/rte/:channelName/:role/:tokenType/:rtcuid/:rtmuid/?expiry=3600
```

response:
``` json
{
  "rtcToken":"rtc-token-djfkaljdla",
  "rtmToken":"rtm-token-djfkaljdla" 
} 
```
