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

const (
	DockerFile PortDetectionAlgorithm = 0
	Compose    PortDetectionAlgorithm = 1
	Source     PortDetectionAlgorithm = 2
)

// All models inside model.go are sorted by name A-Z

type AngularCliJson struct {
	Defaults struct {
		Serve AngularHostPort `json:"serve"`
	} `json:"defaults"`
}

type AngularHostPort struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type AngularJson struct {
	Projects map[string]AngularProjectBody `json:"projects"`
}

type AngularProjectBody struct {
	Architect struct {
		Serve struct {
			Options AngularHostPort `json:"options"`
		} `json:"serve"`
	} `json:"architect"`
}

type ApplicationFileInfo struct {
	Context *context.Context
	Root    string
	Dir     string
	File    string
}

type Component struct {
	Name      string
	Path      string
	Languages []Language
	Ports     []int
}

type DetectionSettings struct {
	BasePath              string
	PortDetectionStrategy []PortDetectionAlgorithm
}

type DevfileFilter struct {
	MinSchemaVersion string
	MaxSchemaVersion string
}

type DevfileScore struct {
	DevfileIndex int
	Score        int
}

type DevfileType struct {
	Name        string
	Language    string
	ProjectType string
	Tags        []string
	Versions    []Version
}

// EnvVar represents an environment variable with a name and a corresponding value.
type EnvVar struct {
	// Name is the name of the environment variable.
	Name string

	// Value is the value associated with the environment variable.
	Value string
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

type MicronautApplicationProps struct {
	Micronaut struct {
		Server struct {
			Port int `yaml:"port,omitempty"`
			SSL  struct {
				Enabled bool `yaml:"enabled,omitempty"`
				Port    int  `yaml:"port,omitempty"`
			} `yaml:"ssl,omitempty"`
		} `yaml:"server,omitempty"`
	} `yaml:"micronaut,omitempty"`
}

type OpenLibertyServerXml struct {
	HttpEndpoint struct {
		HttpPort  string `xml:"httpPort,attr"`
		HttpsPort string `xml:"httpsPort,attr"`
	} `xml:"httpEndpoint"`
}

type PortDetectionAlgorithm int

type PortMatchRule struct {
	Regex     *regexp.Regexp
	ToReplace string
}

type PortMatchRules struct {
	MatchIndexRegexes []PortMatchRule
	MatchRegexes      []PortMatchSubRule
}

type PortMatchSubRule struct {
	Regex    *regexp.Regexp
	SubRegex *regexp.Regexp
}

type QuarkusApplicationYaml struct {
	Quarkus QuarkusHttp `yaml:"quarkus,omitempty"`
}

type QuarkusHttp struct {
	Http QuarkusHttpPort `yaml:"http,omitempty"`
}

type QuarkusHttpPort struct {
	Port             int    `yaml:"port,omitempty"`
	InsecureRequests string `yaml:"insecure-requests,omitempty"`
	SSLPort          int    `yaml:"ssl-port,omitempty"`
}

type SpringApplicationProsServer struct {
	Server struct {
		Port int `yaml:"port,omitempty"`
		Http struct {
			Port int `yaml:"port,omitempty"`
		} `yaml:"http,omitempty"`
	} `yaml:"server,omitempty"`
}

type Version struct {
	SchemaVersion string
	Default       bool
	Version       string
}

type VertxConf struct {
	Port         int                `json:"http.port,omitempty"`
	ServerConfig VertexServerConfig `json:"http.server,omitempty"`
}

type VertexServerConfig struct {
	Port int `json:"http.server.port,omitempty"`
}
