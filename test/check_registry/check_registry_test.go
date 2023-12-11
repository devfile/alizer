package main

import (
	"net/url"
	"testing"

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