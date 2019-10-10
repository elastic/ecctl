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

// Package formatter contains the logic to
// template the CLI output to a specific format
package formatter

import "io"

// Formatter models the formatter actions
type Formatter interface {
	// Format is used from the cmd, it receives the parent and the data to be
	// formatted. The path usually has an effect on how the data will end up
	// being represented.
	Format(path string, data interface{}) error
	// Name obtains the name of the Formatter
	Name() string
}

// New initializes a ChainFormatter defaulting to the
// JSON Formatter if no known formatter is passed, else
// It will use the mentioned formatter and use the JSON
// Formatter to fallback when parsing the output fails.
func New(o io.Writer, name string) *Chain {
	fallbackFormatter := NewChain(NewJSON(o))
	if name == "text" {
		fallbackFormatter = NewChain(
			NewText(&TextConfig{Output: o}),
		).Add(fallbackFormatter)
	}

	return fallbackFormatter
}

// Chainer implements the Single Chain of Responsibility
// Pattern, to be able to chain off multiple Formatters
type Chainer interface {
	Formatter
	// Add adds a the next ChainFormatter to the chain
	Add(Chainer) *Chain
}

// NewChain wraps a formatter into a Chain
func NewChain(f Formatter) *Chain { return &Chain{formatter: f} }

// Chain implements the Chain interface
//
// Allows the formatting responsibility to become a chain.
// It will cascade when formatter.Format returns an error
// and the next property contains another actionable chain
// (next != nil). When the Format is handled without error
// it will return immediately
type Chain struct {
	next      Chainer
	formatter Formatter
}

// Format wraps the Format method of Formatter going down the chain
// when formatter.Format returns an error and the next property
// contains another actionable chain (next != nil). When the Format
// is handled without error it will return immediately
func (f Chain) Format(path string, data interface{}) error {
	err := f.formatter.Format(path, data)
	if err == nil {
		return nil
	}

	if f.next != nil {
		return f.next.Format(path, data)
	}
	return err
}

// Add adds a the next ChainFormatter to the chain
func (f *Chain) Add(next Chainer) *Chain {
	f.next = next
	return f
}

// Name obtains the name of the Formatter
func (f *Chain) Name() string { return f.formatter.Name() }
