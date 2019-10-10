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
	"testing"
)

func TestMin(t *testing.T) {
	type args struct {
		a int8
		b int8
	}
	tests := []struct {
		name string
		args args
		want int8
	}{
		{
			"Min returns the smaller number when the first number is smaller",
			args{
				1, 10,
			},
			1,
		},
		{
			"Min returns the smaller number when the second number is smaller",
			args{
				-50, -100,
			},
			-100,
		},
		{
			"Min returns the smaller number when both numbers are equal",
			args{
				20, 20,
			},
			20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Min(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}
