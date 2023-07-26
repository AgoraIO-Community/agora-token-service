package service

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/AgoraIO-Community/go-tokenbuilder/chatTokenBuilder"
	rtctokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	rtmtokenbuilder2 "github.com/AgoraIO-Community/go-tokenbuilder/rtmtokenbuilder"
	"github.com/gin-gonic/gin"
)

// getToken handles the HTTP request to generate a token based on the provided tokenType.
// It checks the tokenType from the query parameters and calls the appropriate token generation method.
// The generated token is sent as a JSON response to the client.
//
// Parameters:
//   - c: *gin.Context - The Gin context representing the HTTP request and response.
//
// Behavior:
//  1. Retrieves the tokenType from the query parameters, defaulting to "rtc" if not provided.
//  2. Uses a switch statement to handle different tokenType cases:
//     - "rtm": Calls the RtmToken method to generate the RTM token and sends it as a JSON response.
//     - "chat": Calls the ChatToken method to generate the chat token and sends it as a JSON response.
//     - Default: Calls the RtcToken method to generate the RTC token and sends it as a JSON response.
//
// Notes:
//   - The actual token generation methods (RtmToken, ChatToken, and RtcToken) are part of the Service struct.
//   - The generated token is sent as a JSON response with appropriate HTTP status codes.
//
// Example usage:
//
//	router.GET("/getToken", service.getToken)
func (s *Service) getToken(c *gin.Context) {
	// Parse the request body into a TokenRequest struct
	var tokenReq TokenRequest
	if err := c.ShouldBindJSON(&tokenReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var token string
	var tokenErr error

	switch tokenReq.TokenType {
	case "rtc":
		token, tokenErr = s.genRtcToken(tokenReq)
	case "rtm":
		token, tokenErr = s.genRtmToken(tokenReq)
	case "chat":
		token, tokenErr = s.genChatToken(tokenReq)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported tokenType"})
		return
	}
	if tokenErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": tokenErr.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// TokenRequest is a struct representing the JSON payload structure for token generation requests.
// It contains fields necessary for generating different types of tokens (RTC, RTM, or chat) based on the "TokenType".
// The "Channel", "RtcRole", "Uid", and "ExpirationSeconds" fields are used for specific token types.
//
// TokenType options: "rtc" for RTC token, "rtm" for RTM token, and "chat" for chat token.
type TokenRequest struct {
	TokenType         string `json:"tokenType"`         // The token type: "rtc", "rtm", or "chat"
	Channel           string `json:"channel,omitempty"` // The channel name (used for RTC and RTM tokens)
	RtcRole           string `json:"role,omitempty"`    // The role of the user for RTC tokens (publisher or subscriber)
	Uid               string `json:"uid,omitempty"`     // The user ID or account (used for RTC, RTM, and some chat tokens)
	ExpirationSeconds int    `json:"expire,omitempty"`  // The token expiration time in seconds (used for all token types)
}

// genRtcToken generates an RTC token based on the provided TokenRequest and returns it.
//
// Parameters:
//   - tokenRequest: TokenRequest - The TokenRequest struct containing the required information for RTC token generation.
//
// Returns:
//   - string: The generated RTC token.
//   - error: An error if there are any issues during token generation or validation.
//
// Behavior:
//  1. Validates the required fields in the TokenRequest (channel and UID).
//  2. Sets a default expiration time of 3600 seconds (1 hour) if not provided in the request.
//  3. Determines the user's role (publisher or subscriber) based on the "Role" field in the request.
//  4. Generates the RTC token using the rtctokenbuilder2 package.
//
// Notes:
//   - The rtctokenbuilder2 package is used for generating RTC tokens.
//   - The "Role" field can be "publisher" or "subscriber"; other values are considered invalid.
//
// Example usage:
//
//	tokenReq := TokenRequest{
//	    TokenType:  "rtc",
//	    Channel:    "my_channel",
//	    Uid:        "user123",
//	    Role:       "publisher",
//	    ExpirationSeconds: 3600,
//	}
//	token, err := service.genRtcToken(tokenReq)
func (s *Service) genRtcToken(tokenRequest TokenRequest) (string, error) {
	if tokenRequest.Channel == "" {
		return "", errors.New("invalid: missing channel name")
	}
	if tokenRequest.Uid == "" {
		return "", errors.New("invalid: missing user ID or account")
	}

	var userRole rtctokenbuilder2.Role
	if tokenRequest.RtcRole == "publisher" {
		userRole = rtctokenbuilder2.RolePublisher
	} else {
		userRole = rtctokenbuilder2.RoleSubscriber
	}

	if tokenRequest.ExpirationSeconds == 0 {
		tokenRequest.ExpirationSeconds = 3600
	}

	uid64, parseErr := strconv.ParseUint(tokenRequest.Uid, 10, 64)
	if parseErr != nil {
		return rtctokenbuilder2.BuildTokenWithAccount(
			s.appID, s.appCertificate, tokenRequest.Channel,
			tokenRequest.Uid, userRole, uint32(tokenRequest.ExpirationSeconds),
		)
	}

	return rtctokenbuilder2.BuildTokenWithUid(
		s.appID, s.appCertificate, tokenRequest.Channel,
		uint32(uid64), userRole, uint32(tokenRequest.ExpirationSeconds),
	)
}

// genRtmToken generates an RTM (Real-Time Messaging) token based on the provided TokenRequest and returns it.
//
// Parameters:
//   - tokenRequest: TokenRequest - The TokenRequest struct containing the required information for RTM token generation.
//
// Returns:
//   - string: The generated RTM token.
//   - error: An error if there are any issues during token generation or validation.
//
// Behavior:
//  1. Validates the required field in the TokenRequest (UID).
//  2. Sets a default expiration time of 3600 seconds (1 hour) if not provided in the request.
//  3. Generates the RTM token using the rtmtokenbuilder2 package.
//
// Notes:
//   - The rtmtokenbuilder2 package is used for generating RTM tokens.
//   - The "UID" field in TokenRequest is mandatory for RTM token generation.
//
// Example usage:
//
//	tokenReq := TokenRequest{
//	    TokenType:  "rtm",
//	    Uid:        "user123",
//	    ExpirationSeconds: 3600,
//	}
//	token, err := service.genRtmToken(tokenReq)
func (s *Service) genRtmToken(tokenRequest TokenRequest) (string, error) {
	if tokenRequest.Uid == "" {
		return "", errors.New("invalid: missing user ID or account")
	}
	if tokenRequest.ExpirationSeconds == 0 {
		tokenRequest.ExpirationSeconds = 3600
	}

	return rtmtokenbuilder2.BuildToken(
		s.appID, s.appCertificate,
		tokenRequest.Uid,
		uint32(tokenRequest.ExpirationSeconds),
		tokenRequest.Channel,
	)
}

// genChatToken generates a chat token based on the provided TokenRequest and returns it.
//
// Parameters:
//   - tokenRequest: TokenRequest - The TokenRequest struct containing the required information for chat token generation.
//
// Returns:
//   - string: The generated chat token.
//   - error: An error if there are any issues during token generation or validation.
//
// Behavior:
//  1. Sets a default expiration time of 3600 seconds (1 hour) if not provided in the request.
//  2. Determines whether to generate a chat app token or a chat user token based on the "UID" field in the request.
//  3. Generates the chat token using the chatTokenBuilder package.
//
// Notes:
//   - The chatTokenBuilder package is used for generating chat tokens.
//   - If the "UID" field is empty, a chat app token is generated; otherwise, a chat user token is generated.
//
// Example usage:
//
//	// Generate a chat app token
//	tokenReq := TokenRequest{
//	    TokenType:  "chat",
//	    ExpirationSeconds: 3600,
//	}
//	token, err := service.genChatToken(tokenReq)
//
//	// Generate a chat user token
//	tokenReq := TokenRequest{
//	    TokenType:  "chat",
//	    Uid:        "user123",
//	    ExpirationSeconds: 3600,
//	}
//	token, err := service.genChatToken(tokenReq)
func (s *Service) genChatToken(tokenRequest TokenRequest) (string, error) {
	if tokenRequest.ExpirationSeconds == 0 {
		tokenRequest.ExpirationSeconds = 3600
	}

	var chatToken string
	var tokenErr error

	if tokenRequest.Uid == "" {
		chatToken, tokenErr = chatTokenBuilder.BuildChatAppToken(
			s.appID, s.appCertificate, uint32(tokenRequest.ExpirationSeconds),
		)
	} else {
		chatToken, tokenErr = chatTokenBuilder.BuildChatUserToken(
			s.appID, s.appCertificate,
			tokenRequest.Uid,
			uint32(tokenRequest.ExpirationSeconds),
		)

	}

	return chatToken, tokenErr
}
