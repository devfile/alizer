package recognizer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/devfile/alizer/pkg/apis/model"
	"github.com/devfile/alizer/pkg/apis/recognizer"
)

func updateContent(filePath string, data []byte) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(){
		if err := f.Close(); err != nil {
			fmt.Printf("error closing file: %s", err)
		}
	}()
	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}

func getComponentsFromTestProject(t *testing.T, project string) []model.Component {
	testingProjectPath := getTestProjectPath(project)
	return getComponentsFromProjectInner(t, testingProjectPath)
}

func getComponentsFromProjectInner(t *testing.T, testingProjectPath string) []model.Component {
	components, err := recognizer.DetectComponents(testingProjectPath)
	if err != nil {
		t.Error(err)
	}

	return components
}

func getComponentsFromFiles(t *testing.T, files []string, settings model.DetectionSettings) []model.Component {
	ctx := context.Background()
	return recognizer.DetectComponentsFromFilesList(files, settings, &ctx)
}

func isComponentsInProject(t *testing.T, project string, expectedNumber int, expectedLanguage string, expectedProjectName string) {
	components := getComponentsFromTestProject(t, project)
	verifyComponents(t, components, expectedNumber, expectedLanguage, expectedProjectName)
}

func verifyComponents(t *testing.T, components []model.Component, expectedNumber int, expectedLanguage string, expectedProjectName string) {
	hasComponents := len(components) == expectedNumber
	if hasComponents {
		isExpectedComponent := strings.EqualFold(expectedLanguage, components[0].Languages[0].Name)
		if !isExpectedComponent {
			t.Errorf("Project does not use " + expectedLanguage + " language")
		}
		if expectedProjectName != "" {
			isExpectedProjectName := strings.EqualFold(expectedProjectName, components[0].Name)
			if !isExpectedProjectName {
				t.Errorf("Main component has a different project name. Expected " + expectedProjectName + " but it was " + components[0].Name)
			}
		}
	} else {
		t.Errorf("Expected " + strconv.Itoa(expectedNumber) + " of components but it was " + strconv.Itoa(len(components)))
	}
}

func testPortDetectionInProject(t *testing.T, project string, ports []int) {
	components := getComponentsFromTestProject(t, project)
	if len(components) == 0 {
		t.Errorf("No component detected")
	}

	portsDetected := components[0].Ports
	if len(portsDetected) != len(ports) {
		t.Errorf("Length of found ports and expected ports is different")
	}

	found := false
	for _, port := range ports {
		for _, portDetected := range portsDetected {
			if port == portDetected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Port " + strconv.Itoa(port) + " have not been detected")
		}
		found = false
	}
}

func getTestProjectPath(folder string) string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Join(basepath, "..", "..", "resources/projects", folder)
}
