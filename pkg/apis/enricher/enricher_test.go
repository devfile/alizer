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

// Package enricher implements functions that detect the name and ports of a component.
// Uses three general strategies: Dockerfile, Compose, and Source.
// Dockerfile consists of using a dockerfile to extract information.
// Compose consists of using a compose file to extract information.
// Source consists of searching for specific statements of function invocations inside the source code.
package enricher

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetEnvVarsFromDockerFile(t *testing.T) {
	type args struct {
		root string
	}
	tests := []struct {
		name    string
		args    args
		want    []EnvVar
		wantErr bool
	}{
		{
			name: "case 1: dockerfile project without env var",
			args: args{
				root: "../../../resources/projects/dockerfile",
			},
		},
		{
			name: "case 2: dockerfile project with env var",
			args: args{
				root: "../../../resources/projects/dockerfile-with-port-env-var",
			},
			want: []EnvVar{
				{
					Name:  "PORT",
					Value: "11001",
				},
				{
					Name:  "ANOTHER_VAR",
					Value: "another_value",
				},
			},
		},
		{
			name: "case 3: not found project",
			args: args{
				root: "../../../resources/projects/not-existing",
			},
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetEnvVarsFromDockerFile(tt.args.root)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEnvVarsFromDockerFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEnvVarsFromDockerFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLocations(t *testing.T) {
	type args struct {
		root string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "case 1: one level",
			args: args{root: "../../../resources/projects/dockerfile"},
			want: []string{"Dockerfile", "Containerfile"},
		}, {
			name: "case 2: two levels",
			args: args{root: "../../../resources/projects/dockerfile-nested"},
			want: []string{"Dockerfile", "Containerfile", "docker/Dockerfile", "docker/Containerfile"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLocations(tt.args.root); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLocations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPortsFromReader(t *testing.T) {
	tests := []struct {
		name string
		path string
		want []int
	}{
		{
			name: "case 1: dockerfile with ports",
			path: "../../../resources/projects/dockerfile/Dockerfile",
			want: []int{8085},
		},
		{
			name: "case 2: dockerfile without ports",
			path: "../../../resources/projects/dockerfile-no-port/Dockerfile",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanFilePath := filepath.Clean(tt.path)
			file, err := os.Open(cleanFilePath)
			if err != nil {
				t.Errorf("error: %s", err)
			}
			if got := getPortsFromReader(file); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPortsFromReader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_upsertEnvVar(t *testing.T) {
	testEnvVars := []EnvVar{
		{
			Name:  "var1",
			Value: "value1",
		},
	}
	type args struct {
		envVars []EnvVar
		envVar  EnvVar
	}
	tests := []struct {
		name string
		args args
		want []EnvVar
	}{
		{
			name: "case 1: insert",
			args: args{
				envVars: testEnvVars,
				envVar:  EnvVar{Name: "var2", Value: "value2"},
			},
			want: []EnvVar{
				{
					Name:  "var1",
					Value: "value1",
				},
				{
					Name:  "var2",
					Value: "value2",
				},
			},
		},
		{
			name: "case 2: update",
			args: args{
				envVars: testEnvVars,
				envVar:  EnvVar{Name: "var1", Value: "value2"},
			},
			want: []EnvVar{
				{
					Name:  "var1",
					Value: "value2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := upsertEnvVar(tt.args.envVars, tt.args.envVar); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("upsertEnvVar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEnvVarsFromReader(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    []EnvVar
		wantErr bool
	}{
		{
			name: "case 1: dockerfile project without env var",
			path: "../../../resources/projects/dockerfile/Dockerfile",
		}, {
			name: "case 2: dockerfile project with env var",
			path: "../../../resources/projects/dockerfile-with-port-env-var/Dockerfile",
			want: []EnvVar{
				{
					Name:  "PORT",
					Value: "11001",
				},
				{
					Name:  "ANOTHER_VAR",
					Value: "another_value",
				},
			},
		},
		{
			name:    "case 3: not found project",
			path:    "../../../resources/projects/not-existing/Dockerfile",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanFilePath := filepath.Clean(tt.path)
			file, err := os.Open(cleanFilePath)
			if err != nil && !tt.wantErr {
				t.Errorf("error: %s", err)
			}
			got, err := getEnvVarsFromReader(file)
			if (err != nil) != tt.wantErr {
				t.Errorf("getEnvVarsFromReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEnvVarsFromReader() = %v, want %v", got, tt.want)
			}
		})
	}
}
