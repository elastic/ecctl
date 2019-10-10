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

package ecctl

import (
	"reflect"
	"testing"
)

func TestOperationCommentator_Message(t *testing.T) {
	type fields struct {
		extraMessage string
	}
	type args struct {
		m string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"Message when there's no extra message",
			fields{
				"",
			},
			args{
				"My message",
			},
			"My message",
		},
		{
			"Message when there's an extra space in the extra message",
			fields{
				" extra part of message",
			},
			args{
				"My message",
			},
			"My message extra part of message",
		},
		{
			"Message when there's an extra message",
			fields{
				"extra part of message",
			},
			args{
				"My message",
			},
			"My message extra part of message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &operationCommentator{
				extraMessage: tt.fields.extraMessage,
			}
			if got := c.Message(tt.args.m); got != tt.want {
				t.Errorf("operationCommentator.Message() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperationCommentatorSet(t *testing.T) {
	type fields struct {
		extraMessage string
	}
	type args struct {
		m string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"Set effectively sets an extra message",
			fields{
				"",
			},
			args{
				"message",
			},
			"message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &operationCommentator{
				extraMessage: tt.fields.extraMessage,
			}
			c.Set(tt.args.m)
			if c.extraMessage != tt.want {
				t.Errorf("operationCommentator.Set() = %v, want %v", c.extraMessage, tt.want)
			}
		})
	}
}

func TestGetOperationInstance(t *testing.T) {
	tests := []struct {
		name string
		want Commentator
	}{
		{
			"Obtains a comentator",
			&operationCommentator{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOperationInstance(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOperationInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkOperationCommentatorMessage(b *testing.B) {
	benchmarks := []struct {
		name         string
		extraMessage string
		message      string
		N            int
	}{
		{"No extra message", "", "message", 1000},
		{"With extra message", "extra message", "message", 1000},
		{"With space in extra message", " extra message", "message", 1000},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			c := new(operationCommentator)
			for n := 0; n < b.N; n++ {
				c.Set(bm.extraMessage)
				c.Message(bm.message)
			}
		})
	}
}
