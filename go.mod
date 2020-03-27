module github.com/elastic/ecctl

go 1.12

require (
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/blang/semver v3.5.1+incompatible
	github.com/chzyer/logex v1.1.10 // indirect
	github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/elastic/cloud-sdk-go v1.0.0-beta2.0.20200325215536-c64d182f57e4
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/runtime v0.19.14
	github.com/go-openapi/strfmt v0.19.5
	github.com/hashicorp/go-multierror v1.0.0
	github.com/marclop/elasticsearch-cli v0.0.0-20190212133917-c1d1bf9d46e4
	github.com/pkg/errors v0.9.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/spf13/cobra v0.0.6
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.2
	golang.org/x/crypto v0.0.0-20190829043050-9756ffdc2472
)

replace sourcegraph.com/sourcegraph/go-diff v0.5.1 => github.com/sourcegraph/go-diff v0.5.1
