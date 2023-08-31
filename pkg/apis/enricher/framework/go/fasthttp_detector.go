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
	"golang.org/x/mod/modfile"
)

type FastHttpDetector struct{}

func (f FastHttpDetector) GetSupportedFrameworks() []string {
	return []string{"FastHttp"}
}

func (f FastHttpDetector) GetApplicationFileInfos(componentPath string, ctx *context.Context) []model.ApplicationFileInfo {
	files, err := utils.GetCachedFilePathsFromRoot(componentPath, ctx)
	if err != nil {
		return []model.ApplicationFileInfo{}
	}
	return utils.GenerateApplicationFileFromFilters(files, componentPath, ".go", ctx)
}

// DoFrameworkDetection uses a tag to check for the framework name
func (f FastHttpDetector) DoFrameworkDetection(language *model.Language, goMod *modfile.File) {
	if hasFramework(goMod.Require, "github.com/valyala/fasthttp") {
		language.Frameworks = append(language.Frameworks, "FastHttp")
	}
}

func (f FastHttpDetector) DoPortsDetection(component *model.Component, ctx *context.Context) {
	fileContents, err := utils.GetApplicationFileContents(f.GetApplicationFileInfos(component.Path, ctx))
	if err != nil {
		return
	}

	matchRegexRules := model.PortMatchRules{
		MatchIndexRegexes: []model.PortMatchRule{
			{
				Regex:     regexp.MustCompile(`.ListenAndServe\([^,)]*`),
				ToReplace: ".ListenAndServe(",
			},
		},
	}
	for _, fileContent := range fileContents {
		ports := GetPortFromFileGo(matchRegexRules, fileContent)
		if len(ports) > 0 {
			component.Ports = ports
			return
		}
	}
}
