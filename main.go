package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"rtctokenbuilder"

	"github.com/gin-gonic/gin"
)

func main() {

	appID, appIDExists := os.LookupEnv("APP_ID")
	appCertificate, appCertExists := os.LookupEnv("APP_CERT")

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
	api.GET("token/:tokentype/:channelName/:uid/", func(c *gin.Context) {
		// get param values
		channelName := c.Param("channelName")
		uidStr := c.Param("uid")
		tokentype := c.Param("tokentype")
		expireTime := c.DefaultQuery("expiry", "3600")

		log.Printf("tokentype: %s\n", tokentype)

		// declare vars
		var result string // token string
		var err error     // catch-all error

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

		if tokentype == "uid" {
			uid64, err := strconv.ParseUint(uidStr, 10, 64)
			// check if string conversion fails
			if err != nil {
				c.Error(err)
				c.AbortWithStatusJSON(400, gin.H{
					"message": "UID conversion error",
					"status":  400,
				})
				return
			}
			uid := uint32(uid64) // convert uid from uint64 to uint 32
			log.Printf("\nBuilding Token with uid: %d\n", uid)
			result, err = rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, uid, rtctokenbuilder.RoleAttendee, expireTimestamp)
		} else if tokentype == "userAccount" {
			log.Printf("\nBuilding Token with userAccount: %s\n", uidStr)
			result, err = rtctokenbuilder.BuildTokenWithUserAccount(appID, appCertificate, channelName, uidStr, rtctokenbuilder.RoleAttendee, expireTimestamp)
		} else {
			errMsg := "Unknown Tokentype: " + tokentype
			log.Println(errMsg)
			c.AbortWithStatusJSON(400, gin.H{
				"status": 400,
				"error":  errMsg,
			})
			return
		}

		if err != nil {
			log.Println(err) // token failed to generate
			c.Error(err)
			c.AbortWithStatusJSON(400, gin.H{
				"error":  err,
				"status": 400,
			})
		} else {
			log.Println("Token generated\n")
			c.JSON(200, gin.H{
				"token": result,
			})
		}

	})
	api.Run(":8080") // listen and serve on localhost:8080
}
