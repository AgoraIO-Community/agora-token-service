package service

import (
	"fmt"
	"log"
	"strconv"

	rtctokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
)

func (s *Service) generateRtcToken(channelName, uidStr, tokenType string, role rtctokenbuilder2.Role, expireTimestamp uint32) (rtcToken string, err error) {

	if tokenType == "userAccount" {
		log.Printf("Building Token with userAccount: %s\n", uidStr)
		rtcToken, err = rtctokenbuilder2.BuildTokenWithAccount(s.appID, s.appCertificate, channelName, uidStr, role, expireTimestamp)
		return rtcToken, err

	} else if tokenType == "uid" {
		uid64, parseErr := strconv.ParseUint(uidStr, 10, 64)
		// check if conversion fails
		if parseErr != nil {
			err = fmt.Errorf("failed to parse uidStr: %s, to uint causing error: %s", uidStr, parseErr)
			return "", err
		}

		uid := uint32(uid64) // convert uid from uint64 to uint 32
		log.Printf("Building Token with uid: %d\n", uid)
		rtcToken, err = rtctokenbuilder2.BuildTokenWithUid(s.appID, s.appCertificate, channelName, uid, role, expireTimestamp)
		return rtcToken, err
	} else {
		err = fmt.Errorf("failed to generate RTC token for Unknown Tokentype: %s", tokenType)
		log.Println(err)
		return "", err
	}
}
