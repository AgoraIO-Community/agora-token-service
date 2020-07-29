package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	"github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"

	"github.com/gin-gonic/gin"
)

var appID string
var appCertificate string

func main() {

	appIDEnv, appIDExists := os.LookupEnv("APP_ID")
	appCertEnv, appCertExists := os.LookupEnv("APP_CERTIFICATE")

	if !appIDExists || !appCertExists {
		log.Fatal("FATAL ERROR: ENV not properly configured, check appID and appCertificate")
	} else {
		appID = appIDEnv
		appCertificate = appCertEnv
	}

	api := gin.Default()

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api.Use(nocache())
	api.GET("rtc/:tokentype/:channelName/:uid/:role/", getRtcToken)
	api.GET("rtm/:uid/", getRtmToken)
	api.GET("rte/:tokentype/:channelName/:uid/:role/", getBothTokens)
	api.Run(":8080") // listen and serve on localhost:8080
}

func nocache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// set headers
		c.Header("Cache-Control", "private, no-cache, no-store, must-revalidate")
		c.Header("Expires", "-1")
		c.Header("Pragma", "no-cache")
		c.Header("Access-Control-Allow-Origin", "*")
	}
}

func getBothTokens(c *gin.Context) {
	getRtcToken(c)
	// check if first token succeeded
	if !c.IsAborted() {
		getRtmToken(c)
	}

}

func getRtcToken(c *gin.Context) {
	// get param values
	channelName := c.Param("channelName")
	uidStr := c.Param("uid")
	roleStr := c.Param("role")
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
			"message": "Error Generating RTC token: expireTime conversion error",
			"status":  400,
		})
		return
	}

	var role rtctokenbuilder.Role
	if roleStr == "publisher" {
		role = rtctokenbuilder.RolePublisher
	} else {
		role = rtctokenbuilder.RoleSubscriber
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
				"message": "Error Generating RTC token: UID conversion error.",
				"status":  400,
			})
			return
		}
		uid := uint32(uid64) // convert uid from uint64 to uint 32
		log.Printf("\nBuilding Token with uid: %d\n", uid)
		result, err = rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, uid, role, expireTimestamp)
	} else if tokentype == "userAccount" {
		log.Printf("\nBuilding Token with userAccount: %s\n", uidStr)
		result, err = rtctokenbuilder.BuildTokenWithUserAccount(appID, appCertificate, channelName, uidStr, role, expireTimestamp)
	} else {
		errMsg := "Error Generating RTC token: Unknown Tokentype: " + tokentype
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
		errMsg := "Error Generating RTC token: " + err.Error()
		c.AbortWithStatusJSON(400, gin.H{
			"status": 400,
			"error":  errMsg,
		})
	} else {
		log.Println("Token generated")
		c.JSON(200, gin.H{
			"rtcToken": result,
		})
	}
}

func getRtmToken(c *gin.Context) {
	// get param values
	uidStr := c.Param("uid")
	expireTime := c.DefaultQuery("expiry", "3600")

	log.Printf("rtm token\n")

	// declare vars
	var result string // token string
	var err error     // catch-all error

	expireTime64, err := strconv.ParseUint(expireTime, 10, 64)
	// check if string conversion fails
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Error Generating RTM token: expireTime conversion error",
			"status":  400,
		})
		return
	}

	// set timestamps
	expireTimeInSeconds := uint32(expireTime64)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp := currentTimestamp + expireTimeInSeconds

	result, err = rtmtokenbuilder.BuildToken(appID, appCertificate, uidStr, rtmtokenbuilder.RoleRtmUser, expireTimestamp)

	if err != nil {
		log.Println(err) // token failed to generate
		c.Error(err)
		errMsg := "Error Generating RTM token: " + err.Error()
		c.AbortWithStatusJSON(400, gin.H{
			"error":  errMsg,
			"status": 400,
		})
	} else {
		log.Println("Token generated")
		c.JSON(200, gin.H{
			"rtmToken": result,
		})
	}
}
