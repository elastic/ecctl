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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

var (
	validStorageClasses = []string{"standard", "reduced_redundancy", "standard_ia"}
	validCannedACLs     = []string{
		"private",
		"public-read",
		"public-read-write",
		"authenticated-read",
		"log-delivery-write",
		"bucket-owner-read",
		"bucket-owner-full-control",
	}
	validProtocols = []string{"http", "https"}
)

var (
	// required setting errors
	errRegionCannotBeEmpty    = errors.New("region cannot be empty")
	errBucketCannotBeEmpty    = errors.New("bucket cannot be empty")
	errAccessKeyCannotBeEmpty = errors.New("access key cannot be empty")
	errSecretKeyCannotBeEmpty = errors.New("secret key cannot be empty")

	// Optional setting errors
	errInvalidStorageClass = fmt.Errorf("storage class must be one of %s", validStorageClasses)
	errInvalidCannedACL    = fmt.Errorf("canned acl must be one of %s", validCannedACLs)
	errInvalidProtocol     = fmt.Errorf("protocol must be one of %s", validProtocols)
	errInvalidEndpoint     = errors.New("endpoint is not valid")

	// Parser errors
	errFailedParsingConfig = errors.New("failed to parse config format")
	errReaderCannotBeNil   = errors.New("reader cannot be nil")
)

// ParseS3Config reads the contents of an io.Reader and tries to parse its
// contents as YAML or JSON, returns an error if parsing fails in both formats.
func ParseS3Config(input io.Reader) (S3Config, error) {
	var config S3Config
	if input == nil {
		return config, errReaderCannotBeNil
	}

	var buf = new(bytes.Buffer)
	if _, err := buf.ReadFrom(input); err != nil {
		return config, nil
	}

	if err := yaml.Unmarshal(buf.Bytes(), &config); err == nil {
		return config, nil
	}

	if err := json.Unmarshal(buf.Bytes(), &config); err != nil {
		return config, errors.Wrap(err, errFailedParsingConfig.Error())
	}
	return config, nil
}

// S3Config is used to configure an S3 snapshot repository
// Full list of settings in the Elasticsearch official documentation:
// https://www.elastic.co/guide/en/elasticsearch/plugins/current/repository-s3-repository.html
// https://www.elastic.co/guide/en/elasticsearch/plugins/current/repository-s3-client.html
// nolint
type S3Config struct {
	// Required settings
	Region    string `json:"region,omitempty"`
	Bucket    string `json:"bucket,omitempty"`
	AccessKey string `json:"access_key,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`

	// Optional settings
	BasePath             string `json:"base_path,omitempty"`
	Compress             bool   `json:"compress,omitempty"`
	ServerSideEncryption bool   `json:"server_side_encryption,omitempty"`

	// Advanced Settings
	ChunkSize    string `json:"chunk_size,omitempty"`
	CannedACL    string `json:"canned_acl,omitempty"`
	StorageClass string `json:"storage_class,omitempty"`

	// Client settings
	// See http://docs.aws.amazon.com/general/latest/gr/rande.html#s3_region.
	Endpoint        string        `json:"endpoint,omitempty"`
	Protocol        string        `json:"protocol,omitempty"`
	Timeout         time.Duration `json:"timeout,omitempty"`
	MaxRetries      int           `json:"max_retries,omitempty"`
	ThrottleRetries bool          `json:"throttle_retries,omitempty"`
}

// S3TypeConfig is used by the text formatter to wrwap the S3 config with the
// type field.
type S3TypeConfig struct {
	Type     string   `json:"type"`
	Settings S3Config `json:"settings"`
}

// Validate ensures that S3Config is
func (c S3Config) Validate() error {
	var merr = multierror.NewPrefixed("s3 configuration")
	if err := validateRequiredSettings(c); err != nil {
		merr = merr.Append(err)
	}
	if e := validateOptionalSettings(c); e != nil {
		merr = merr.Append(e)
	}

	return merr.ErrorOrNil()
}

func validateRequiredSettings(c S3Config) error {
	var merr = multierror.NewPrefixed("required setting")
	if c.Region == "" {
		merr = merr.Append(errRegionCannotBeEmpty)
	}
	if c.Bucket == "" {
		merr = merr.Append(errBucketCannotBeEmpty)
	}
	if c.AccessKey == "" {
		merr = merr.Append(errAccessKeyCannotBeEmpty)
	}
	if c.SecretKey == "" {
		merr = merr.Append(errSecretKeyCannotBeEmpty)
	}

	return merr.ErrorOrNil()
}

func validateOptionalSettings(c S3Config) error {
	var merr = multierror.NewPrefixed("optional setting")
	if c.StorageClass != "" && !stringInSlice(c.StorageClass, validStorageClasses) {
		merr = merr.Append(errInvalidStorageClass)
	}
	if c.CannedACL != "" && !stringInSlice(c.CannedACL, validCannedACLs) {
		merr = merr.Append(errInvalidCannedACL)
	}
	if c.Endpoint != "" && !govalidator.IsURL(c.Endpoint) {
		merr = merr.Append(errInvalidEndpoint)
	}
	if c.Protocol != "" && !stringInSlice(c.Protocol, validProtocols) {
		merr = merr.Append(errInvalidProtocol)
	}

	return merr.ErrorOrNil()
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
