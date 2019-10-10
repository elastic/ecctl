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

	"github.com/asaskevich/govalidator"
)

const kebabcase = "kebabcase"

// FieldsOfStruct returns a list of the kebabcased struct property names that
// have the tag `kebabcase:"description"` set. The map that will be returned
// will have the form of:
// "skip-data-migration": "description"
func FieldsOfStruct(p interface{}) map[string]string {
	var reflectedStruct = reflect.ValueOf(p)
	if reflect.Ptr == reflectedStruct.Kind() {
		reflectedStruct = reflectedStruct.Elem()
	}

	var field = reflectedStruct.Type().Field
	var fields = make(map[string]string)
	for i := 0; i < reflectedStruct.NumField(); i++ {
		if t := field(i).Tag.Get(kebabcase); t != "" {
			var key = UnderscoreToDashes(field(i).Name)
			fields[key] = t
		}
	}
	return fields
}

// Set will populate a property from its dashed version to the specified value
// i.e. passing `skip-snapshot` as a key will populate the `SkipSnapshot` prop.
func Set(p interface{}, key string, value interface{}) {
	if reflect.Ptr != reflect.ValueOf(p).Kind() {
		panic("must receive a struct reference not a struct value")
	}

	var prop = reflect.ValueOf(p).Elem().FieldByName(
		govalidator.UnderscoreToCamelCase(
			DashesToUnderscore(key),
		),
	)

	if prop.IsValid() && prop.CanSet() {
		prop.Set(reflect.ValueOf(value))
	}
}
