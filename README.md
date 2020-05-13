# agora-token-server
![Go](https://github.com/digitallysavvy/agora-token-server/workflows/Go/badge.svg?branch=master)
Written in Golang, using Gin middleware to create a RESTful webservice for generating user tokens for use with the Agora.io platform.

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
