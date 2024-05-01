package service

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Define a struct to match the JSON response structure
type AuthResponse struct {
	Token        string `json:"token"`
	FitnessAppID int    `json:"fitness_app_id"`
	ExpiresAt    string `json:"expires_at"`
	UpdatedAt    string `json:"updated_at"`
	CreatedAt    string `json:"created_at"`
}

type AuthRequest struct {
	AuthenticationToken struct {
		FitnessAppID int `json:"fitness_app_id"`
	} `json:"authentication_token"`
}

func authMiddleware(authorizationServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the security header
		securityHeader := c.GetHeader("WF-USER-TOKEN")
		if securityHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Security header is required"})
			return
		}

		// Prepare the JSON payload
		payload := AuthRequest{}
		payload.AuthenticationToken.FitnessAppID = 1 // Set statically or from some logic

		// Marshal the payload into JSON
		jsonData, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error preparing request"})
			return
		}

		// Create and send the request
		client := &http.Client{Timeout: time.Second * 10}
		req, err := http.NewRequest("POST", authorizationServiceURL, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("WF-USER-TOKEN", securityHeader)
		if err != nil {
			log.Printf("Failed to create request: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate security header"})
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to send request to authorization service: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate security header"})
			return
		}
		defer resp.Body.Close()

		// Read and parse the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate security header"})
			return
		}

		// Decode the JSON response
		var authResponse AuthResponse
		if err = json.Unmarshal(body, &authResponse); err != nil {
			log.Printf("Error decoding JSON response: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error decoding authorization response. Make sure WF-USER-TOKEN is set correctly"})
			return
		}

		// Check if the token is present
		if authResponse.Token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - No token returned"})
			return
		}

		// Store the token in Gin context and proceed
		c.Set("userToken", authResponse.Token)
		c.Next()
	}
}
