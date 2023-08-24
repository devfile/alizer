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
				root: "../../../resources/projects/dockerfile",
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
