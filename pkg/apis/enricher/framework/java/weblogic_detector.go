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

type WebLogicDetector struct{}

func (o WebLogicDetector) GetSupportedFrameworks() []string {
	return []string{"WebLogic"}
}

// DoFrameworkDetection uses the groupId and artifactId to check for the framework name
func (o WebLogicDetector) DoFrameworkDetection(language *model.Language, config string) {
	if hasFwk, _ := hasFramework(config, "com.oracle.weblogic", ""); hasFwk {
		language.Frameworks = append(language.Frameworks, "WebLogic")
	}
}

func (o WebLogicDetector) DoPortsDetection(component *model.Component, ctx *context.Context) {
}
