package recognizer

import (
	"context"
	"reflect"
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
