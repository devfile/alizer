package enricher

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/devfile/alizer/pkg/apis/model"
	"github.com/stretchr/testify/assert"
)

func TestOpenLibertyDetector_DoPortsDetection(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "openliberty-test")
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

	detector := OpenLibertyDetector{}

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

		assert.Equal(t, []int{9080, 9443}, component.Ports)
	})
}

func writeTestFile(t *testing.T, dir, filename, content string) {
	path := filepath.Join(dir, filename)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}
}
