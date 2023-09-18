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

package enricher

import (
	"context"
	"regexp"

	"github.com/devfile/alizer/pkg/apis/model"
	"github.com/devfile/alizer/pkg/utils"
)

type VueDetector struct{}

func (v VueDetector) GetSupportedFrameworks() []string {
	return []string{"Vue"}
}

func (v VueDetector) GetApplicationFileInfos(componentPath string, ctx *context.Context) []model.ApplicationFileInfo {
	return []model.ApplicationFileInfo{
		{
			Context: ctx,
			Root:    componentPath,
			Dir:     "",
			File:    "vue.config.js",
		},
	}
}

// DoFrameworkDetection uses a tag to check for the framework name
func (v VueDetector) DoFrameworkDetection(language *model.Language, config string) {
	if hasFramework(config, "vue") {
		language.Frameworks = append(language.Frameworks, "Vue")
	}
}

// DoPortsDetection searches for the port in package.json, .env file, and vue.config.js
func (v VueDetector) DoPortsDetection(component *model.Component, ctx *context.Context) {
	regexes := []string{`--port (\d*)`, `PORT=(\d*)`}
	// check if --port or PORT is set in start script in package.json
	port := getPortFromStartScript(component.Path, regexes)
	if utils.IsValidPort(port) {
		component.Ports = []int{port}
	}

	// check if --port or PORT is set in dev script in package.json
	port = getPortFromDevScript(component.Path, regexes)
	if utils.IsValidPort(port) {
		component.Ports = []int{port}
	}

	// check if port is set on .env file
	port = utils.GetPortValueFromEnvFile(component.Path, `PORT=(\d*)`)
	if utils.IsValidPort(port) {
		component.Ports = []int{port}
		return
	}

	// check if port is set on as env var inside a dockerfile
	ports, err := utils.GetEnvVarPortValueFromDockerfile(component.Path, []string{"PORT"})
	if err == nil {
		component.Ports = ports
		return
	}

	//check inside the vue.config.js file
	appFileInfos := v.GetApplicationFileInfos(component.Path, ctx)
	if len(appFileInfos) == 0 {
		return
	}

	fileBytes, err := utils.GetApplicationFileBytes(appFileInfos[0])
	if err != nil {
		return
	}

	re := regexp.MustCompile(`port:\s*(\d+)*`)
	component.Ports = utils.FindAllPortsSubmatch(re, string(fileBytes), 1)
}
