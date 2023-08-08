package recognizer

/*******************************************************************************
 * Copyright (c) 2022 Red Hat, Inc.
 * Distributed under license by Red Hat, Inc. All rights reserved.
 * This program is made available under the terms of the
 * Eclipse Public License v2.0 which accompanies this distribution,
 * and is available at http://www.eclipse.org/legal/epl-v20.html
 *
 * Contributors:
 * Red Hat, Inc.
 ******************************************************************************/
import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/devfile/alizer/pkg/apis/model"
	"github.com/stretchr/testify/assert"
)

func TestDetectQuarkusDevfile(t *testing.T) {
	detectDevfiles(t, "quarkus", []string{"java-quarkus"})
}

func TestDetectMicronautDevfile(t *testing.T) {
	detectDevfiles(t, "micronaut", []string{"java-maven"})
}

func TestDetectNodeJSDevfile(t *testing.T) {
	detectDevfiles(t, "nodejs-ex", []string{"nodejs"})
}

func TestDetectDotNet50Devfile(t *testing.T) {
	detectDevfile(t, "dotnet5.0", []string{"dotnet50"})
}

func TestDetectDotNet60Devfile(t *testing.T) {
	detectDevfile(t, "dotnet6.0", []string{"dotnet60"})
}

func TestDetectDotNetCore31Devfile(t *testing.T) {
	detectDevfile(t, "dotnetcore3.1", []string{"dotnetcore31"})
}

func TestDetectDjangoDevfile(t *testing.T) {
	detectDevfiles(t, "django", []string{"python-django"})
}

func TestDetectFlaskDevfile(t *testing.T) {
	detectDevfiles(t, "flask", []string{"python"})
}

func TestDetectWildflyDevfile(t *testing.T) {
	detectDevfiles(t, "wildfly", []string{"java-wildfly"})
}

func TestDetectDjangoDevfileUsingLanguages(t *testing.T) {
	languages := []model.Language{
		{
			Name: "Python",
			Aliases: []string{
				"python3",
			},
			Weight: 88.23,
			Frameworks: []string{
				"Django",
			},
			Tools:          []string{},
			CanBeComponent: false,
		},
		{
			Name: "Shell",
			Aliases: []string{
				"sh",
			},
			Weight:         11.77,
			Frameworks:     []string{},
			Tools:          []string{},
			CanBeComponent: false,
		},
	}
	detectDevfilesUsingLanguages(t, "", languages, []string{"python-django"})
}

func TestDetectQuarkusDevfileUsingLanguages(t *testing.T) {
	detectDevfilesUsingLanguages(t, "quarkus", []model.Language{}, []string{"java-quarkus"})
}

func TestDetectMicronautDevfileUsingLanguages(t *testing.T) {
	detectDevfilesUsingLanguages(t, "micronaut", []model.Language{}, []string{"java-maven"})
}

func TestDetectNodeJSDevfileUsingLanguages(t *testing.T) {
	detectDevfilesUsingLanguages(t, "nodejs-ex", []model.Language{}, []string{"nodejs"})
}

func TestDetectGoDevfile(t *testing.T) {
	detectDevfiles(t, "golang-gin-app", []string{"go"})
}

func TestDetectAngularDevfile(t *testing.T) {
	detectDevfiles(t, "angularjs", []string{"Angular"})
}

func TestDetectVertxDevfile(t *testing.T) {
	detectDevfile(t, "vertx", []string{"java-vertx"})
}

func TestDetectNextJsDevfile(t *testing.T) {
	detectDevfiles(t, "nextjs-app", []string{"Next.js"})
}

func TestDetectNuxtJsDevfile(t *testing.T) {
	detectDevfiles(t, "nuxtjs-app", []string{"nodejs-nuxtjs", "nodejs-vue"})
}

func TestDetectVueDevfile(t *testing.T) {
	detectDevfiles(t, "vue-app", []string{"nodejs-vue"})
}

func TestDetectReactJSDevfile(t *testing.T) {
	detectDevfiles(t, "reactjs", []string{"nodejs-react"})
}

func TestDetectSvelteDevfile(t *testing.T) {
	detectDevfiles(t, "svelte-app", []string{"nodejs-svelte"})
}

func TestDetectSpringDevfile(t *testing.T) {
	detectDevfiles(t, "spring", []string{"java-spring", "java-springboot"})
}

func TestDetectLaravelDevfile(t *testing.T) {
	detectDevfiles(t, "laravel", []string{"php-laravel"})
}

func TestGetUrlWithVersions(t *testing.T) {
	tests := []struct {
		name             string
		minSchemaVersion string
		maxSchemaVersion string
		testUrl          string
		expectedError    error
		expectedUrl      string
	}{
		{

			name:             "Case 1: Url with valid min and max versions",
			minSchemaVersion: "2.0.0",
			maxSchemaVersion: "2.2.0",
			testUrl:          "http://localhost:5000/",
			expectedError:    nil,
		},
		{
			name:             "Case 2: Url with valid min version",
			minSchemaVersion: "2.0.0",
			maxSchemaVersion: "",
			testUrl:          "http://localhost:5000",
			expectedError:    nil,
		},
		{
			name:             "Case 3: Url with valid max version",
			minSchemaVersion: "",
			maxSchemaVersion: "2.2.0",
			testUrl:          "http://localhost:5000/",
			expectedError:    nil,
		},
		{
			name:             "Case 4: Url with max version lower than min version",
			minSchemaVersion: "2.2.0",
			maxSchemaVersion: "2.1.0",
			testUrl:          "http://localhost:5000/v2index",
			expectedError:    fmt.Errorf("max-schema-version cannot be lower than min-schema-version"),
		},
		{
			name:             "Case 5: Url with max version lower than minimum allowed version",
			minSchemaVersion: "0.1.0",
			maxSchemaVersion: "1.1.0",
			testUrl:          "http://localhost:5000/v2index",
			expectedError:    fmt.Errorf("min and/or max version are lower than the minimum allowed version (2.0.0)"),
		},
		{
			name:             "Case 6: Url with max-schema-version lower than minimum allowed version & min-schema-version",
			minSchemaVersion: "2.1.0",
			maxSchemaVersion: "1.1.0",
			testUrl:          "http://localhost:5000/v2index",
			expectedError:    fmt.Errorf("max-schema-version cannot be lower than min-schema-version"),
		},
		{
			name:             "Case 7: Url with min version lower than minimum allowed version",
			minSchemaVersion: "1.1.0",
			maxSchemaVersion: "",
			testUrl:          "http://localhost:5000/v2index",
			expectedError:    fmt.Errorf("min version is lower than the minimum allowed version (2.0.0)"),
		},
		{
			name:             "Case 8: Url with max version lower than minimum allowed version",
			minSchemaVersion: "",
			maxSchemaVersion: "1.1.0",
			testUrl:          "http://localhost:5000/v2index",
			expectedError:    fmt.Errorf("max version is lower than the minimum allowed version (2.0.0)"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetUrlWithVersions(tt.testUrl, tt.minSchemaVersion, tt.maxSchemaVersion)
			if err != nil {
				assert.EqualValues(t, tt.expectedError, err)
			} else {
				assert.EqualValues(t, getExceptedVersionsUrl(tt.testUrl, tt.minSchemaVersion, tt.maxSchemaVersion, tt.expectedError), result)
			}
		})
	}
}

func TestMatchDevfiles(t *testing.T) {
	tests := []struct {
		name                string
		filter              model.DevfileFilter
		mockFunc            func()
		url                 string
		path                string
		expectedDevfileType model.DevfileType
		expectingErr        bool
	}{
		{
			name:   "Match devfile success",
			filter: model.DevfileFilter{},
			path:   "some-path",
			url:    "some-url",
			expectedDevfileType: model.DevfileType{
				Name:        "mocked-stack",
				Language:    "python",
				ProjectType: "python",
				Tags:        []string{"python"},
				Versions:    []model.Version{},
			},
			expectingErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			// mock DownloadDevfileTypesFromRegistry
			DownloadDevfileTypesFromRegistry = func(url string, filter model.DevfileFilter) ([]model.DevfileType, error) {
				return []model.DevfileType{tc.expectedDevfileType}, nil
			}
			selectDevfiles = func(path string, devfileTypesFromRegistry []model.DevfileType) ([]model.DevfileType, error) {
				return []model.DevfileType{tc.expectedDevfileType}, nil
			}
			devfileType, err := MatchDevfiles("test-path", "some-url", model.DevfileFilter{})

			errExist := err != nil
			if tc.expectingErr != errExist {
				tt.Errorf("Error expectation not met, want %v, got %v", tc.expectingErr, errExist)
			}

			if !reflect.DeepEqual(tc.expectedDevfileType.Name, devfileType[0].Name) {
				tt.Errorf("Error, user profile expectation not met, want %+v, got %+v", tc.expectedDevfileType, devfileType)
			}
		})
	}
}

func getExceptedVersionsUrl(url, minSchemaVersion, maxSchemaVersion string, err error) string {
	if err != nil {
		return ""
	} else if minSchemaVersion != "" && maxSchemaVersion != "" {
		return fmt.Sprintf("%s?minSchemaVersion=%s&maxSchemaVersion=%s", url, minSchemaVersion, maxSchemaVersion)
	} else if minSchemaVersion != "" {
		return fmt.Sprintf("%s?minSchemaVersion=%s", url, minSchemaVersion)
	} else if maxSchemaVersion != "" {
		return fmt.Sprintf("%s?maxSchemaVersion=%s", url, maxSchemaVersion)
	} else {
		return url
	}
}

func detectDevfiles(t *testing.T, projectName string, devfilesName []string) {
	detectDevfilesFunc := func(devfileTypes []model.DevfileType) ([]int, error) {
		testingProjectPath := getTestProjectPath(projectName)
		return SelectDevFilesFromTypes(testingProjectPath, devfileTypes)
	}
	detectDevfilesInner(t, devfilesName, detectDevfilesFunc)
}

func detectDevfile(t *testing.T, projectName string, devfilesName []string) {
	detectDevfilesFunc := func(devfileTypes []model.DevfileType) ([]int, error) {
		testingProjectPath := getTestProjectPath(projectName)
		devfileIndex, err := SelectDevFileFromTypes(testingProjectPath, devfileTypes)
		return []int{devfileIndex}, err
	}
	detectDevfilesInner(t, devfilesName, detectDevfilesFunc)
}

func detectDevfilesUsingLanguages(t *testing.T, projectName string, languages []model.Language, devfileName []string) {
	if projectName != "" {
		testingProjectPath := getTestProjectPath(projectName)
		var err error
		languages, err = Analyze(testingProjectPath)
		if err != nil {
			t.Error(err)
		}
	}
	detectDevfileFunc := func(devfileTypes []model.DevfileType) ([]int, error) {
		return SelectDevfilesUsingLanguagesFromTypes(languages, devfileTypes)
	}
	detectDevfilesInner(t, devfileName, detectDevfileFunc)
}

func detectDevfilesInner(t *testing.T, expectedDevfilesName []string, detectFuncInner func([]model.DevfileType) ([]int, error)) {
	devfileTypes := getDevfileTypes()
	foundDevfilesIndexes, err := detectFuncInner(devfileTypes)
	if err != nil {
		t.Error(err)
	}

	found := false
	for _, expectedDevfile := range expectedDevfilesName {
		found = false
		for _, foundDevfileIndex := range foundDevfilesIndexes {
			if devfileTypes[foundDevfileIndex].Name == expectedDevfile {
				found = true
			}
		}
		if !found {
			t.Error("Expected value " + expectedDevfile + " but it was not found")
			return
		}
	}
}

func getDevfileTypes() []model.DevfileType {
	return []model.DevfileType{
		{
			Name:        "java",
			Language:    "java",
			ProjectType: "java",
			Tags:        make([]string, 0),
		},
		{
			Name:        "udi",
			Language:    "Polyglot",
			ProjectType: "default",
			Tags: []string{
				"Java",
				"Maven",
				"Scala",
				"PHP",
				".NET",
				"Node.js",
				"Go",
				"Python",
				"Pip",
				"ubi8",
			},
		},
		{
			Name:        "dotnet50",
			Language:    ".NET",
			ProjectType: "dotnet",
			Tags: []string{
				".Net",
				".Net 5.0",
			},
		},
		{
			Name:        "dotnet60",
			Language:    ".NET",
			ProjectType: "dotnet",
			Tags: []string{
				".Net",
				".Net 6.0",
			},
		},
		{
			Name:        "dotnetcore31",
			Language:    ".NET",
			ProjectType: "dotnet",
			Tags: []string{
				".Net",
				".Net Core App 3.1",
			},
		},
		{
			Name:        "java-quarkus",
			Language:    "java",
			ProjectType: "quarkus",
			Tags: []string{
				"Java",
				"Quarkus",
			},
		},
		{
			Name:        "java-maven",
			Language:    "java",
			ProjectType: "java",
			Tags: []string{
				"Java",
				"Maven",
			},
		},
		{
			Name:        "java-spring",
			Language:    "java",
			ProjectType: "spring",
			Tags: []string{
				"Java",
				"Spring",
			},
		},
		{
			Name:        "java-springboot",
			Language:    "java",
			ProjectType: "springboot",
			Tags: []string{
				"Java",
				"Spring Boot",
			},
		},
		{
			Name:        "java-vertx",
			Language:    "java",
			ProjectType: "vertx",
			Tags: []string{
				"Java",
				"Vert.x",
			},
		},
		{
			Name:        "java-wildfly",
			Language:    "java",
			ProjectType: "wildfly",
			Tags: []string{
				"Java",
				"Wildfly",
			},
		},
		{
			Name:        "nodejs",
			Language:    "nodejs",
			ProjectType: "nodejs",
			Tags: []string{
				"NodeJS",
				"Express",
			},
		},
		{
			Name:        "python-django",
			Language:    "python",
			ProjectType: "django",
			Tags: []string{
				"Python",
				"pip",
			},
		},
		{
			Name:        "python",
			Language:    "python",
			ProjectType: "Python",
			Tags: []string{
				"Python",
				"pip",
				"Flask",
			},
		},
		{
			Name:        "go",
			Language:    "go",
			ProjectType: "go",
			Tags: []string{
				"go",
			},
		},
		{
			Name:        "Angular",
			Language:    "TypeScript",
			ProjectType: "Angular",
			Tags: []string{
				"NodeJS",
				"Angular",
			},
		},
		{
			Name:        "Next.js",
			Language:    "TypeScript",
			ProjectType: "Next.js",
			Tags: []string{
				"Node.js",
				"Next.js",
			},
		},
		{
			Name:        "nodejs-react",
			Language:    "TypeScript",
			ProjectType: "React",
			Tags: []string{
				"Node.js",
				"React",
			},
		},
		{
			Name:        "nodejs-svelte",
			Language:    "TypeScript",
			ProjectType: "Svelte",
			Tags: []string{
				"Node.js",
				"Svelte",
			},
		},
		{
			Name:        "nodejs-nuxtjs",
			Language:    "TypeScript",
			ProjectType: "Nuxt.js",
			Tags: []string{
				"Node.js",
				"Nuxt.js",
			},
		},
		{
			Name:        "nodejs-vue",
			Language:    "TypeScript",
			ProjectType: "Vue",
			Tags: []string{
				"Node.js",
				"Vue",
			},
		},
		{
			Name:        "php-laravel",
			Language:    "PHP",
			ProjectType: "Laravel",
			Tags: []string{
				"PHP",
				"Composer",
				"Laravel",
			},
		},
	}
}

func getTestProjectPath(folder string) string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return filepath.Join(basepath, "..", "..", "..", "resources/projects", folder)
}
