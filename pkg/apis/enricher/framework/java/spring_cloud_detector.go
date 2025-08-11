/*******************************************************************************
 * Copyright (c) 2021 Red Hat, Inc.
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
	"path/filepath"

	"github.com/devfile/alizer/pkg/apis/model"
	"github.com/devfile/alizer/pkg/utils"
)

type SpringCloudDetector struct{}

func (s SpringCloudDetector) GetSupportedFrameworks() []string {
	return []string{"Spring Cloud"}
}

func (s SpringCloudDetector) GetApplicationFileInfos(componentPath string, ctx *context.Context) []model.ApplicationFileInfo {
	return []model.ApplicationFileInfo{
		{
			Context: ctx,
			Root:    componentPath,
			Dir:     "src/main/resources",
			File:    "application.properties",
		},
		{
			Context: ctx,
			Root:    componentPath,
			Dir:     "src/main/resources",
			File:    "application.yml",
		},
		{
			Context: ctx,
			Root:    componentPath,
			Dir:     "src/main/resources",
			File:    "application.yaml",
		},
	}
}

// DoFrameworkDetection uses the groupId to check for the framework name
func (s SpringCloudDetector) DoFrameworkDetection(language *model.Language, config string) {
	if hasFwk, _ := hasFramework(config, "org.springframework.cloud", ""); hasFwk {
		language.Frameworks = append(language.Frameworks, s.GetSupportedFrameworks()...)
	}
}

// DoPortsDetection searches for ports in the env var and
// src/main/resources/application.properties, or src/main/resources/application.yaml
func (s SpringCloudDetector) DoPortsDetection(component *model.Component, ctx *context.Context) {
	// case: port is set on env var
	ports := getSpringPortsFromEnvs()
	if len(ports) > 0 {
		component.Ports = ports
		return
	}

	// check if port is set on env var of dockerfile
	ports = getSpringPortsFromEnvDockerfile(component.Path)
	if len(ports) > 0 {
		component.Ports = ports
		return
	}

	// check if port is set inside application file
	appFileInfos := s.GetApplicationFileInfos(component.Path, ctx)
	if len(appFileInfos) == 0 {
		return
	}

	applicationFile := utils.GetAnyApplicationFilePath(component.Path, appFileInfos, ctx)
	if applicationFile == "" {
		return
	}

	var err error
	if filepath.Ext(applicationFile) == ".yml" || filepath.Ext(applicationFile) == ".yaml" {
		ports, err = getServerPortsFromYamlFile(applicationFile)
	} else {
		ports, err = getServerPortsFromPropertiesFile(applicationFile)
	}
	if err != nil {
		return
	}
	component.Ports = ports
}
