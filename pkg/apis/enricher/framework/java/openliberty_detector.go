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
	"encoding/xml"
	"strings"

	"github.com/devfile/alizer/pkg/apis/model"
	"github.com/devfile/alizer/pkg/utils"
)

type OpenLibertyDetector struct{}

func (o OpenLibertyDetector) GetSupportedFrameworks() []string {
	return []string{"OpenLiberty"}
}

func (o OpenLibertyDetector) GetApplicationFileInfos(componentPath string, ctx *context.Context) []model.ApplicationFileInfo {
	return []model.ApplicationFileInfo{
		{
			Context: ctx,
			Root:    componentPath,
			Dir:     "",
			File:    "server.xml",
		},
		{
			Context: ctx,
			Root:    componentPath,
			Dir:     "src/main/liberty/config",
			File:    "server.xml",
		},
	}
}

// DoFrameworkDetection uses the groupId to check for the framework name
func (o OpenLibertyDetector) DoFrameworkDetection(language *model.Language, config string) {
	if hasFwk, _ := hasFramework(config, "io.openliberty", ""); hasFwk {
		language.Frameworks = append(language.Frameworks, "OpenLiberty")
	}
}

// DoPortsDetection searches for the port in src/main/liberty/config/server.xml and /server.xml
func (o OpenLibertyDetector) DoPortsDetection(component *model.Component, ctx *context.Context) {
	appFileInfos := o.GetApplicationFileInfos(component.Path, ctx)
	if len(appFileInfos) == 0 {
		return
	}

	for _, appFileInfo := range appFileInfos {
		fileBytes, err := utils.GetApplicationFileBytes(appFileInfo)
		if err != nil {
			continue
		}

		var data model.OpenLibertyServerXml
		err = xml.Unmarshal(fileBytes, &data)
		if err != nil {
			continue
		}

		variables := make(map[string]string)
		for _, v := range data.Variables {
			variables[v.Name] = v.DefaultValue
		}

		httpPort := resolvePort(data.HttpEndpoint.HttpPort, variables)
		httpsPort := resolvePort(data.HttpEndpoint.HttpsPort, variables)

		ports := utils.GetValidPorts([]string{httpPort, httpsPort})
		if len(ports) > 0 {
			component.Ports = ports
			return
		}
	}
}

// resolvePort resolves the port value by checking if it is a variable and if it is, it returns the value of the variable
func resolvePort(portValue string, variables map[string]string) string {
	if strings.HasPrefix(portValue, "${") && strings.HasSuffix(portValue, "}") {
		varName := strings.Trim(portValue, "${}")
		if value, exists := variables[varName]; exists {
			return value
		}
	}
	return portValue
}
