package service

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	rtmtokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"

	"github.com/gin-gonic/gin"
)

func (s *Service) getRtcToken(c *gin.Context) {
	log.Println("Generating RTC token")
	// get param values
	channelName, tokenType, uidStr, _, role, expire, err := s.parseRtcParams(c)

	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Error Generating RTC token: " + err.Error(),
			"status":  400,
		})
		return
	}

	rtcToken, tokenErr := s.generateRtcToken(channelName, uidStr, tokenType, role, expire)

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
	uidStr, expire, err := s.parseRtmParams(c)

	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Error Generating RTM token: " + err.Error(),
			"status":  400,
		})
		return
	}

	rtmToken, tokenErr := rtmtokenbuilder2.BuildToken(s.appID, s.appCertificate, uidStr, expire, "")

	if tokenErr != nil {
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

func (s *Service) getChatToken(c *gin.Context) {
	log.Println("Generating Chat token")
	// get param values
	uidStr, tokenType, expireTimestamp, err := s.parseChatParams(c)

	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Error Generating Chat token: " + err.Error(),
			"status":  400,
		})
		return
	}

	chatToken, tokenErr := s.generateChatToken(uidStr, tokenType, expireTimestamp)

	if tokenErr != nil {
		c.Error(tokenErr)
		errMsg := "Error Generating Chat token: " + tokenErr.Error()
		c.AbortWithStatusJSON(400, gin.H{
			"error":  errMsg,
			"status": 400,
		})
	} else {
		log.Println("Chat Token generated")
		c.JSON(200, gin.H{
			"chatToken": chatToken,
		})
	}
}

func (s *Service) getRtcRtmToken(c *gin.Context) {
	log.Println("Generating RTC and RTM tokens")
	// get rtc param values
	channelName, tokenType, uidStr, rtmuid, role, expire, rtcParamErr := s.parseRtcParams(c)

	if rtcParamErr == nil && rtmuid == "" {
		rtcParamErr = fmt.Errorf("failed to parse rtm user ID. Cannot be empty or \"0\"")
	}
	if rtcParamErr != nil {
		c.Error(rtcParamErr)
		c.AbortWithStatusJSON(400, gin.H{
			"message": "Error Generating RTC and RTM token: " + rtcParamErr.Error(),
			"status":  400,
		})
		return
	}
	// generate the rtcToken
	rtcToken, rtcTokenErr := s.generateRtcToken(channelName, uidStr, tokenType, role, expire)
	// generate rtmToken
	rtmToken, rtmTokenErr := rtmtokenbuilder2.BuildToken(s.appID, s.appCertificate, rtmuid, expire, channelName)

	if rtcTokenErr != nil {
		c.Error(rtcTokenErr)
		errMsg := "Error Generating RTC token - " + rtcTokenErr.Error()
		c.AbortWithStatusJSON(400, gin.H{
			"status": 400,
			"error":  errMsg,
		})
	} else if rtmTokenErr != nil {
		c.Error(rtmTokenErr)
		errMsg := "Error Generating RTM token - " + rtmTokenErr.Error()
		c.AbortWithStatusJSON(400, gin.H{
			"status": 400,
			"error":  errMsg,
		})
	} else {
		log.Println("RTC and RTM Tokens generated")
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
	}
}

// Add CORSMiddleware to handle CORS requests and set the necessary headers
func (s *Service) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if !s.isOriginAllowed(origin) {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Origin not allowed",
			})
			c.Abort()
			return
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func (s *Service) isOriginAllowed(origin string) bool {
	if s.allowOrigin == "*" {
		return true
	}

	allowedOrigins := strings.Split(s.allowOrigin, ",")
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}

	return false
}
