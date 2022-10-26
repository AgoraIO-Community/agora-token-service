package service

import (
	"log"

	"github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"

	"github.com/gin-gonic/gin"
)

func (s *Service) getRtcToken(c *gin.Context) {
	log.Println("Generating RTC token")
	// get param values
	channelName, tokenType, uidStr, _, role, expireTimestamp, err := s.parseRtcParams(c)

	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Error Generating RTC token: " + err.Error(),
			"status":  400,
		})
		return
	}

	rtcToken, tokenErr := s.generateRtcToken(channelName, uidStr, tokenType, role, expireTimestamp)

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

func (s *Service) getRtmToken(c *gin.Context) {
	log.Println("Generating RTM token")
	// get param values
	uidStr, expireTimestamp, err := s.parseRtmParams(c)

	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Error Generating RTC token: " + err.Error(),
			"status":  400,
		})
		return
	}

	rtmToken, tokenErr := rtmtokenbuilder.BuildToken(s.appID, s.appCertificate, uidStr, rtmtokenbuilder.RoleRtmUser, expireTimestamp)

	if tokenErr != nil {
		log.Println(tokenErr) // token failed to generate
		c.Error(tokenErr)
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

func (s *Service) getRtcRtmToken(c *gin.Context) {
	log.Println("Generating RTC and RTM tokens")
	// get rtc param values
	channelName, tokenType, uidStr, rtmuid, role, expireTimestamp, rtcParamErr := s.parseRtcParams(c)

	if rtcParamErr != nil {
		c.Error(rtcParamErr)
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Error Generating RTC token: " + rtcParamErr.Error(),
			"status":  400,
		})
		return
	}
	// generate the rtcToken
	rtcToken, rtcTokenErr := s.generateRtcToken(channelName, uidStr, tokenType, role, expireTimestamp)
	// generate rtmToken
	rtmToken, rtmTokenErr := rtmtokenbuilder.BuildToken(s.appID, s.appCertificate, rtmuid, rtmtokenbuilder.RoleRtmUser, expireTimestamp)

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

func (s *Service) nocache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// set headers
		c.Header("Cache-Control", "private, no-cache, no-store, must-revalidate")
		c.Header("Expires", "-1")
		c.Header("Pragma", "no-cache")
		c.Header("Access-Control-Allow-Origin", "*")
	}
}
