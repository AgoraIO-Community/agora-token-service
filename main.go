package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/digitallysavvy/agora-token-server/rtctokenbuilder"
	"github.com/digitallysavvy/agora-token-server/rtmtokenbuilder"

	"github.com/gin-gonic/gin"
)

func main() {

	appID, appIDExists := os.LookupEnv("APP_ID")
	appCertificate, appCertExists := os.LookupEnv("APP_CERTIFICATE")

	if !appIDExists || !appCertExists {
		log.Fatal("FATAL ERROR: ENV not properly configured, check appID and appCertificate")
	}

	api := gin.Default()

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// This handler will match  with or without a tokentype
	api.GET("/token/:channelName/:uid/", func(c *gin.Context) {
		// get param values
		channelName := c.Param("channelName")
		uidStr := c.Param("uid")
		expireTime := c.DefaultQuery("expiry", "3600")

		// declare vars
		var rtcToken, rtmToken string // token strings
		var err error                 // catch-all error

		// check if uid is set to 0
		if uidStr == "0" {
			uidStr = ""
		}

		expireTime64, err := strconv.ParseUint(expireTime, 10, 64)
		// check if string conversion fails
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(400, gin.H{
				"message": "expireTime conversion error",
				"status":  400,
			})
			return
		}

		// set timestamps
		expireTimeInSeconds := uint32(expireTime64)
		currentTimestamp := uint32(time.Now().UTC().Unix())
		expireTimestamp := currentTimestamp + expireTimeInSeconds

		rtcToken, err = rtctokenbuilder.BuildTokenWithUserAccount(appID, appCertificate, channelName, uidStr, rtctokenbuilder.RoleAttendee, expireTimestamp)

		if err != nil {
			log.Println(err) // token failed to generate
			c.Error(err)
			c.AbortWithStatusJSON(400, gin.H{
				"error":  err,
				"status": 400,
			})
		}

		rtmToken, err = rtmtokenbuilder.BuildToken(appID, appCertificate, uidStr, rtmtokenbuilder.RoleRtmUser, expireTimestamp)

		if err != nil {
			log.Println(err) // token failed to generate
			c.Error(err)
			c.AbortWithStatusJSON(400, gin.H{
				"error":  err,
				"status": 400,
			})
		} else {
			log.Println("Token generated")
			c.JSON(200, gin.H{
				"rtcToken": rtcToken,
				"rtmToken": rtmToken,
			})
		}

	})

	// listen and serve on localhost:8080
	api.Run(":8080")
}
