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

type WebSphereDetector struct{}

func (o WebSphereDetector) GetSupportedFrameworks() []string {
	return []string{"WebSphere"}
}

// DoFrameworkDetection uses the groupId and artifactId to check for the framework name
func (o WebSphereDetector) DoFrameworkDetection(language *model.Language, config string) {
	hasWebSphereFwk, _ := hasFramework(config, "com.ibm.websphere.appserver", "")
	hasOpenLibertyFwk, _ := hasFramework(config, "io.openliberty", "")
	if hasWebSphereFwk && !hasOpenLibertyFwk {
		language.Frameworks = append(language.Frameworks, "WebSphere")
	}
}

func (o WebSphereDetector) DoPortsDetection(component *model.Component, ctx *context.Context) {
}
