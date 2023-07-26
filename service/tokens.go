package service

import (
	"fmt"
	"log"
	"strconv"

	"github.com/AgoraIO-Community/go-tokenbuilder/chatTokenBuilder"
	rtctokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
)

// generateRtcToken generates an RTC token for the video conferencing application based on the provided parameters.
//
// Parameters:
//   - channelName: string - The name of the video conferencing channel.
//   - uidStr: string - The user ID for the RTC token, represented as a string.
//   - tokenType: string - The type of RTC token. Can be "userAccount" or "uid".
//   - role: rtctokenbuilder2.Role - The role of the user (Publisher or Subscriber) in the video conferencing.
//   - expireDelta: uint32 - The duration of the token's validity in seconds.
//
// Returns:
//   - rtcToken: string - The generated RTC token as a string.
//   - err: error - Any error that occurred during token generation. Nil if token generation was successful.
//
// Behavior:
//  1. Checks the tokenType to determine whether to build the token using the userAccount or uid.
//  2. If the tokenType is "userAccount", builds the RTC token using the user account (uidStr).
//  3. If the tokenType is "uid", parses uidStr to an unsigned 64-bit integer and converts it to uint32.
//  4. Builds the RTC token using the numeric user ID (uid) and the provided role and expireDelta.
//  5. If the tokenType is neither "userAccount" nor "uid", returns an error indicating the unknown tokenType.
//
// Notes:
//   - The `rtctokenbuilder2.Role` type represents the role of a user in the video conferencing token.
//     It might be an enumeration or constant representing different roles (e.g., Publisher and Subscriber).
//
// Example usage:
//
//	rtcToken, err := generateRtcToken("channel123", "user123", "userAccount", rtctokenbuilder2.RolePublisher, 3600)
func (s *Service) generateRtcToken(channelName, uidStr, tokenType string, role rtctokenbuilder2.Role, expireDelta uint32) (rtcToken string, err error) {

	if tokenType == "userAccount" {
		log.Printf("Building Token for userAccount: %s\n", uidStr)
		rtcToken, err = rtctokenbuilder2.BuildTokenWithAccount(s.appID, s.appCertificate, channelName, uidStr, role, expireDelta)
		return rtcToken, err
	} else if tokenType == "uid" {
		uid64, parseErr := strconv.ParseUint(uidStr, 10, 64)
		// check if conversion fails
		if parseErr != nil {
			err = fmt.Errorf("failed to parse uidStr: %s, to uint causing error: %s", uidStr, parseErr)
			return "", err
		}

		uid := uint32(uid64) // convert uid from uint64 to uint 32
		log.Printf("Building Token for uid: %d\n", uid)
		rtcToken, err = rtctokenbuilder2.BuildTokenWithUid(s.appID, s.appCertificate, channelName, uid, role, expireDelta)
		return rtcToken, err
	} else {
		err = fmt.Errorf("failed to generate RTC token for Unknown Tokentype: %s", tokenType)
		log.Println(err)
		return "", err
	}
}

func (s *Service) generateChatToken(uidStr string, tokenType string, expireTimestamp uint32) (chatToken string, err error) {

	if tokenType == "userAccount" {
		log.Printf("Building Token with userAccount: %s\n", uidStr)
		chatToken, err = chatTokenBuilder.BuildChatUserToken(s.appID, s.appCertificate, uidStr, expireTimestamp)
		return chatToken, err

	} else if tokenType == "app" {
		chatToken, err = chatTokenBuilder.BuildChatAppToken(s.appID, s.appCertificate, expireTimestamp)
		return chatToken, err
	} else {
		err = fmt.Errorf("failed to generate Chat token for Unknown token type: %s", tokenType)
		log.Println(err)
		return "", err
	}
}
