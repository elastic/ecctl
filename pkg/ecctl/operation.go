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

import "bytes"

var instance *operationCommentator

type operationCommentator struct {
	extraMessage string
}

// Commentator contains the actions for any commentator
type Commentator interface {
	// Set sets an extra message in the commentator
	Set(m string)
	// Message formats the message so it returns any extra message
	// after the message it receives.
	Message(m string) string
}

// GetOperationInstance obtains the operation Commenter singleton
func GetOperationInstance() Commentator {
	if instance == nil {
		instance = &operationCommentator{}
	}
	return instance
}

// Set sets an extra message in the commentator
func (c *operationCommentator) Set(m string) {
	c.extraMessage = m
}

// Message formats the message so it returns any extra message
// after the message it receives.
func (c *operationCommentator) Message(m string) string {
	if c.extraMessage == "" {
		return m
	}

	// Using a buffer for performance reasons
	var finalMessageBuf bytes.Buffer
	finalMessageBuf.WriteString(m)
	if string(c.extraMessage[0]) != " " {
		finalMessageBuf.WriteString(" ")
	}

	finalMessageBuf.WriteString(c.extraMessage)
	return finalMessageBuf.String()
}
