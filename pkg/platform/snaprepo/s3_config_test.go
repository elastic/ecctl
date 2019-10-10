// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package snaprepo

import (
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseS3Config(t *testing.T) {
	var invalidReaderContent = `such invalid {}content`
	var simpleJSONConfig = `{"region": "us-east-1"}`
	var simpleYAMLConfig = `
region: us-east-1
bucket: mybucket
`[1:]

	type args struct {
		input io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    S3Config
		wantErr bool
	}{
		{
			name: "returns an error if the reader is nil",
			args: args{
				input: nil,
			},
			wantErr: true,
		},
		{
			name: "returns an error if the reader's content is not a valid yaml|json",
			args: args{
				input: strings.NewReader(invalidReaderContent),
			},
			wantErr: true,
		},
		{
			name: "Correctly reads a JSON formatted reader",
			args: args{
				input: strings.NewReader(simpleJSONConfig),
			},
			want: S3Config{
				Region: "us-east-1",
			},
			wantErr: false,
		},
		{
			name: "Correctly reads a YAML formatted reader",
			args: args{
				input: strings.NewReader(simpleYAMLConfig),
			},
			want: S3Config{
				Region: "us-east-1",
				Bucket: "mybucket",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseS3Config(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseS3Config() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseS3Config() = %v, want %v", got, tt.want)
			}
		})
	}
}

func parseDuration(t *testing.T, s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestS3Config_Validate(t *testing.T) {
	type fields struct {
		Region               string
		Bucket               string
		AccessKey            string
		SecretKey            string
		BasePath             string
		ChunkSize            string
		CannedACL            string
		StorageClass         string
		Endpoint             string
		Protocol             string
		Timeout              time.Duration
		MaxRetries           int
		Compress             bool
		ServerSideEncryption bool
		ThrottleRetries      bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Succeeds when a simple config is specified",
			fields: fields{
				Region:    "us-east-1",
				Bucket:    "mybucket",
				AccessKey: "myaccesskey",
				SecretKey: "mysupersecretkey",
			},
			wantErr: false,
		},
		{
			name:    "Fails when a simple config is invalid",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "Succeeds when a complex config is specified",
			fields: fields{
				Region:               "us-east-1",
				Bucket:               "mybucket",
				AccessKey:            "myaccesskey",
				SecretKey:            "mysupersecretkey",
				BasePath:             "apath",
				Compress:             true,
				ServerSideEncryption: true,
				Timeout:              parseDuration(t, "60s"),
				CannedACL:            "private",
				StorageClass:         "standard",
				Protocol:             "http",
			},
			wantErr: false,
		},
		{
			name: "Fails when when a complex config is invalid",
			fields: fields{
				Region:               "us-east-1",
				Bucket:               "mybucket",
				AccessKey:            "myaccesskey",
				SecretKey:            "mysupersecretkey",
				BasePath:             "apath",
				Compress:             true,
				ServerSideEncryption: true,
				Endpoint:             "a very invalid endpoint",
				CannedACL:            "an invalid value",
				StorageClass:         "another invalid value",
				Protocol:             "yet another invalid value",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := S3Config{
				Region:               tt.fields.Region,
				Bucket:               tt.fields.Bucket,
				AccessKey:            tt.fields.AccessKey,
				SecretKey:            tt.fields.SecretKey,
				BasePath:             tt.fields.BasePath,
				Compress:             tt.fields.Compress,
				ServerSideEncryption: tt.fields.ServerSideEncryption,
				ChunkSize:            tt.fields.ChunkSize,
				CannedACL:            tt.fields.CannedACL,
				StorageClass:         tt.fields.StorageClass,
				Endpoint:             tt.fields.Endpoint,
				Protocol:             tt.fields.Protocol,
				Timeout:              tt.fields.Timeout,
				MaxRetries:           tt.fields.MaxRetries,
				ThrottleRetries:      tt.fields.ThrottleRetries,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("S3Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
