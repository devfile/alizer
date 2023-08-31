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

package model

import (
	"context"
	"regexp"
)

type PortDetectionAlgorithm int

const (
	DockerFile PortDetectionAlgorithm = 0
	Compose    PortDetectionAlgorithm = 1
	Source     PortDetectionAlgorithm = 2
)

type DetectionSettings struct {
	BasePath              string
	PortDetectionStrategy []PortDetectionAlgorithm
}

type Language struct {
	Name                    string
	Aliases                 []string
	Weight                  float64
	Frameworks              []string
	Tools                   []string
	CanBeComponent          bool
	CanBeContainerComponent bool
}

type Component struct {
	Name      string
	Path      string
	Languages []Language
	Ports     []int
}

type Version struct {
	SchemaVersion string
	Default       bool
	Version       string
}

type DevfileType struct {
	Name        string
	Language    string
	ProjectType string
	Tags        []string
	Versions    []Version
}

type DevfileFilter struct {
	MinSchemaVersion string
	MaxSchemaVersion string
}

type ApplicationFileInfo struct {
	Context *context.Context
	Root    string
	Dir     string
	File    string
}

type PortMatchRules struct {
	MatchIndexRegexes []PortMatchRule
	MatchRegexes      []PortMatchSubRule
}

type PortMatchRule struct {
	Regex     *regexp.Regexp
	ToReplace string
}

type PortMatchSubRule struct {
	Regex    *regexp.Regexp
	SubRegex *regexp.Regexp
}

type DevfileScore struct {
	DevfileIndex int
	Score        int
}

// EnvVar represents an environment variable with a name and a corresponding value.
type EnvVar struct {
	// Name is the name of the environment variable.
	Name string

	// Value is the value associated with the environment variable.
	Value string
}
