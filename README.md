# Agora Token Webservice

<p align="center">
  <img src="https://github.com/AgoraIO-Community/agora-token-service/workflows/Go/badge.svg?branch=main">
  <img src="https://github.com/AgoraIO-Community/agora-token-service/actions/workflows/dockerimage.yml/badge.svg?branch=main">
  <a href="https://github.com/AgoraIO-Community/agora-token-service/releases/latest"><img src="https://github.com/AgoraIO-Community/agora-token-service/actions/workflows/release.yml/badge.svg?release=latest"></a>
</p>

Written in Golang, using [Gin framework](https://github.com/gin-gonic/gin) to create a RESTful webservice for generating user tokens for use with the [Agora.io](https://www.agora.io) platform. 

Agora Advanced Guide: [Token Management](https://docs.agora.io/en/video-calling/develop/authentication-workflow).

## One-Click Deployments

| Railway | Render | Heroku |
|:-:|:-:|:-:|
| [![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template/NKYzQA?referralCode=waRWUT) | [![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/AgoraIO-Community/agora-token-service) | [![Deploy to Heroku](https://www.herokucdn.com/deploy/button.svg)](https://www.heroku.com/deploy/?template=https://github.com/AgoraIO-Community/agora-token-service) |

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
{"rtcToken":"007rtc-token-djfkaljdla"} 
```

### RTM Token ###

The `rtm` token endpoint requires the user's `rtmuid`. 
`expiry(optional)` Pass an integer to represent the privelege lifetime in seconds.
**endpoint structure** 
```
/rtm/:rtmuid/?expiry=3600
```

response:
``` json
{"rtmToken":"007rtm-token-djfkaljdla"} 
```

### RTM + RTC Tokens ###
The `rte` token endpoint generates both the `rtc` and `rtm` tokens with a single request. This endpoint requires a `tokenType` (uid || userAccount), `channelName`, the user's `rtcuid` (type varies `String/Int` based on `tokenType`) and `rtmuid` which is a `String`. Omitting `rtmuid` will assume it's the same as `rtcuid`.
`expiry(optional)` Pass an integer to represent the token lifetime in seconds.

**endpoint structure** 
```
/rte/:channelName/:role/:tokenType/:rtcuid/:rtmuid/?expiry=3600
```

response:
``` json
{
  "rtcToken":"007rtc-token-djfkaljdla",
  "rtmToken":"007rtm-token-djfkaljdla" 
} 
```

### Chat Tokens ###

#### endpoint structure ####

app privileges:
```
chat/app/?expiry=3600
```

user privileges:
```
/chat/account/:chatid/?expiry=3600
```

`expiry` is an optional parameter for both.

response:
``` json
{
  "chatToken":"007chat-token-djfkaljdla"
} 
```

## Contributions

Contributions are welcome, please test any changes to the Go code with the following command:

```sh
APP_ID=<YOUR_APP_ID> APP_CERTIFICATE=<YOUR_APP_CERT> go test -cover github.com/AgoraIO-Community/agora-token-service/service
```
