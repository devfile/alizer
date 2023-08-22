package recognizer

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devfile/alizer/pkg/apis/model"
	"github.com/stretchr/testify/assert"
)

func TestDetectComponentsWithoutPortDetection(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		components []model.Component
		want       bool
	}{
		{
			name: "Case 1: Func successful",
			path: "/testPath",
			components: []model.Component{{
				Name: "test_name",
				Path: "test/path",
				Languages: []model.Language{{
					Name:                    "lang",
					Aliases:                 []string{"alias"},
					Weight:                  0.59,
					Frameworks:              []string{"framework"},
					Tools:                   []string{"tool"},
					CanBeComponent:          true,
					CanBeContainerComponent: true,
				}},
				Ports: []int{12345},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock detectComponentsWithPathAndPortStartegy
			detectComponentsWithPathAndPortStartegy = func(path string, portDetectionStrategy []model.PortDetectionAlgorithm, ctx *context.Context) ([]model.Component, error) {
				return tt.components, nil
			}
			result, err := DetectComponentsWithoutPortDetection("somePath")
			if err != nil {
				t.Errorf("Error: %t", err)
			}
			assert.EqualValues(t, tt.components, tt.components, result)
		})
	}
}

func Test_isAnyComponentInDirectPath(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		components []model.Component
		want       bool
	}{
		{
			name: "Case 1: path should match",
			path: "/alizer/resources/projects/ocparcade/arkanoid/",
			components: []model.Component{{
				Name:      "",
				Path:      "/alizer/resources/projects/ocparcade/arkanoid/",
				Languages: nil,
				Ports:     nil,
			}},
			want: true,
		},
		{
			name: "Case 2: path should match",
			path: "/alizer/resources/projects/ocparcade/arkanoid/arkanoid/",
			components: []model.Component{{
				Name:      "",
				Path:      "/alizer/resources/projects/ocparcade/arkanoid/arkanoid/",
				Languages: nil,
				Ports:     nil,
			}},
			want: true,
		},
		{
			name: "Case 3: path should not match",
			path: "/alizer/resources/projects/ocparcade/arkanoid/",
			components: []model.Component{{
				Name:      "",
				Path:      "/alizer/resources/projects/ocparcade/arkanoid/arkanoid/",
				Languages: nil,
				Ports:     nil,
			}},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAnyComponentInDirectPath(tt.path, tt.components)
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("Got: %t, want: %t", result, tt.want)
			}
		})
	}
}

func Test_detectComponentsWithSettings(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		expectedComponents []model.Component
		expectingError     bool
	}{
		{
			name: "Case 1: detect components",
			path: "../../../resources/projects/beego",
			expectedComponents: []model.Component{
				{
					Name: "beego",
					Path: "../../../resources/projects/beego",
				},
			},
			expectingError: false,
		},
		{
			name:               "Case 2: invalid path",
			path:               "../../../resources/projects/notexisting",
			expectedComponents: []model.Component{},
			expectingError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detectionSettings, err := getDetectionSettings(tt.path)
			if err != nil {
				t.Errorf("test failed: Couldn't locate current working directory")
			}
			ctx := context.Background()
			result, err := detectComponentsWithSettings(detectionSettings, &ctx)
			if tt.expectingError {
				if err == nil {
					t.Errorf("test Failed: Was expecting path not found")
				}
				assert.EqualValues(t, 0, len(result))
			} else {
				if len(result) != 1 {
					t.Errorf("expected 1 component for %s dir", tt.path)
				}
				expectedPath, err := getAbsolutePath(tt.expectedComponents[0].Path)
				if err != nil {
					t.Errorf("test failed: %s", err)
				}
				assert.EqualValues(t, tt.expectedComponents[0].Name, result[0].Name)
				assert.EqualValues(t, expectedPath, result[0].Path)
			}
		})
	}
}

func TestDetectComponentsWithSettings(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		expectedComponents []model.Component
		expectingError     bool
	}{
		{
			name: "Case 1: detect components",
			path: "../../../resources/projects/beego",
			expectedComponents: []model.Component{
				{
					Name: "beego",
					Path: "../../../resources/projects/beego",
				},
			},
			expectingError: false,
		},
		{
			name:               "Case 2: invalid path",
			path:               "../../../resources/projects/notexisting",
			expectedComponents: []model.Component{},
			expectingError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detectionSettings, err := getDetectionSettings(tt.path)
			if err != nil {
				t.Errorf("test failed: Couldn't locate current working directory")
			}
			result, err := DetectComponentsWithSettings(detectionSettings)
			if tt.expectingError {
				if err == nil {
					t.Errorf("test Failed: Was expecting path not found")
				}
				assert.EqualValues(t, 0, len(result))
			} else {
				if len(result) != 1 {
					t.Errorf("expected 1 component for %s dir", tt.path)
				}
				expectedPath, err := getAbsolutePath(tt.expectedComponents[0].Path)
				if err != nil {
					t.Errorf("test failed: %s", err)
				}
				assert.EqualValues(t, tt.expectedComponents[0].Name, result[0].Name)
				assert.EqualValues(t, expectedPath, result[0].Path)
			}
		})
	}
}

func Test_detectComponentsInRootWithSettings(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		expectedComponents []model.Component
		expectingError     bool
	}{
		{
			name: "Case 1: detect components",
			path: "../../../resources/projects/beego",
			expectedComponents: []model.Component{
				{
					Name: "beego",
					Path: "../../../resources/projects/beego",
				},
			},
			expectingError: false,
		},
		{
			name:               "Case 2: invalid path",
			path:               "../../../resources/projects/notexisting",
			expectedComponents: []model.Component{},
			expectingError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detectionSettings, err := getDetectionSettings(tt.path)
			if err != nil {
				t.Errorf("test failed: Couldn't locate current working directory")
			}
			ctx := context.Background()
			result, err := detectComponentsInRootWithSettings(detectionSettings, &ctx)
			if tt.expectingError {
				if err == nil {
					t.Errorf("test Failed: Was expecting path not found")
				}
				assert.EqualValues(t, 0, len(result))
			} else {
				if len(result) != 1 {
					t.Errorf("expected 1 component for %s dir", tt.path)
				}
				expectedPath, err := getAbsolutePath(tt.expectedComponents[0].Path)
				if err != nil {
					t.Errorf("test failed: %s", err)
				}
				assert.EqualValues(t, tt.expectedComponents[0].Name, result[0].Name)
				assert.EqualValues(t, expectedPath, result[0].Path)
			}
		})
	}
}

func TestDetectComponentsInRootWithSettings(t *testing.T) {
	tests := []struct {
		name               string
		path               string
		expectedComponents []model.Component
		expectingError     bool
	}{
		{
			name: "Case 1: detect components",
			path: "../../../resources/projects/beego",
			expectedComponents: []model.Component{
				{
					Name: "beego",
					Path: "../../../resources/projects/beego",
				},
			},
			expectingError: false,
		},
		{
			name:               "Case 2: invalid path",
			path:               "../../../resources/projects/notexisting",
			expectedComponents: []model.Component{},
			expectingError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detectionSettings, err := getDetectionSettings(tt.path)
			if err != nil {
				t.Errorf("test failed: Couldn't locate current working directory")
			}
			result, err := DetectComponentsInRootWithSettings(detectionSettings)
			if tt.expectingError {
				if err == nil {
					t.Errorf("test Failed: Was expecting path not found")
				}
				assert.EqualValues(t, 0, len(result))
			} else {
				if len(result) != 1 {
					t.Errorf("expected 1 component for %s dir", tt.path)
				}
				expectedPath, err := getAbsolutePath(tt.expectedComponents[0].Path)
				if err != nil {
					t.Errorf("test failed: %s", err)
				}
				assert.EqualValues(t, tt.expectedComponents[0].Name, result[0].Name)
				assert.EqualValues(t, expectedPath, result[0].Path)
			}
		})
	}
}

func getDetectionSettings(path string) (model.DetectionSettings, error) {
	absPath, err := getAbsolutePath(path)
	if err != nil {
		return model.DetectionSettings{}, err
	}
	return model.DetectionSettings{
		BasePath:              absPath,
		PortDetectionStrategy: []model.PortDetectionAlgorithm{model.Compose, model.DockerFile, model.Source},
	}, nil
}

func getAbsolutePath(path string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return strings.Join([]string{filepath.Clean(filepath.Join(pwd, path)), "/"}, ""), nil
}
