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
package recognizer

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/devfile/alizer/pkg/apis/model"
	"github.com/devfile/alizer/pkg/utils"
)

// component detection: dotnet, f-sharp, net-vb
func TestComponentDetectionOnDotNet(t *testing.T) {
	isComponentsInProject(t, "s2i-dotnetcore-ex", 1, "c#", "app")
}

func TestComponentDetectionOnFSharp(t *testing.T) {
	isComponentsInProject(t, "net-fsharp", 1, "f#", "net-fsharp")
}

func TestComponentDetectionOnVBNet(t *testing.T) {
	isComponentsInProject(t, "net-vb", 1, "Visual Basic .NET", "net-vb")
}

// port detection: dotnet
// not supported

// component detection: go
func TestComponentDetectionOnBeego(t *testing.T) {
	isComponentsInProject(t, "beego", 1, "Go", "beego")
}

func TestComponentDetectionOnEcho(t *testing.T) {
	isComponentsInProject(t, "echo", 1, "Go", "echo")
}

func TestComponentDetectionOnFastHTTP(t *testing.T) {
	isComponentsInProject(t, "fasthttp", 1, "Go", "fasthttp")
}

func TestComponentDetectionOnGin(t *testing.T) {
	isComponentsInProject(t, "golang-gin-app", 1, "Go", "golang-gin-app")
}

func TestComponentDetectionOnFiber(t *testing.T) {
	isComponentsInProject(t, "golang-fiber", 1, "Go", "golang-fiber")
}

func TestComponentDetectionOnMux(t *testing.T) {
	isComponentsInProject(t, "golang-mux", 1, "Go", "golang-mux")
}

// port detection: go
func TestPortDetectionGoBeego(t *testing.T) {
	testPortDetectionInProject(t, "beego", []int{1999})
}

func TestPortDetectionGoEcho(t *testing.T) {
	testPortDetectionInProject(t, "echo", []int{8585})
}

func TestPortDetectionGoFasthttp(t *testing.T) {
	testPortDetectionInProject(t, "fasthttp", []int{2999})
}

func TestPortDetectionGoGin(t *testing.T) {
	testPortDetectionInProject(t, "golang-gin-app", []int{8789})
}

func TestPortDetectionGo(t *testing.T) {
	testPortDetectionInProject(t, "golang-runtime", []int{8080})
}

func TestPortDetectionGoMux(t *testing.T) {
	testPortDetectionInProject(t, "golang-mux", []int{8000})
}

func TestPortDetectionGoFiber(t *testing.T) {
	testPortDetectionInProject(t, "golang-fiber", []int{3000})
}

// component detection: java
func TestComponentDetectionOnJBossEAP(t *testing.T) {
	isComponentsInProject(t, "jboss-eap", 1, "java", "jboss-eap")
}

func TestComponentDetectionOnMicronaut(t *testing.T) {
	isComponentsInProject(t, "micronaut", 1, "java", "myMicronautProject")
}

func TestComponentDetectionOnOpenLiberty(t *testing.T) {
	isComponentsInProject(t, "open-liberty", 1, "java", "openliberty")
}

func TestComponentDetectionOnQuarkus(t *testing.T) {
	isComponentsInProject(t, "quarkus", 1, "java", "code-with-quarkus-maven")
}

func TestComponentDetectionOnSpring(t *testing.T) {
	isComponentsInProject(t, "spring", 1, "java", "spring")
}

func TestComponentDetectionOnSpringCloud(t *testing.T) {
	isComponentsInProject(t, "spring-cloud", 1, "java", "spring-cloud")
}

func TestComponentDetectionOnVertx(t *testing.T) {
	isComponentsInProject(t, "vertx", 1, "java", "http-vertx")
}

func TestComponentDetectionOnWildFly(t *testing.T) {
	isComponentsInProject(t, "wildfly", 1, "java", "wildfly")
}

func TestComponentDetectionOnJakartaEE(t *testing.T) {
	isComponentsInProject(t, "jakartaee", 1, "java", "jakartaee-app")
}

// port detection: java
func TestPortDetectionJavaJBossEAP(t *testing.T) {
	testPortDetectionInProject(t, "jboss-eap", []int{8380})
}

func TestPortDetectionJavaMicronaut(t *testing.T) {
	testPortDetectionInProject(t, "micronaut", []int{4444})
}

func TestPortDetectionJavaMicronautFromEnvs(t *testing.T) {
	os.Setenv("MICRONAUT_SERVER_PORT", "1234")
	testPortDetectionInProject(t, "micronaut", []int{1234})
	os.Unsetenv("MICRONAUT_SERVER_PORT")
}

func TestPortDetectionJavaMicronautFromEnvsWithSSLEnabled(t *testing.T) {
	os.Setenv("MICRONAUT_SERVER_PORT", "1345")
	// Case MICRONAUT_SERVER_SSL_ENABLED is set to something else
	os.Setenv("MICRONAUT_SERVER_SSL_ENABLED", "something else")
	os.Setenv("MICRONAUT_SERVER_SSL_PORT", "1456")
	testPortDetectionInProject(t, "micronaut", []int{1345})
	os.Setenv("MICRONAUT_SERVER_SSL_ENABLED", "true")
	testPortDetectionInProject(t, "micronaut", []int{1345, 1456})
	os.Unsetenv("MICRONAUT_SERVER_PORT")
	os.Unsetenv("MICRONAUT_SERVER_SSL_PORT")
	os.Unsetenv("MICRONAUT_SERVER_SSL_ENABLED")
}

func TestPortDetectionJavaMicronautFromDockerfile(t *testing.T) {
	testPortDetectionInProject(t, "micronaut-dockerfile-one-port", []int{1345})
}

func TestPortDetectionJavaMicronautFromDockerfileWithSSLEnabled(t *testing.T) {
	testPortDetectionInProject(t, "micronaut-dockerfile-ssl-enabled", []int{1456, 1345})
}

func TestPortDetectionOnOpenLiberty(t *testing.T) {
	testPortDetectionInProject(t, "open-liberty", []int{9080, 9443})
	testOpenLibertyDetector_DoPortsDetection(t, "open-liberty", []int{9080, 9443})
}

func TestPortDetectionJavaQuarkus(t *testing.T) {
	testPortDetectionInProject(t, "quarkus", []int{9898})
}

func TestPortDetectionJavaQuarkusWithEnv(t *testing.T) {
	os.Setenv("QUARKUS_HTTP_SSL_PORT", "1456")
	// Check that only true is accepted as value
	os.Setenv("QUARKUS_HTTP_INSECURE_REQUESTS", "wrong")
	testPortDetectionInProject(t, "quarkus", []int{1456})
	os.Setenv("QUARKUS_HTTP_INSECURE_REQUESTS", "true")
	os.Setenv("QUARKUS_HTTP_PORT", "1235")
	testPortDetectionInProject(t, "quarkus", []int{1456, 1235})
	os.Unsetenv("QUARKUS_HTTP_SSL_PORT")
	os.Unsetenv("QUARKUS_HTTP_INSECURE_REQUESTS")
	os.Unsetenv("QUARKUS_HTTP_PORT")
}

func TestPortDetectionJavaQuarkusDockerfileOnePort(t *testing.T) {
	testPortDetectionInProject(t, "quarkus-dockerfile-one-port", []int{1345})
}

func TestPortDetectionJavaQuarkusDockerfileMultiplePorts(t *testing.T) {
	testPortDetectionInProject(t, "quarkus-dockerfile-multiple-ports", []int{1345, 1456})
}

func TestPortDetectionSpring(t *testing.T) {
	testPortDetectionInProject(t, "spring", []int{9012})
}

func TestPortDetectionSpringDockerfileSimple(t *testing.T) {
	testPortDetectionInProject(t, "spring-dockerfile-simple", []int{1345})
}

func TestPortDetectionJavaVertxHttpPort(t *testing.T) {
	testPortDetectionInProject(t, "vertx", []int{2321})
}

func TestPortDetectionJavaVertxServerPort(t *testing.T) {
	testPortDetectionInProject(t, "vertx-port-server", []int{5555})
}

func TestPortDetectionJavaWildfly(t *testing.T) {
	testPortDetectionInProject(t, "wildfly", []int{8085})
}

// component detection: javascript, typescript
func TestComponentDetectionOnJavascript(t *testing.T) {
	isComponentsInProject(t, "nodejs-ex", 1, "javascript", "nodejs-starter")
}

func TestComponentDetectionOnAngular(t *testing.T) {
	isComponentsInProject(t, "angularjs", 1, "typescript", "angularjs")
}

func TestComponentDetectionOnExpress(t *testing.T) {
	isComponentsInProject(t, "expressjs", 1, "javascript", "expressjs")
}

func TestComponentDetectionOnNextJs(t *testing.T) {
	isComponentsInProject(t, "nextjs-app", 1, "typescript", "nextjs-app")
}

func TestComponentDetectionOnNuxtJs(t *testing.T) {
	isComponentsInProject(t, "nuxtjs-app", 1, "typescript", "nuxt-app")
}

func TestComponentDetectionOnReactJs(t *testing.T) {
	isComponentsInProject(t, "reactjs", 1, "javascript", "simple-react-template")
}

func TestComponentDetectionOnSvelteJs(t *testing.T) {
	isComponentsInProject(t, "svelte-app", 1, "javascript", "svelte-app")
}

func TestComponentDetectionOnVue(t *testing.T) {
	isComponentsInProject(t, "vue-app", 1, "typescript", "vue-app")
}

// port detection: javascript, typescript
func TestPortDetectionAngularPortInStartScript(t *testing.T) {
	testPortDetectionInProject(t, "angularjs", []int{8780})
}

func TestPortDetectionJavascriptExpressClear(t *testing.T) {
	testPortDetectionInProject(t, "expressjs", []int{7777})
}

func TestPortDetectionJavascriptExpressEnv(t *testing.T) {
	os.Setenv("TEST_EXPRESS_ENV", "1111")
	testPortDetectionInProject(t, "expressjs-env", []int{1111})
	os.Unsetenv("TEST_EXPRESS_ENV")
}

func TestPortDetectionJavascriptExpressEnvOROperatorWithEnvVar(t *testing.T) {
	os.Setenv("TEST_EXPRESS_ENV", "1111")
	testPortDetectionInProject(t, "expressjs-env-logical-or", []int{1111, 8080})
	os.Unsetenv("TEST_EXPRESS_ENV")
}

func TestPortDetectionJavascriptExpressDockerfileEnvOROperatorWithEnvVar(t *testing.T) {
	os.Setenv("TEST_EXPRESS_DOCKERFILE_ENV", "1111")
	testPortDetectionInProject(t, "expressjs-dockerfile-env-logical-or", []int{1111, 8080})
	os.Unsetenv("TEST_EXPRESS_DOCKERFILE_ENV")
	testPortDetectionInProject(t, "expressjs-dockerfile-env-logical-or", []int{1345, 8080})
}

func TestPortDetectionJavascriptExpressEnvOROperatorWithoutEnvVar(t *testing.T) {
	testPortDetectionInProject(t, "expressjs-env-logical-or", []int{8080})
}
func TestPortDetectionJavascriptExpressVariable(t *testing.T) {
	testPortDetectionInProject(t, "expressjs-variable", []int{3000})
}

func TestPortDetectionJavascriptExpressDockerfileEnvVar(t *testing.T) {
	testPortDetectionInProject(t, "expressjs-dockerfile-env", []int{1345})
}

func TestPortDetectionNextJsPortInStartScript(t *testing.T) {
	testPortDetectionInProject(t, "nextjs-app", []int{8610})
}

func TestPortDetectionNuxtJsPortInConfigFile(t *testing.T) {
	testPortDetectionInProject(t, "nuxtjs-app", []int{8787})
}

func TestPortDetectionJavascriptReactEnvVariable(t *testing.T) {
	oldValue := os.Getenv("PORT")
	os.Setenv("PORT", "2121")
	testPortDetectionInProject(t, "reactjs", []int{2121})
	if oldValue == "" {
		os.Unsetenv("PORT")
	} else {
		os.Setenv("PORT", oldValue)
	}
}

func TestPortDetectionJavascriptReactEnvFile(t *testing.T) {
	testPortDetectionInProject(t, "reactjs", []int{1231})
}

func TestPortDetectionJavascriptReactEnvDockerfile(t *testing.T) {
	testPortDetectionInProject(t, "reactjs-dockerfile-env", []int{4526})
}

func TestPortDetectionJavascriptReactScript(t *testing.T) {
	testPortDetectionInProject(t, "reactjs-script", []int{5353})
}

func TestPortDetectionSvelteJsPortInStartScript(t *testing.T) {
	testPortDetectionInProject(t, "svelte-app", []int{8282})
}

func TestPortDetectionVuePortInStartScript(t *testing.T) {
	testPortDetectionInProject(t, "vue-app", []int{8282})
}

func TestPortDetectionVuePortEnvDockerfile(t *testing.T) {
	testPortDetectionInProject(t, "vue-app-dockerfile-simple", []int{4526})
}

// component detection: php
func TestComponentDetectionOnLaravel(t *testing.T) {
	isComponentsInProject(t, "laravel", 1, "PHP", "laravel")
}

// port detection: php
func TestPortDetectionPHPLaravel(t *testing.T) {
	testPortDetectionInProject(t, "laravel", []int{9988})
}

// component detection: python
func TestComponentDetectionOnDjango(t *testing.T) {
	isComponentsInProject(t, "django", 1, "python", "django")
}

func TestComponentDetectionOnFlask(t *testing.T) {
	isComponentsInProject(t, "flask", 1, "python", "flask")
}

// port detection: python
func TestPortDetectionDjango(t *testing.T) {
	testPortDetectionInProject(t, "django", []int{3543})
}

func TestPortDetectionFlask(t *testing.T) {
	testPortDetectionInProject(t, "flask", []int{3000})
}

func TestPortDetectionFlaskAssignedVariable(t *testing.T) {
	testPortDetectionInProject(t, "flask-port-assigned-variable", []int{3001})
}

func TestPortDetectionFlaskStringValue(t *testing.T) {
	testPortDetectionInProject(t, "flask-port-string-value", []int{})
}

// component detection: corner cases
func TestComponentDetectionNoResult(t *testing.T) {
	components := getComponentsFromTestProject(t, "simple")
	if len(components) > 0 {
		t.Errorf("Expected 0 components but found %v", strconv.Itoa(len(components)))
	}
}

func TestComponentDetectionOnDoubleComponents(t *testing.T) {
	isComponentsInProject(t, "double-components", 2, "javascript", "")
}

func TestComponentDetectionWithDoubleComponentsDockerFile(t *testing.T) {
	isComponentsInProject(t, "dockerfile-double-components", 2, "python", "")
}

func TestComponentDetectionWithNestedComponents(t *testing.T) {
	isComponentsInProject(t, "dockerfile-nested-inner-double-components", 2, "python", "")
}

func TestComponentDetectionWithNestedComponents2(t *testing.T) {
	isComponentsInProject(t, "dockerfile-nested-inner", 1, "python", "")
}

func TestComponentDetectionWithGitIgnoreRule(t *testing.T) {
	testingProjectPath := getTestProjectPath("component-wrapped-in-folder")
	settings := model.DetectionSettings{
		BasePath: testingProjectPath,
	}
	ctx := context.Background()
	files, err := utils.GetCachedFilePathsFromRoot(testingProjectPath, &ctx)
	if err != nil {
		t.Error(err)
	}

	components := getComponentsFromFiles(t, files, settings)

	if len(components) != 1 {
		t.Errorf("Expected 1 components but found %v", strconv.Itoa(len(components)))
	}

	//now add a gitIgnore with a rule to exclude the only component found
	gitIgnorePath := filepath.Join(testingProjectPath, ".gitignore")
	err = updateContent(gitIgnorePath, []byte("**/quarkus/"))
	if err != nil {
		t.Error(err)
	}
	files, err = utils.GetFilePathsFromRoot(testingProjectPath)
	if err != nil {
		t.Error(err)
	}
	componentsWithUpdatedGitIgnore := getComponentsFromFiles(t, files, settings)
	//delete gitignore file
	os.Remove(gitIgnorePath)

	if len(componentsWithUpdatedGitIgnore) != 0 {
		t.Errorf("Expected 0 components but found %v", strconv.Itoa(len(componentsWithUpdatedGitIgnore)))
	}
}

func TestComponentDetectionMultiProjects(t *testing.T) {
	components := getComponentsFromTestProject(t, "")
	nComps := 72
	if len(components) != nComps {
		t.Errorf("Expected %v components but found %v", strconv.Itoa(nComps), strconv.Itoa(len(components)))
	}
}

// port detection: other cases
func TestPortDetectionWithDockerComposeExpose(t *testing.T) {
	testPortDetectionInProject(t, "docker-compose-expose", []int{3000, 8000})
}

func TestPortDetectionWithDockerComposeShortSyntaxPorts(t *testing.T) {
	testPortDetectionInProject(t, "docker-compose-ports-short-syntax", []int{3000, 1234})
}

func TestPortDetectionWithDockerComposeLongSyntaxPorts(t *testing.T) {
	testPortDetectionInProject(t, "docker-compose-ports-long-syntax", []int{6060})
}

func TestPortDetectionWithDockerFile(t *testing.T) {
	testPortDetectionInProject(t, "dockerfile-simple", []int{8085})
}

func TestPortDetectionWithContainerFile(t *testing.T) {
	testPortDetectionInProject(t, "containerfile-simple", []int{8085})
}

func TestPortDetectionWithOrphanContainerFile(t *testing.T) {
	testPortDetectionInProject(t, "containerfile-orphan", []int{8090})
}

func TestPortDetectionWithOrphanDockerFile(t *testing.T) {
	testPortDetectionInProject(t, "dockerfile-orphan", []int{8085})
}

func TestPortDetectionWithSecondLevelDockerFile(t *testing.T) {
	testPortDetectionInProject(t, "dockerfile-nested", []int{8085})
}

func TestPortDetectionWithSecondLevelContainerFile(t *testing.T) {
	testPortDetectionInProject(t, "containerfile-nested", []int{8085})
}
