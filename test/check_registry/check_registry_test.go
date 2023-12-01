package main

import (
	"net/url"
	"testing"
	"net/http"
	"net/http/httptest"
	"io"
	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestGetProjectReplacements(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Case 1: Validated replacement urls",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectReplacements := getProjectReplacements()
			for _, projectReplacement := range projectReplacements {
				_, err := url.ParseRequestURI(projectReplacement.ReplacementRepo)
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestGetRegistries(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Case 1: Validated registries urls",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registries := getRegistries()
			for _, registry := range registries {
				_, err := url.ParseRequestURI(registry.Url)
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestGetStarterProjects(t *testing.T) {
	tests := []struct {
		name        string
		devfileUrl  string
		expectedUrl string
	}{
		{
			name:        "Case 1: Validated registries urls",
			devfileUrl:  "https://registry.devfile.io/devfiles/java-maven",
			expectedUrl: "https://github.com/odo-devfiles/springboot-ex.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			starterProjects, err := getStarterProjects(tt.devfileUrl)
			if err != nil {
				t.Error(err)
			}
			for _, starterProject := range starterProjects {
				assert.EqualValues(t, starterProject.Repo, tt.expectedUrl)
			}
		})
	}
}

func TestResponseBodyClosingStarterProjects(t *testing.T) {
	tests := []struct {
		name             string
		mockServerConfig func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "Successful response body closing",
			mockServerConfig: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("response body content"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			// Setup a mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(tc.mockServerConfig))
			defer server.Close()

			// Perform a request to the mock server
			resp, err := http.Get(server.URL)
			assert.NoError(t, err, "Unexpected error during HTTP GET")

			// Close the response body
			func() {
				defer func() {
					if err := resp.Body.Close(); err != nil {
						fmt.Printf("error closing file: %s", err)
					}
				}()

				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err, "Unexpected error reading response body")
				assert.Equal(t, "response body content", string(body), "Unexpected response body content")
			}()
		})
	}
}
