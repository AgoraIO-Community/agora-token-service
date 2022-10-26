package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	"github.com/gin-gonic/gin"
)

func (s *Service) parseRtcParams(c *gin.Context) (channelName, tokenType, uidStr string, rtmuid string, role rtctokenbuilder.Role, expireTimestamp uint32, err error) {
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
	if rtmuid == "" {
		if uidStr == "0" {
			err = fmt.Errorf("Failed to parse rtm user ID. Cannot be empty or \"0\"")
		}
		rtmuid = uidStr
	}

	if roleStr == "publisher" {
		role = rtctokenbuilder.RolePublisher
	} else {
		// Making an assumption that !publisher == subscriber
		role = rtctokenbuilder.RoleSubscriber
	}

	expireTime := c.DefaultQuery("expiry", "3600")
	expireTime64, parseErr := strconv.ParseUint(expireTime, 10, 64)
	if parseErr != nil {
		// if string conversion fails return an error
		if err != nil {
			err = fmt.Errorf("%s. Also failed to parse expireTime: %s, causing error: %s", err, expireTime, parseErr)
		} else {
			err = fmt.Errorf("Failed to parse expireTime: %s, causing error: %s", expireTime, parseErr)
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

	// set timestamps
	expireTimeInSeconds := uint32(expireTime64)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp = currentTimestamp + expireTimeInSeconds

	// check if string conversion fails
	return uidStr, expireTimestamp, err
}
