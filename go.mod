module github.com/elastic/ecctl

go 1.13

require (
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/elastic/cloud-sdk-go v1.2.0
	github.com/go-openapi/runtime v0.19.26
	github.com/go-openapi/strfmt v0.20.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.2
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
)

replace sourcegraph.com/sourcegraph/go-diff v0.5.1 => github.com/sourcegraph/go-diff v0.5.1
