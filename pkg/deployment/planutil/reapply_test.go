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

package planutil

import (
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/models"
)

func TestReapplyParamsValidate(t *testing.T) {
	type fields struct {
		ID                   string
		HidePlan             bool
		Default              bool
		Rolling              bool
		GrowAndShrink        bool
		RollingGrowAndShrink bool
		RollingAll           bool
		Reallocate           bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"ID is invalid",
			fields{
				ID:                   "123131",
				Default:              false,
				Rolling:              false,
				GrowAndShrink:        false,
				RollingGrowAndShrink: false,
				RollingAll:           false,
			},
			true,
		},
		{
			"Default Strategy is set to true",
			fields{
				ID:                   "ee44fcf653f4283cecaa6bbc23c66c4e",
				Default:              true,
				Rolling:              false,
				GrowAndShrink:        false,
				RollingGrowAndShrink: false,
				RollingAll:           false,
			},
			false,
		},
		{
			"Default Strategy is set to false",
			fields{
				ID:                   "ee44fcf653f4283cecaa6bbc23c66c4e",
				Default:              false,
				Rolling:              false,
				GrowAndShrink:        false,
				RollingGrowAndShrink: false,
				RollingAll:           false,
			},
			false,
		},
		{
			"Rolling Strategy is valid",
			fields{
				ID:                   "ee44fcf653f4283cecaa6bbc23c66c4e",
				Default:              false,
				Rolling:              true,
				GrowAndShrink:        false,
				RollingGrowAndShrink: false,
				RollingAll:           false,
			},
			false,
		},
		{
			"Grow&Shrink Strategy is valid",
			fields{
				ID:                   "ee44fcf653f4283cecaa6bbc23c66c4e",
				Default:              false,
				Rolling:              false,
				GrowAndShrink:        true,
				RollingGrowAndShrink: false,
				RollingAll:           false,
			},
			false,
		},
		{
			"RollingGrow&Shrink Strategy is valid",
			fields{
				ID:                   "ee44fcf653f4283cecaa6bbc23c66c4e",
				Default:              false,
				Rolling:              false,
				GrowAndShrink:        false,
				RollingGrowAndShrink: true,
				RollingAll:           false,
			},
			false,
		},
		{
			"RollingAll Strategy is valid",
			fields{
				ID:                   "ee44fcf653f4283cecaa6bbc23c66c4e",
				Default:              false,
				Rolling:              false,
				GrowAndShrink:        false,
				RollingGrowAndShrink: false,
				RollingAll:           true,
			},
			false,
		},
		{
			"Strategy is invalid",
			fields{
				ID:                   "ee44fcf653f4283cecaa6bbc23c66c4e",
				Default:              false,
				Rolling:              true,
				GrowAndShrink:        true,
				RollingGrowAndShrink: false,
				RollingAll:           false,
			},
			true,
		},
		{
			"Strategy is invalid",
			fields{
				ID:                   "ee44fcf653f4283cecaa6bbc23c66c4e",
				Default:              false,
				Rolling:              true,
				GrowAndShrink:        true,
				RollingGrowAndShrink: true,
				RollingAll:           false,
			},
			true,
		},
		{
			"Strategy is invalid",
			fields{
				ID:                   "ee44fcf653f4283cecaa6bbc23c66c4e",
				Default:              false,
				Rolling:              true,
				GrowAndShrink:        false,
				RollingGrowAndShrink: true,
				RollingAll:           false,
			},
			true,
		},
		{
			"Strategy is invalid GrowAndShrink AND RollingGrowAndShrink",
			fields{
				ID:                   "ee44fcf653f4283cecaa6bbc23c66c4e",
				Default:              false,
				Rolling:              false,
				GrowAndShrink:        true,
				RollingGrowAndShrink: true,
				RollingAll:           false,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ReapplyParams{
				ID:                   tt.fields.ID,
				HidePlan:             tt.fields.HidePlan,
				Reallocate:           tt.fields.Reallocate,
				Default:              tt.fields.Default,
				Rolling:              tt.fields.Rolling,
				GrowAndShrink:        tt.fields.GrowAndShrink,
				RollingGrowAndShrink: tt.fields.RollingGrowAndShrink,
				RollingAll:           tt.fields.RollingAll,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ReapplyParams.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReapplyParamsStrategy(t *testing.T) {
	type fields struct {
		ID                   string
		HidePlan             bool
		Reallocate           bool
		Default              bool
		Rolling              bool
		GrowAndShrink        bool
		RollingGrowAndShrink bool
		RollingAll           bool
	}
	tests := []struct {
		name   string
		fields fields
		want   *models.PlanStrategy
	}{
		{
			"Strategy is nil when no strategy is specified",
			fields{
				ID: "ee44fcf653f4283cecaa6bbc23c66c4e",
			},
			nil,
		},
		{
			"Strategy is default when flag is specified",
			fields{
				ID:                   "ee44fcf653f4283cecaa6bbc23c66c4e",
				Default:              true,
				Rolling:              false,
				GrowAndShrink:        false,
				RollingGrowAndShrink: false,
			},
			DefaultPlanStrategy,
		},
		{
			"Strategy is rolling when flag is specified",
			fields{
				ID:      "ee44fcf653f4283cecaa6bbc23c66c4e",
				Rolling: true,
			},
			RollingByNameStrategy,
		},
		{
			"Strategy is grow and shrink when flag is specified",
			fields{
				ID:            "ee44fcf653f4283cecaa6bbc23c66c4e",
				GrowAndShrink: true,
			},
			GrowAndShrinkStrategy,
		},
		{
			"Strategy is rolling grow and shrink when flag is specified",
			fields{
				ID:                   "ee44fcf653f4283cecaa6bbc23c66c4e",
				RollingGrowAndShrink: true,
			},
			RollingGrowAndShrinkStrategy,
		},
		{
			"Strategy is rolling all when flag is specified",
			fields{
				ID:         "ee44fcf653f4283cecaa6bbc23c66c4e",
				RollingAll: true,
			},
			MajorUpgradeStrategy,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ReapplyParams{
				ID:                   tt.fields.ID,
				HidePlan:             tt.fields.HidePlan,
				Reallocate:           tt.fields.Reallocate,
				Default:              tt.fields.Default,
				Rolling:              tt.fields.Rolling,
				GrowAndShrink:        tt.fields.GrowAndShrink,
				RollingGrowAndShrink: tt.fields.RollingGrowAndShrink,
				RollingAll:           tt.fields.RollingAll,
			}
			if got := p.Strategy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReapplyParams.strategy() = %v, want %v", got, tt.want)
			}
		})
	}
}
