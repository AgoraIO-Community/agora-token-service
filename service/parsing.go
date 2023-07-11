package service

import (
	"fmt"
	"strconv"
	"time"

	rtctokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	"github.com/gin-gonic/gin"
)

func (s *Service) parseRtcParams(c *gin.Context) (channelName, tokenType, uidStr string, rtmuid string, role rtctokenbuilder2.Role, expireTimestamp uint32, err error) {
	// get param values
	channelName = c.Param("channelName")
	roleStr := c.Param("role")
	tokenType = c.Param("tokenType")
	uidStr = c.Param("rtcuid")
	rtmuid = c.Param("rtmuid")

	if uidStr == "" {
		// If the uid is missing, just set to 0,
		// meaning it allows for any user ID
		uidStr = "0"
	}
	if rtmuid == "" && uidStr != "0" {
		rtmuid = uidStr
	}

	if roleStr == "publisher" {
		role = rtctokenbuilder2.RolePublisher
	} else {
		// Making an assumption that !publisher == subscriber
		role = rtctokenbuilder2.RoleSubscriber
	}

	expireTime := c.DefaultQuery("expiry", "3600")
	expireTime64, parseErr := strconv.ParseUint(expireTime, 10, 64)
	if parseErr != nil {
		// if string conversion fails return an error
		if err != nil {
			err = fmt.Errorf("%s. Also failed to parse expireTime: %s, causing error: %s", err, expireTime, parseErr)
		} else {
			err = fmt.Errorf("failed to parse expireTime: %s, causing error: %s", expireTime, parseErr)
		}
	}

	// set timestamps
	expireTimeInSeconds := uint32(expireTime64)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp = currentTimestamp + expireTimeInSeconds

	return channelName, tokenType, uidStr, rtmuid, role, expireTimestamp, err
}

func (s *Service) parseRtmParams(c *gin.Context) (uidStr string, expireTimestamp uint32, err error) {
	// get param values
	uidStr = c.Param("rtmuid")
	expireTime := c.DefaultQuery("expiry", "3600")

	expireTime64, parseErr := strconv.ParseUint(expireTime, 10, 64)
	if parseErr != nil {
		// if string conversion fails return an error
		err = fmt.Errorf("failed to parse expireTime: %s, causing error: %s", expireTime, parseErr)
	}

	if uidStr == "" || uidStr == "0" {
		err = fmt.Errorf("invalid RTM User ID: \"%s\"", uidStr)
	}

	// set timestamps
	expireTimeInSeconds := uint32(expireTime64)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp = currentTimestamp + expireTimeInSeconds

	// check if string conversion fails
	return uidStr, expireTimestamp, err
}

func (s *Service) parseChatParams(c *gin.Context) (uidStr string, tokenType string, expireTimestamp uint32, err error) {
	// get param values
	uidStr = c.Param("chatid")
	tokenType = c.Param("tokenType")
	expireTime := c.DefaultQuery("expiry", "3600")
	if tokenType == "account" {
		tokenType = "userAccount"
	}
	if tokenType != "userAccount" && tokenType != "app" {
		err = fmt.Errorf("token type must be \"userAccount\" or \"app\"")
		return uidStr, tokenType, expireTimestamp, err
	}
	if uidStr == "" && tokenType != "app" {
		err = fmt.Errorf("userAccount type requires chat ID")
		return uidStr, tokenType, expireTimestamp, err
	}
	expireTime64, parseErr := strconv.ParseUint(expireTime, 10, 64)
	if parseErr != nil {
		// if string conversion fails return an error
		err = fmt.Errorf("failed to parse expireTime: %s, causing error: %s", expireTime, parseErr)
	}
	// set timestamps
	expireTimeInSeconds := uint32(expireTime64)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp = currentTimestamp + expireTimeInSeconds

	// check if string conversion fails
	return uidStr, tokenType, expireTimestamp, err
}
