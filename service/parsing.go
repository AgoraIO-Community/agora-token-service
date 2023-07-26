package service

import (
	"fmt"
	"strconv"
	"strings"

	rtctokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	"github.com/gin-gonic/gin"
)

// parseRtcParams extracts and parses the required parameters from the given Gin context (HTTP request).
// It retrieves various parameters such as channelName, tokenType, rtcuid, rtmuid, role, and expiry,
// and performs necessary conversions and validations.
//
// Parameters:
//
//	c: *gin.Context - The Gin context representing the HTTP request and response.
//
// Returns:
//   - channelName: string - The name of the video conferencing channel.
//   - tokenType: string - The type of RTC token.
//   - uidStr: string - The user ID for the RTC token.
//   - rtmuid: string - The user ID for the RTM (Real-Time Messaging) token.
//   - role: rtctokenbuilder2.Role - The role of the user (Publisher or Subscriber) in the video conferencing.
//   - expire: uint32 - The expiration time of the token in seconds.
//   - err: error - Any error that occurred during parameter parsing. Nil if parsing was successful.
//
// Behavior:
//  1. Retrieves the values of channelName, roleStr, tokenType, rtcuid, and rtmuid from the Gin context.
//  2. Sets uidStr to "0" if it is empty, implying that any user ID is allowed.
//  3. If rtmuid is empty and uidStr is not "0", it sets rtmuid to uidStr.
//  4. Determines the role based on the value of roleStr. "publisher" maps to RolePublisher, else RoleSubscriber.
//  5. Parses the expiry time from the query parameter "expiry" and converts it to uint32.
//  6. If string conversion fails for the expiry time, sets err to an error with the failure information.
//
// Notes:
//   - The `rtctokenbuilder2.Role` type represents the role of a user in the video conferencing token.
//     It might be an enumeration or constant representing different roles (e.g., Publisher and Subscriber).
//
// Example usage:
//
//	channelName, tokenType, uidStr, rtmuid, role, expire, err := parseRtcParams(context)
func (s *Service) parseRtcParams(c *gin.Context) (channelName, tokenType, uidStr string, rtmuid string, role rtctokenbuilder2.Role, expire uint32, err error) {
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
		err = fmt.Errorf("failed to parse expireTime: %s, causing error: %s", expireTime, parseErr)
	}
	expire = uint32(expireTime64)

	return channelName, tokenType, uidStr, rtmuid, role, expire, err
}

func (s *Service) parseRtmParams(c *gin.Context) (uidStr string, expire uint32, err error) {
	// get param values
	uidStr = c.Param("rtmuid")
	expireTime := c.DefaultQuery("expiry", "3600")
	expireTime64, parseErr := strconv.ParseUint(expireTime, 10, 64)
	if parseErr != nil {
		// if string conversion fails return an error
		err = fmt.Errorf("failed to parse expireTime: %s, causing error: %s", expireTime, parseErr)
	}
	expire = uint32(expireTime64)

	if uidStr == "" || uidStr == "0" {
		err = fmt.Errorf("invalid RTM User ID: \"%s\"", uidStr)
	}

	// check if string conversion fails
	return uidStr, expire, err
}

func (s *Service) parseChatParams(c *gin.Context) (uidStr string, tokenType string, expire uint32, err error) {
	// get param values
	uidStr = c.Param("chatid")
	urlSplit := strings.Split(c.Request.URL.Path, "/")
	for i := range urlSplit {
		if urlSplit[i] == "chat" {
			tokenType = urlSplit[i+1]
			break
		}
	}
	expireTime := c.DefaultQuery("expiry", "3600")
	expireTime64, parseErr := strconv.ParseUint(expireTime, 10, 64)
	if parseErr != nil {
		// if string conversion fails return an error
		err = fmt.Errorf("failed to parse expireTime: %s, causing error: %s", expireTime, parseErr)
	}
	expire = uint32(expireTime64)
	if tokenType == "account" {
		tokenType = "userAccount"
	}
	if uidStr == "" && tokenType != "app" {
		err = fmt.Errorf("userAccount type requires chat ID")
		return uidStr, tokenType, expire, err
	}

	// check if string conversion fails
	return uidStr, tokenType, expire, err
}
