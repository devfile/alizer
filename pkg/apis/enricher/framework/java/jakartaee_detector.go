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

	"github.com/devfile/alizer/pkg/apis/model"
)

type JakartaEEDetector struct{}

func (j JakartaEEDetector) GetSupportedFrameworks() []string {
	return []string{"JakartaEE"}
}

func (j JakartaEEDetector) GetApplicationFileInfos(componentPath string, ctx *context.Context) []model.ApplicationFileInfo {
	return []model.ApplicationFileInfo{}
}

// DoFrameworkDetection uses the groupId to check for the framework name
func (j JakartaEEDetector) DoFrameworkDetection(language *model.Language, config string) {
	jakartaGroupIds := []string{
		"jakarta.platform",    // Jakarta EE Platform BOM
		"jakarta.servlet",     // Servlet API
		"jakarta.ws.rs",       // JAX-RS (RESTful Web Services)
		"jakarta.persistence", // JPA (Persistence)
		"jakarta.enterprise",  // CDI (Contexts and Dependency Injection)
	}

	for _, groupId := range jakartaGroupIds {
		if hasFwk, _ := hasFramework(config, groupId, ""); hasFwk {
			language.Frameworks = append(language.Frameworks, "JakartaEE")
			return
		}
	}
}

// DoPortsDetection is not implemented for JakartaEE as port configuration varies by runtime
func (j JakartaEEDetector) DoPortsDetection(component *model.Component, ctx *context.Context) {
}
