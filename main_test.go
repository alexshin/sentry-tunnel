package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTunnel(t *testing.T) {
	route := getRoot("abc", []string{"1", "2", "3"}, "http")

	wrongHostPayload := `{"sent_at":"2021-10-14T17:10:40.136Z","sdk":{"name":"sentry.javascript.browser","version":"6.13.3"},"dsn":"http://public@abc1/1"}
        {"type":"session"}
        {"sid":"751d80dc94e34cd282a2cf1fe698a8d2","init":true,"started":"2021-10-14T17:10:40.135Z","timestamp":"2021-10-14T17:10:40.135Z","status":"ok","errors":0,"attrs":{"release":"test_project@1.0"}`

	wrongProjectPayload := `{"sent_at":"2021-10-14T17:10:40.136Z","sdk":{"name":"sentry.javascript.browser","version":"6.13.3"},"dsn":"http://public@abc/5"}
        {"type":"session"}
        {"sid":"751d80dc94e34cd282a2cf1fe698a8d2","init":true,"started":"2021-10-14T17:10:40.135Z","timestamp":"2021-10-14T17:10:40.135Z","status":"ok","errors":0,"attrs":{"release":"test_project@1.0"}`

	t.Run("Non-POST methods are not supported", func(t *testing.T) {

		for _, m := range []string{http.MethodGet, http.MethodDelete, http.MethodHead, http.MethodPut} {
			request, _ := http.NewRequest(m, "/bugs", strings.NewReader("abc"))
			response := httptest.NewRecorder()

			route(response, request)

			if response.Code != 405 {
				t.Errorf("Only POST method supported")
			}
			if !strings.Contains(response.Body.String(), "supported") {
				t.Errorf("Only POST method supported")
			}
		}
	})

	t.Run("POST method is supported", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/bugs", strings.NewReader("abc"))
		response := httptest.NewRecorder()

		route(response, request)

		if response.Code == 405 {
			t.Errorf("Only POST method supported")
		}
	})

	t.Run("Request body should contain dsn", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/bugs", strings.NewReader("abc"))
		response := httptest.NewRecorder()

		route(response, request)

		if response.Code != 400 {
			t.Errorf("You must provide dsn")
		}

		if strings.Contains(response.Body.String(), "DSN") {
			t.Errorf("You must provide dsn")
		}
	})

	t.Run("Request body should contain dsn", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/bugs", strings.NewReader("abc"))
		response := httptest.NewRecorder()

		route(response, request)

		if response.Code != 400 {
			t.Errorf("You must provide dsn")
		}

		if strings.Contains(response.Body.String(), "DSN") {
			t.Errorf("You must provide dsn")
		}
	})

	t.Run("Request with wrong host", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/bugs", strings.NewReader(wrongHostPayload))
		response := httptest.NewRecorder()

		route(response, request)

		if response.Code != 400 {
			t.Errorf("Hostname is not allowed")
		}

		if strings.Contains(response.Body.String(), "DSN") {
			t.Errorf("Hostname is not allowed")
		}
	})

	t.Run("Request with wrong projectId", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/bugs", strings.NewReader(wrongProjectPayload))
		response := httptest.NewRecorder()

		route(response, request)

		if response.Code != 400 {
			t.Errorf("Project is not allowed")
		}

		if strings.Contains(response.Body.String(), "DSN") {
			t.Errorf("Project is not allowed")
		}
	})
}
