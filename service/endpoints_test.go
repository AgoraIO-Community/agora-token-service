package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRtcValidAndInvalid(t *testing.T) {

	// Create a new gin engine for testing
	reqValid, err := http.NewRequest(http.MethodGet, "/rtc/fsda/publisher/uid/0/?expiry=600", nil)
	if err != nil {
		t.Fatal(err)
	}
	reqInvalid, err := http.NewRequest(http.MethodGet, "/rtc/fsda/publisher/uid/test/?expiry=600", nil)
	if err != nil {
		t.Fatal(err)
	}
	reqInvalidExp, err := http.NewRequest(http.MethodGet, "/rtc/fsda/publisher/uid/0/?expiry=failing", nil)
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

	resp = httptest.NewRecorder()

	// Call the endpoint
	testService.Server.Handler.ServeHTTP(resp, reqInvalidExp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
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

func TestRteValidAndInvalid(t *testing.T) {

	// Create a new gin engine for testing
	reqValid, err := http.NewRequest(http.MethodGet, "/rte/channelName/publisher/uid/0/rtmid/?expiry=600", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a new gin engine for testing
	reqValid2, err := http.NewRequest(http.MethodGet, "/rte/channelName/publisher/uid/2345/?expiry=600", nil)
	if err != nil {
		t.Fatal(err)
	}
	reqInvalid, err := http.NewRequest(http.MethodGet, "/rte/channelName/publisher/uid/0/?expiry=600", nil)
	if err != nil {
		t.Fatal(err)
	}
	reqInvalidExp, err := http.NewRequest(http.MethodGet, "/rte/channelName/publisher/uid/2345/?expiry=failing", nil)
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
	testService.Server.Handler.ServeHTTP(resp, reqValid2)

	assert.Equal(t, http.StatusOK, resp.Code)

	resp = httptest.NewRecorder()

	// Call the endpoint
	testService.Server.Handler.ServeHTTP(resp, reqInvalid)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	resp = httptest.NewRecorder()

	// Call the endpoint
	testService.Server.Handler.ServeHTTP(resp, reqInvalidExp)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
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
