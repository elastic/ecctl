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

// text.go contains the cli templates to
// be able to set the output to `text`

package formatter

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"text/template"

	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/pkg/formatter/templates"
)

func TestText_format(t *testing.T) {
	type abc struct {
		A string
		B string
		C string
	}

	type fields struct {
		output          io.Writer
		templater       *template.Template
		override        bool
		assetLoaderFunc AssetLoaderFunc
	}
	type args struct {
		data interface{}
		text string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    string
	}{
		{
			name: "Passing a struct that complies the format succeeds",
			fields: fields{
				assetLoaderFunc: templates.Asset,
				output:          new(bytes.Buffer),
				templater:       template.New("text"),
			},
			args: args{
				data: &abc{"a", "b", "c"},
				text: `{{define "default"}}{{.A}}  {{.B}} {{.C}}{{end}}`,
			},
			wantErr: false,
			want:    "a  b c",
		},
		{
			name: "Passing a struct that does not comply the format fails",
			fields: fields{
				assetLoaderFunc: templates.Asset,
				output:          new(bytes.Buffer),
				templater:       template.New("text"),
			},
			args: args{
				&abc{"a", "b", "c"},
				`{{define "default"}}{{.W}}  {{.S}} {{.Z}}{{end}}`,
			},
			wantErr: true,
			want:    "",
		},
		{
			name: `If override is set, the "override" template name is used`,
			fields: fields{
				assetLoaderFunc: templates.Asset,
				output:          new(bytes.Buffer),
				templater:       template.New("text"),
				override:        true,
			},
			args: args{
				&abc{"a", "b", "c"},
				`{{define "default"}}{{.A}}  {{.B}} {{.C}}{{end}}{{define "override"}}{{.A}}{{.B}} {{.C}}{{end}}`,
			},
			wantErr: false,
			want:    "ab c",
		},
		{
			name: `Failure to parse template returns an error`,
			fields: fields{
				assetLoaderFunc: templates.Asset,
				output:          new(bytes.Buffer),
				templater:       template.New("text"),
				override:        true,
			},
			args: args{
				&abc{"a", "b", "c"},
				`{{Invalid Template}}`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Text{
				output:          tt.fields.output,
				templater:       tt.fields.templater,
				override:        tt.fields.override,
				assetLoaderFunc: tt.fields.assetLoaderFunc,
			}
			if err := f.format(tt.args.text, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Text.format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := tt.fields.output.(*bytes.Buffer).String()
			if !strings.EqualFold(got, tt.want) {
				t.Errorf("Text.format() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatBytes(t *testing.T) {
	type args struct {
		cap   int32
		human bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Format empty",
			args{
				int32(0),
				true,
			},
			"-",
		},
		{
			"Format Megabytes",
			args{
				int32(768),
				true,
			},
			"768  MB",
		},
		{
			"Format Gigabytes",
			args{
				int32(2048),
				true,
			},
			"2    GB",
		},
		{
			"Format Terabytes",
			args{
				int32(2048000),
				true,
			},
			"2    TB",
		},
		{
			"Format NonHuman",
			args{
				int32(2048000),
				false,
			},
			"2048000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatBytes(tt.args.cap, tt.args.human); got != tt.want {
				t.Errorf("formatBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestText_FormatAllocators(t *testing.T) {
	var overrideBuffer = new(bytes.Buffer)

	type fields struct {
		o               *bytes.Buffer
		t               *template.Template
		assetLoaderFunc AssetLoaderFunc
		funcs           template.FuncMap
		// overrides the default template. Still relies on the template asset
		// to loop over the properties
		override bool
		// fallback base template to use when a file asset is not found in
		// "path" unless specified `defaultOverrideFormat` is
		// used.
		fallback string
	}
	type args struct {
		path string
		data interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    string
	}{
		{
			"format allocator list",
			fields{
				assetLoaderFunc: templates.Asset,
				o:               &bytes.Buffer{},
				t:               template.New("text").Funcs(defaultTemplateFuncs),
			},
			args{
				path: "allocator/list",
				data: &models.AllocatorOverview{
					Zones: []*models.AllocatorZoneInfo{
						{
							Allocators: []*models.AllocatorInfo{
								{
									AllocatorID: ec.String("MyAllocatorID"),
									Capacity: &models.AllocatorCapacity{
										Memory: &models.AllocatorCapacityMemory{
											Total: ec.Int32(4096),
											Used:  ec.Int32(2048),
										},
									},
									Status: &models.AllocatorHealthStatus{
										Connected:       ec.Bool(true),
										MaintenanceMode: ec.Bool(false),
									},
									ZoneID: ec.String("zone1"),
									HostIP: ec.String("1.1.1.1"),
								},
							},
							ZoneID: ec.String("zone1"),
						},
						{
							Allocators: []*models.AllocatorInfo{
								{
									AllocatorID: ec.String("MyAllocatorID"),
									Capacity: &models.AllocatorCapacity{
										Memory: &models.AllocatorCapacityMemory{
											Total: ec.Int32(4096),
											Used:  ec.Int32(2048),
										},
									},
									Status: &models.AllocatorHealthStatus{
										Connected:       ec.Bool(true),
										MaintenanceMode: ec.Bool(false),
									},
									ZoneID: ec.String("zone2"),
									HostIP: ec.String("1.1.1.2"),
									Instances: []*models.AllocatedInstanceStatus{
										{},
									},
								},
							},
							ZoneID: ec.String("zone2"),
						},
					},
				},
			},
			false,
			`
ZONE    ALLOCATOR ID    HOST IP   CAPACITY   FREE      INSTANCES   CONNECTED   MAINTENANCE
zone1   MyAllocatorID   1.1.1.1   4    GB    2    GB   0           true        false
zone2   MyAllocatorID   1.1.1.2   4    GB    2    GB   1           true        false
`[1:],
		},
		{
			"format allocator list, with a disconnected allocator with Instances",
			fields{
				assetLoaderFunc: templates.Asset,
				o:               &bytes.Buffer{},
				t:               template.New("text").Funcs(defaultTemplateFuncs),
			},
			args{
				path: "allocator/list",
				data: &models.AllocatorOverview{
					Zones: []*models.AllocatorZoneInfo{
						{
							Allocators: []*models.AllocatorInfo{
								{
									AllocatorID: ec.String("MyAllocatorID"),
									Capacity: &models.AllocatorCapacity{
										Memory: &models.AllocatorCapacityMemory{
											Total: ec.Int32(4096),
											Used:  ec.Int32(2048),
										},
									},
									Status: &models.AllocatorHealthStatus{
										Connected:       ec.Bool(true),
										MaintenanceMode: ec.Bool(false),
									},
									ZoneID: ec.String("zone1"),
									HostIP: ec.String("1.1.1.1"),
								},
							},
							ZoneID: ec.String("zone1"),
						},
						{
							Allocators: []*models.AllocatorInfo{
								{
									AllocatorID: ec.String("MyAllocatorID"),
									Capacity: &models.AllocatorCapacity{
										Memory: &models.AllocatorCapacityMemory{
											Total: ec.Int32(4096),
											Used:  ec.Int32(2048),
										},
									},
									Status: &models.AllocatorHealthStatus{
										Connected:       ec.Bool(false),
										MaintenanceMode: ec.Bool(false),
									},
									ZoneID: ec.String("zone2"),
									HostIP: ec.String("1.1.1.2"),
									Instances: []*models.AllocatedInstanceStatus{
										{
											ClusterID:    ec.String("1234"),
											ClusterName:  "Instance",
											InstanceName: ec.String("instance-0"),
											ClusterType:  ec.String("elasticsearch"),
											NodeMemory:   ec.Int32(1024),
										},
										{
											ClusterID:    ec.String("1235"),
											ClusterName:  "Instance",
											InstanceName: ec.String("instance-0"),
											ClusterType:  ec.String("elasticsearch"),
											NodeMemory:   ec.Int32(1024),
										},
									},
								},
							},
							ZoneID: ec.String("zone2"),
						},
					},
				},
			},
			false,
			`
ZONE    ALLOCATOR ID    HOST IP   CAPACITY   FREE      INSTANCES   CONNECTED   MAINTENANCE
zone1   MyAllocatorID   1.1.1.1   4    GB    2    GB   0           true        false
zone2   MyAllocatorID   1.1.1.2   4    GB    2    GB   2           false       false
`[1:],
		},
		{
			"format allocator where a command has no available template from the assets",
			fields{
				assetLoaderFunc: templates.Asset,
				o:               overrideBuffer,
				t:               template.New("text").Funcs(defaultTemplateFuncs),
				override:        true,
				fallback:        defaultOverrideFormat,
				funcs:           newOverridenTemplateFuncs(overrideBuffer, "{{ .AllocatorID }}"),
			},
			args{
				path: "command",
				data: &models.AllocatorInfo{
					AllocatorID: ec.String("MyAllocatorID"),
				},
			},
			false,
			"MyAllocatorID\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Text{
				output:          tt.fields.o,
				templater:       tt.fields.t,
				override:        tt.fields.override,
				fallback:        tt.fields.fallback,
				funcs:           tt.fields.funcs,
				assetLoaderFunc: tt.fields.assetLoaderFunc,
			}
			if err := f.Format(tt.args.path, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Text.Format() error = \n%v, wantErr \n%v", err, tt.wantErr)
				return
			}
			if tt.fields.o.String() != tt.want {
				t.Errorf("Text.Format() got = \n%v, want \n%v", tt.fields.o.String(), tt.want)
			}
		})
	}
}

func newOverridenTemplateFuncs(o io.Writer, t string) template.FuncMap {
	var overridenTemplateFuncs = defaultTemplateFuncs
	overridenTemplateFuncs["executeTemplate"] = executeTemplateFunc(o, t)
	return overridenTemplateFuncs
}

func TestNewText(t *testing.T) {
	type args struct {
		c *TextConfig
	}
	tests := []struct {
		name string
		args args
		want *Text
	}{
		{
			name: "New text formatter without any override uses the default funcs",
			args: args{
				c: &TextConfig{
					Output: new(bytes.Buffer),
				},
			},
			want: &Text{
				output:    new(bytes.Buffer),
				templater: template.New("text"),
				funcs:     defaultTemplateFuncs,
				fallback:  defaultOverrideFormat,
				padding:   defaultPadding,
			},
		},
		{
			name: "New text formatter with an override adds 2 more functions",
			args: args{
				c: &TextConfig{
					Output:   new(bytes.Buffer),
					Override: "{{ .Name }}",
				},
			},
			want: &Text{
				output:    new(bytes.Buffer),
				templater: template.New("text"),
				funcs:     newOverridenTemplateFuncs(os.Stdout, "{{ .Name }}"),
				override:  true,
				fallback:  defaultOverrideFormat,
				padding:   defaultPadding,
			},
		},
		{
			name: "New text formatter with a fallback overide overrides defaultOverrideFormat",
			args: args{
				c: &TextConfig{
					Output:   new(bytes.Buffer),
					Override: "{{ .Name }}",
					Fallback: "{{ range . }}{{executeTemplate .}}{{end}}",
				},
			},
			want: &Text{
				output:    new(bytes.Buffer),
				templater: template.New("text"),
				funcs:     newOverridenTemplateFuncs(os.Stdout, "{{ .Name }}"),
				override:  true,
				fallback:  "{{ range . }}{{executeTemplate .}}{{end}}",
				padding:   defaultPadding,
			},
		},
		{
			name: "New text formatter with a fallback overide overrides defaultOverrideFormat",
			args: args{
				c: &TextConfig{
					Output:   new(bytes.Buffer),
					Override: "{{ .Name }}",
					Fallback: "{{ range . }}{{executeTemplate .}}{{end}}",
				},
			},
			want: &Text{
				output:    new(bytes.Buffer),
				templater: template.New("text"),
				funcs:     newOverridenTemplateFuncs(os.Stdout, "{{ .Name }}"),
				override:  true,
				fallback:  "{{ range . }}{{executeTemplate .}}{{end}}",
				padding:   defaultPadding,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewText(tt.args.c)
			// Can't deepequal funcs.
			got.assetLoaderFunc = nil
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewText() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_executeTemplateFunc(t *testing.T) {
	type fields struct {
		templateFormat string
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       string
		wantOutput string
	}{
		{
			name: "Parsing and executing the template succeeds",
			fields: fields{
				templateFormat: "{{ .AllocatorID }}",
			},
			args: args{
				data: &models.AllocatorInfo{
					AllocatorID: ec.String("MyAllocatorID"),
				},
			},
			want:       "MyAllocatorID",
			wantOutput: "",
		},
		{
			name: "Parsing the template succeeds but execution fails",
			fields: fields{
				templateFormat: "{{ .AnUnknownProperty }}",
			},
			args: args{
				data: &models.AllocatorInfo{
					AllocatorID: ec.String("MyAllocatorID"),
				},
			},
			want:       "",
			wantOutput: fmt.Sprintln(`template: template:1:3: executing "template" at <.AnUnknownProperty>: can't evaluate field AnUnknownProperty in type *models.AllocatorInfo`),
		},
		{
			name: "Parsing the template fails",
			fields: fields{
				templateFormat: "{{ WhellThisFailsOnTemplateParse }}",
			},
			args: args{
				data: nil,
			},
			want:       "",
			wantOutput: fmt.Sprintln(`template: template:1: function "WhellThisFailsOnTemplateParse" not defined`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := new(bytes.Buffer)
			if got := executeTemplateFunc(output, tt.fields.templateFormat)(tt.args.data); got != tt.want {
				t.Errorf("Return of executeTemplateFunc() = %v, want %v", got, tt.want)
			}
			if gotOutput := output.String(); gotOutput != tt.wantOutput {
				t.Errorf("Output of executeTemplateFunc() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
