package component

import (
	"testing"

	"github.com/devfile/alizer/pkg/apis/model"
	"github.com/stretchr/testify/assert"
)

func TestGetPortStrategy(t *testing.T) {

	tests := []struct {
		name                 string
		expectedResult       []model.PortDetectionAlgorithm
		noPortDetectionValue bool
	}{
		{
			name:                 "Case 1: default port detection",
			expectedResult:       []model.PortDetectionAlgorithm{model.DockerFile, model.Compose, model.Source},
			noPortDetectionValue: false,
		},
		{
			name:                 "Case 2: no port detection active",
			expectedResult:       []model.PortDetectionAlgorithm{},
			noPortDetectionValue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			noPortDetection = tt.noPortDetectionValue
			result := getPortDetectionStrategy()
			assert.EqualValues(t, tt.expectedResult, result)
		})
	}
}
