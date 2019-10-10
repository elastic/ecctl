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

package util

import (
	"reflect"
	"testing"
)

type ReapplyParams struct {
	ID                   string
	HidePlan             bool `kebabcase:"doesn't print the plan before reapplying"`
	Default              bool `kebabcase:"overwrites the strategy to the default one"`
	Rolling              bool `kebabcase:"overwrites the strategy to rolling"`
	GrowAndShrink        bool `kebabcase:"overwrites the strategy to grow and shrink"`
	RollingGrowAndShrink bool `kebabcase:"overwrites the strategy to rolling grow and shrink (one at a time)"`
	RollingAll           bool `kebabcase:"overwrites the strategy to apply the change in all the instances at a time (causes downtime)"`
	Reallocate           bool `kebabcase:"forces creation of new instances"`
	ExtendedMaintenance  bool `kebabcase:"stops routing to the cluster instances after the plan has been applied"`
	OverrideFailsafe     bool `kebabcase:"overrides failsafe at the constructor level that prevent bad things to happen"`
}

func TestReapplyParamsFieldsOfStruct(t *testing.T) {
	tests := []struct {
		p    interface{}
		name string
		want map[string]string
	}{
		{
			name: "returns all of the struct properties kebabcased",
			p:    ReapplyParams{},
			want: map[string]string{
				"hide-plan":               "doesn't print the plan before reapplying",
				"default":                 "overwrites the strategy to the default one",
				"rolling":                 "overwrites the strategy to rolling",
				"grow-and-shrink":         "overwrites the strategy to grow and shrink",
				"rolling-grow-and-shrink": "overwrites the strategy to rolling grow and shrink (one at a time)",
				"rolling-all":             "overwrites the strategy to apply the change in all the instances at a time (causes downtime)",
				"reallocate":              "forces creation of new instances",
				"extended-maintenance":    "stops routing to the cluster instances after the plan has been applied",
				"override-failsafe":       "overrides failsafe at the constructor level that prevent bad things to happen",
			},
		},
		{
			name: "returns all of the pointer to struct properties kebabcased",
			p:    new(ReapplyParams),
			want: map[string]string{
				"hide-plan":               "doesn't print the plan before reapplying",
				"default":                 "overwrites the strategy to the default one",
				"rolling":                 "overwrites the strategy to rolling",
				"grow-and-shrink":         "overwrites the strategy to grow and shrink",
				"rolling-grow-and-shrink": "overwrites the strategy to rolling grow and shrink (one at a time)",
				"rolling-all":             "overwrites the strategy to apply the change in all the instances at a time (causes downtime)",
				"reallocate":              "forces creation of new instances",
				"extended-maintenance":    "stops routing to the cluster instances after the plan has been applied",
				"override-failsafe":       "overrides failsafe at the constructor level that prevent bad things to happen",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FieldsOfStruct(tt.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReapplyParams.FieldsOfStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReapplyParamsSet(t *testing.T) {
	type args struct {
		p     interface{}
		key   string
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want *ReapplyParams
	}{
		{
			name: "Set hide-plan correctly",
			args: args{
				p:     new(ReapplyParams),
				key:   "hide-plan",
				value: true,
			},
			want: &ReapplyParams{
				HidePlan: true,
			},
		},
		{
			name: "Set extended-maintenance correctly",
			args: args{
				p:     new(ReapplyParams),
				key:   "extended-maintenance",
				value: true,
			},
			want: &ReapplyParams{
				ExtendedMaintenance: true,
			},
		},
		{
			name: "Set properties correctly",
			args: args{
				p: &ReapplyParams{
					ExtendedMaintenance: false,
				},
				key:   "extended-maintenance",
				value: true,
			},
			want: &ReapplyParams{
				ExtendedMaintenance: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Set(tt.args.p, tt.args.key, tt.args.value)
			if !reflect.DeepEqual(tt.args.p, tt.want) {
				t.Errorf("ReapplyParams.FieldsOfStruct() = %v, want %v", tt.args.p, tt.want)
			}
		})
	}
}
