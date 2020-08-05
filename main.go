package main

import (
	"fmt"
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
	api.GET("rtc/:channelName/:role/:tokentype/:uid/", getRtcToken)
	api.GET("rtm/:uid/", getRtmToken)
	api.GET("rte/:channelName/:role/:tokentype/:uid/", getBothTokens)
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

func validateRtcParams(c *gin.Context) (channelName, roleStr, tokentype, uidStr string, expireTime64 uint64, err error) {
	// get param values
	channelName = c.Param("channelName")
	roleStr = c.Param("role")
	tokentype = c.Param("tokentype")
	uidStr = c.Param("uid")
	expireTime := c.DefaultQuery("expiry", "3600")

	var parseErr error
	expireTime64, parseErr = strconv.ParseUint(expireTime, 10, 64)
	if err != nil {
		// with formatting
		err = fmt.Errorf("failed to parse expireTime: %s, causing error: %s", expireTime, parseErr)
	}
	// check if string conversion fails
	return channelName, roleStr, tokentype, uidStr, expireTime64, err
}

func validateRtmParams(c *gin.Context) (uidStr string, expireTime64 uint64, err error) {
	// get param values
	uidStr = c.Param("uid")
	expireTime := c.DefaultQuery("expiry", "3600")

	var parseErr error
	expireTime64, parseErr = strconv.ParseUint(expireTime, 10, 64)
	if err != nil {
		// with formatting
		err = fmt.Errorf("failed to parse expireTime: %s, causing error: %s", expireTime, parseErr)
	}
	// check if string conversion fails
	return uidStr, expireTime64, err
}

func generateRtmToken(uidStr string, expireTime64 uint64) (rtmToken string, err error) {
	// set timestamps
	expireTimeInSeconds := uint32(expireTime64)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp := currentTimestamp + expireTimeInSeconds

	rtmToken, tokenErr := rtmtokenbuilder.BuildToken(appID, appCertificate, uidStr, rtmtokenbuilder.RoleRtmUser, expireTimestamp)
	return rtmToken, tokenErr
}

func generateRtcToken(channelName, uidStr, roleStr, tokentype string, expireTime64 uint64) (rtcToken string, err error) {
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

	if tokentype == "userAccount" {
		log.Printf("Building Token with userAccount: %s\n", uidStr)
		rtcToken, err = rtctokenbuilder.BuildTokenWithUserAccount(appID, appCertificate, channelName, uidStr, role, expireTimestamp)
		return rtcToken, err

	} else if tokentype == "uid" {
		uid64, parseErr := strconv.ParseUint(uidStr, 10, 64)
		// check if conversion fails
		if parseErr != nil {
			err = fmt.Errorf("failed to parse uidStr: %s, to uint causing error: %s", uidStr, parseErr)
			return "", err
		}

		uid := uint32(uid64) // convert uid from uint64 to uint 32
		log.Printf("Building Token with uid: %d\n", uid)
		rtcToken, err = rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, uid, role, expireTimestamp)
		return rtcToken, err

	} else {
		err = fmt.Errorf("failed to generate RTC token for Unknown Tokentype: %s", tokentype)
		log.Println(err)
		return "", err
	}
}

func getRtcToken(c *gin.Context) {
	log.Printf("rtc token\n")
	// get param values
	channelName, roleStr, tokentype, uidStr, expireTime64, err := validateRtcParams(c)

	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Error Generating RTC token: " + err.Error(),
			"status":  400,
		})
		return
	}

	rtcToken, tokenErr := generateRtcToken(channelName, uidStr, roleStr, tokentype, expireTime64)

	if tokenErr != nil {
		log.Println(tokenErr) // token failed to generate
		c.Error(tokenErr)
		errMsg := "Error Generating RTC token - " + tokenErr.Error()
		c.AbortWithStatusJSON(400, gin.H{
			"status": 400,
			"error":  errMsg,
		})
	} else {
		log.Println("RTC Token generated")
		c.JSON(200, gin.H{
			"rtcToken": rtcToken,
		})
	}
}

func getRtmToken(c *gin.Context) {
	log.Printf("rtm token\n")
	// get param values
	uidStr, expireTime64, err := validateRtmParams(c)

	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Error Generating RTC token: " + err.Error(),
			"status":  400,
		})
		return
	}

	rtmToken, tokenErr := generateRtmToken(uidStr, expireTime64)

	if tokenErr != nil {
		log.Println(err) // token failed to generate
		c.Error(err)
		errMsg := "Error Generating RTM token: " + tokenErr.Error()
		c.AbortWithStatusJSON(400, gin.H{
			"error":  errMsg,
			"status": 400,
		})
	} else {
		log.Println("RTM Token generated")
		c.JSON(200, gin.H{
			"rtmToken": rtmToken,
		})
	}
}

func getBothTokens(c *gin.Context) {
	log.Printf("dual token\n")
	// get rtc param values
	channelName, roleStr, tokentype, uidStr, expireTime64, rtcParamErr := validateRtcParams(c)

	if rtcParamErr != nil {
		c.Error(rtcParamErr)
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Error Generating RTC token: " + rtcParamErr.Error(),
			"status":  400,
		})
		return
	}
	// generate the rtcToken
	rtcToken, rtcTokenErr := generateRtcToken(channelName, uidStr, roleStr, tokentype, expireTime64)
	// generate rtmToken
	rtmToken, rtmTokenErr := generateRtmToken(uidStr, expireTime64)

	if rtcTokenErr != nil {
		log.Println(rtcTokenErr) // token failed to generate
		c.Error(rtcTokenErr)
		errMsg := "Error Generating RTC token - " + rtcTokenErr.Error()
		c.AbortWithStatusJSON(400, gin.H{
			"status": 400,
			"error":  errMsg,
		})
	} else if rtmTokenErr != nil {
		log.Println(rtmTokenErr) // token failed to generate
		c.Error(rtmTokenErr)
		errMsg := "Error Generating RTC token - " + rtmTokenErr.Error()
		c.AbortWithStatusJSON(400, gin.H{
			"status": 400,
			"error":  errMsg,
		})
	} else {
		log.Println("RTC Token generated")
		c.JSON(200, gin.H{
			"rtcToken": rtcToken,
			"rtmToken": rtmToken,
		})
	}

}
