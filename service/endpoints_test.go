package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRtcValidAndInvalid(t *testing.T) {

	tests := []UrlCodePair{
		{"/rtc/fsda/publisher/uid/0/?expiry=600", http.StatusOK},
		{"/rtc/fsda/publisher/uid//?expiry=600", http.StatusOK},
		{"/rtc/fsda/publisher/uid/test/?expiry=600", http.StatusBadRequest},
		{"/rtc/fsda/publisher/uid/0/?expiry=failing", http.StatusBadRequest},
	}
	for _, httpTest := range tests {
		testApi, err := http.NewRequest(http.MethodGet, httpTest.url, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp := httptest.NewRecorder()
		// Call the endpoint
		testService.Server.Handler.ServeHTTP(resp, testApi)
		assert.Equal(t, httpTest.code, resp.Code, resp.Body)
	}
}

func TestRtmValidAndInvalid(t *testing.T) {

	// Create a new gin engine for testing
	reqValid, err := http.NewRequest(http.MethodGet, "/rtm/username/?expiry=600", nil)
	if err != nil {
		t.Fatal(err)
	}
	reqInvalid, err := http.NewRequest(http.MethodGet, "/rtm/0/?expiry=600", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder to inspect the response
	resp := httptest.NewRecorder()

	// Call the endpoint
	testService.Server.Handler.ServeHTTP(resp, reqValid)

	assert.Equal(t, http.StatusOK, resp.Code)

	resp = httptest.NewRecorder()

	// Call the endpoint
	testService.Server.Handler.ServeHTTP(resp, reqInvalid)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

type UrlCodePair struct {
	url  string
	code int
}

func TestChatValidAndInvalid(t *testing.T) {

	tests := []UrlCodePair{
		{"/chat/app/", http.StatusOK},
		{"/chat/account/username/", http.StatusOK},
		{"/chat/account/", http.StatusNotFound},
		{"/chat/invalid/", http.StatusNotFound},
		{"/chat/account/username/?expiry=600", http.StatusOK},
		{"/chat/account/username/?expiry=fail", http.StatusBadRequest},
	}
	for _, httpTest := range tests {
		testApi, err := http.NewRequest(http.MethodGet, httpTest.url, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp := httptest.NewRecorder()
		// Call the endpoint
		testService.Server.Handler.ServeHTTP(resp, testApi)
		assert.Equal(t, httpTest.code, resp.Code, resp.Body)
	}
}

func TestRteValidAndInvalid(t *testing.T) {

	tests := []UrlCodePair{
		{"/rte/channelName/publisher/uid/0/rtmid/?expiry=600", http.StatusOK},
		{"/rte/channelName/publisher/uid/2345/?expiry=600", http.StatusOK},
		{"/rte/channelName/subscriber/uid/0/rtmid/?expiry=600", http.StatusOK},
		{"/rte/channelName/subscriber/uid/2345/?expiry=600", http.StatusOK},
		{"/rte/channelName/publisher/uid/0/?expiry=600", http.StatusBadRequest},
		{"/rte/channelName/publisher/uid/2345/?expiry=failing", http.StatusBadRequest},
	}
	for _, httpTest := range tests {
		testApi, err := http.NewRequest(http.MethodGet, httpTest.url, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp := httptest.NewRecorder()
		// Call the endpoint
		testService.Server.Handler.ServeHTTP(resp, testApi)
		assert.Equal(t, httpTest.code, resp.Code, resp.Body)
	}
}

func TestTokentypes(t *testing.T) {

	// Create a new gin engine for testing
	reqValidUid, err := http.NewRequest(http.MethodGet, "/rtc/fsda/publisher/uid/0/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a new gin engine for testing
	reqValidUserAcc, err := http.NewRequest(http.MethodGet, "/rtc/fsda/publisher/userAccount/0/", nil)
	if err != nil {
		t.Fatal(err)
	}
	reqInvalid, err := http.NewRequest(http.MethodGet, "/rtc/fsda/publisher/nonsense/0/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder to inspect the response
	resp := httptest.NewRecorder()

	// Call the endpoint
	testService.Server.Handler.ServeHTTP(resp, reqValidUid)

	assert.Equal(t, http.StatusOK, resp.Code)

	resp = httptest.NewRecorder()

	// Call the endpoint
	testService.Server.Handler.ServeHTTP(resp, reqValidUserAcc)

	assert.Equal(t, http.StatusOK, resp.Code)

	resp = httptest.NewRecorder()

	// Call the endpoint
	testService.Server.Handler.ServeHTTP(resp, reqInvalid)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
