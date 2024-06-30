package http_server

import (
	"github.com/prusiksoftware/monorepo/obscurio/analytics"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestServer_healthzReady(t *testing.T) {

	t.Run("default", func(t *testing.T) {
		s := NewHTTPServer(nil)
		request := httptest.NewRequest("GET", "http://localhost/healthz/ready", nil)
		response := httptest.NewRecorder()

		s.healthzReady(response, request)

		require.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("not ready", func(t *testing.T) {
		s := NewHTTPServer(nil)
		s.SetReady(false)
		request := httptest.NewRequest("GET", "http://localhost/healthz/ready", nil)
		response := httptest.NewRecorder()

		s.healthzReady(response, request)

		require.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("ready", func(t *testing.T) {
		s := NewHTTPServer(nil)
		s.SetReady(true)
		request := httptest.NewRequest("GET", "http://localhost/healthz/ready", nil)
		response := httptest.NewRecorder()

		s.healthzReady(response, request)

		require.Equal(t, http.StatusOK, response.Code)
	})
}

func TestServer_healthzLive(t *testing.T) {

	t.Run("default", func(t *testing.T) {
		s := NewHTTPServer(nil)
		request := httptest.NewRequest("GET", "http://localhost/healthz/live", nil)
		response := httptest.NewRecorder()

		s.healthzLive(response, request)

		require.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("not live", func(t *testing.T) {
		s := NewHTTPServer(nil)
		s.SetLive(false)
		request := httptest.NewRequest("GET", "http://localhost/healthz/live", nil)
		response := httptest.NewRecorder()

		s.healthzLive(response, request)

		require.Equal(t, http.StatusInternalServerError, response.Code)
	})

	t.Run("live", func(t *testing.T) {
		s := NewHTTPServer(nil)
		s.SetLive(true)
		request := httptest.NewRequest("GET", "http://localhost/healthz/live", nil)
		response := httptest.NewRecorder()

		s.healthzLive(response, request)

		require.Equal(t, http.StatusOK, response.Code)
	})
}

func TestServer_debug(t *testing.T) {
	at := analytics.New(60)
	at.TrackQuery(
		"profile",
		"original",
		"default",
		map[analytics.DurationType]time.Duration{
			"query": time.Second,
		})
	durations := make(map[analytics.DurationType]time.Duration)
	durations["query"] = time.Second
	at.TrackQuery("profile", "original", "default", durations)

	s := NewHTTPServer(at)
	request := httptest.NewRequest("GET", "http://localhost/debug", nil)
	response := httptest.NewRecorder()

	s.debug(response, request)

	require.Equal(t, http.StatusOK, response.Code)
}
