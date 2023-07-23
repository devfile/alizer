package recognizer

import (
	"github.com/devfile/alizer/pkg/apis/model"
	"reflect"
	"testing"
)

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
