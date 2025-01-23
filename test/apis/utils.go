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

	framework "github.com/devfile/alizer/pkg/apis/enricher/framework/java"
	"github.com/devfile/alizer/pkg/apis/model"
	"github.com/devfile/alizer/pkg/apis/recognizer"
	"github.com/stretchr/testify/assert"
)

func updateContent(filePath string, data []byte) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
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
			t.Errorf("Project does not use %v language ", expectedLanguage)
		}
		if expectedProjectName != "" {
			isExpectedProjectName := strings.EqualFold(expectedProjectName, components[0].Name)
			if !isExpectedProjectName {
				t.Errorf("Main component has a different project name. Expected %v but it was %v", expectedProjectName, components[0].Name)
			}
		}
	} else {
		t.Errorf("Expected %v of components but it was %v", strconv.Itoa(expectedNumber), strconv.Itoa(len(components)))
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
			t.Errorf("Port %v have not been detected", strconv.Itoa(port))
		}
		found = false
	}
}

func getTestProjectPath(folder string) string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Join(basepath, "..", "..", "resources/projects", folder)
}

func testOpenLibertyDetector_DoPortsDetection(t *testing.T, projectType string, expectedPorts []int) {
	tempDir, err := os.MkdirTemp("", projectType+"-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	hardcodedXML := `
	<server>
		<httpEndpoint host="*" httpPort="1234" httpsPort="1235" id="defaultHttpEndpoint"/>
	</server>
	`
	writeTestFile(t, tempDir, "hardcoded_server.xml", hardcodedXML)

	variableXML := `
	<server>
		<variable name="default.http.port" defaultValue="9080"/>
		<variable name="default.https.port" defaultValue="9443"/>
		<httpEndpoint host="*" httpPort="${default.http.port}" httpsPort="${default.https.port}" id="defaultHttpEndpoint"/>
	</server>
	`
	writeTestFile(t, tempDir, "variable_server.xml", variableXML)

	mixedXML := `
	<server>
		<variable name="default.http.port" defaultValue="9080"/>
		<httpEndpoint host="*" httpPort="${default.http.port}" httpsPort="1235" id="defaultHttpEndpoint"/>
	</server>
	`
	writeTestFile(t, tempDir, "mixed_server.xml", mixedXML)

	emptyVarXML := `
	<server>
		<variable name="default.http.port" defaultValue=""/>
		<variable name="default.https.port" defaultValue="9443"/>
		<httpEndpoint host="*" httpPort="${default.http.port}" httpsPort="${default.https.port}" id="defaultHttpEndpoint"/>
	</server>
	`
	writeTestFile(t, tempDir, "empty_var_server.xml", emptyVarXML)

	detector := framework.OpenLibertyDetector{}
	ctx := context.TODO()

	t.Run("Hardcoded Ports", func(t *testing.T) {
		component := &model.Component{
			Path: filepath.Join(tempDir, "hardcoded_server.xml"),
		}
		detector.DoPortsDetection(component, &ctx)
		assert.Equal(t, []int{1234, 1235}, component.Ports)
	})

	t.Run("Variable-based Ports", func(t *testing.T) {
		component := &model.Component{
			Path: filepath.Join(tempDir, "variable_server.xml"),
		}
		detector.DoPortsDetection(component, &ctx)
		assert.Equal(t, expectedPorts, component.Ports)
	})

	t.Run("Mixed Hardcoded and Variable Ports", func(t *testing.T) {
		component := &model.Component{
			Path: filepath.Join(tempDir, "mixed_server.xml"),
		}
		detector.DoPortsDetection(component, &ctx)
		assert.Equal(t, []int{9080, 1235}, component.Ports)
	})

	t.Run("Empty Variable Port", func(t *testing.T) {
		component := &model.Component{
			Path: filepath.Join(tempDir, "empty_var_server.xml"),
		}
		detector.DoPortsDetection(component, &ctx)
		assert.Equal(t, []int{9443}, component.Ports)
	})
}

func writeTestFile(t *testing.T, dir, filename, content string) {
	path := filepath.Join(dir, filename)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
}
