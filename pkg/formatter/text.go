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

package formatter

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/elastic/ecctl/pkg/formatter/templates"
)

var defaultTemplateFuncs = template.FuncMap{
	"trim":                      strings.TrimSpace,
	"trimToLen":                 trimToLen,
	"toS3TypeConfig":            toS3TypeConfig,
	"rpad":                      rpad,
	"rpadTrim":                  rpadTrim,
	"formatBytes":               formatBytes,
	"tab":                       tab,
	"formatClusterBytes":        formatClusterBytes,
	"substr":                    substr,
	"computeClusterCapacity":    computeClusterCapacity,
	"equal":                     equal,
	"derefInt":                  derefInt,
	"derefBool":                 derefBool,
	"displayAllocator":          displayAllocator,
	"getFailedPlanStepName":     getFailedPlanStepName,
	"getClusterName":            getClusterName,
	"formatTopologyInfo":        formatTopologyInfo,
	"computePlanDuration":       computePlanDuration,
	"getESCurrentOrPendingPlan": getESCurrentOrPendingPlan,
	"centiCentsToCents":         centiCentsToCents,
	"computeApmPlanDuration":    computeApmPlanDuration,
	"getApmFailedPlanStepName":  getApmFailedPlanStepName,
}

const (
	defaultPadding = 3
)

var (
	// defaultOverrideFormat is used when no fallback has been specified in
	// TextConfig
	defaultOverrideFormat = `
{{- define "override" }}{{executeTemplate .}}
{{ end }}`[1:]
)

// AssetLoaderFunc loads and returns the asset for the given name. If an error
// is returned, the asset could not be found or could not be loaded.
type AssetLoaderFunc func(name string) ([]byte, error)

// TextConfig is the configuration to use for a text formatter
type TextConfig struct {
	// Output device where the template execution will be written.
	Output io.Writer
	// Override is the format to be used for the commands that are sent
	Override string
	// Fallback base template to use when a file asset is not found in
	// "text/parent/child.gotmpl" unless specified `defaultOverrideFormat` is
	// used.
	Fallback string
	// Padding represents the number of spaces to use in between the text
	// columns.
	Padding int

	// AssetLoaderFunction is the asset loader function to be used. It must
	// return a byte encoded Go template parseable by template.Template.Parse.
	AssetLoaderFunction AssetLoaderFunc
}

func (c *TextConfig) fillValues() {
	if c.Fallback == "" {
		c.Fallback = defaultOverrideFormat
	}

	if c.AssetLoaderFunction == nil {
		c.AssetLoaderFunction = templates.Asset
	}

	if c.Padding <= 0 {
		c.Padding = defaultPadding
	}
}

// NewText acts as the factory for formatter.Text
func NewText(c *TextConfig) *Text {
	var templateFuncs = defaultTemplateFuncs
	templateFuncs["executeTemplate"] = executeTemplateFunc(c.Output, c.Override)

	c.fillValues()

	return &Text{
		output:          c.Output,
		padding:         c.Padding,
		templater:       template.New("text"),
		funcs:           templateFuncs,
		override:        c.Override != "",
		fallback:        c.Fallback,
		assetLoaderFunc: c.AssetLoaderFunction,
	}
}

// Text formats into text
type Text struct {
	// Output device where the template execution will be written.
	output io.Writer
	// Templater to use
	templater *template.Template
	// funcs are the funtions that are passed to the template on format
	funcs template.FuncMap
	// overrides the default template. Still relies on the template asset
	// to loop over the properties
	override bool
	// fallback base template to use when a file asset is not found in
	// "text/parent/child.gotmpl" unless specified `defaultOverrideFormat` is
	// used.
	fallback string
	// padding represents the number of spaces to use in between the text
	// columns.
	padding int

	assetLoaderFunc AssetLoaderFunc
}

// format formats the data according to the template
// used during struct initialization
func (f *Text) format(text string, data interface{}) error {
	t, err := f.templater.Funcs(f.funcs).Parse(text)
	if err != nil {
		return err
	}

	var name = "default"
	if f.override {
		name = "override"
	}

	w := tabwriter.NewWriter(f.output, 2, 4, 3, ' ', 0)
	defer func() { _ = w.Flush() }()

	return t.ExecuteTemplate(w, name, data)
}

// Name obtains the name of the formatter
func (f *Text) Name() string { return f.templater.Name() }

// Format is used from the cmd for conveniency, it receives a path and the data
// which needs to be used in the template. Uses the name of the formatter as
// the first argument of the template name.
// If the asset is not found, it will default to defaultOverrideFormat
func (f *Text) Format(path string, data interface{}) error {
	if filepath.Ext(path) == "" {
		path += ".gotmpl"
	}

	tformat, err := f.assetLoaderFunc(filepath.Join(f.Name(), path))
	if err != nil {
		tformat = []byte(f.fallback)
	}
	return f.format(string(tformat), data)
}

func executeTemplateFunc(output io.Writer, templateFormat string) func(data interface{}) string {
	return func(data interface{}) string {
		t, err := template.New("template").Funcs(defaultTemplateFuncs).Parse(templateFormat)
		if err != nil {
			_, _ = fmt.Fprintln(output, err)
			return ""
		}
		var b = new(bytes.Buffer)
		if err := t.Execute(b, data); err != nil {
			_, _ = fmt.Fprintln(output, err)
			return ""
		}

		return b.String()
	}
}
