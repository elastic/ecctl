module github.com/elastic/ecctl

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20180720115003-f9ffefc3facf
	github.com/blang/semver v3.5.1+incompatible
	github.com/chzyer/logex v1.1.10 // indirect
	github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/elastic/cloud-sdk-go v1.0.0-bc2
	github.com/elastic/uptd v1.0.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/runtime v0.19.0
	github.com/go-openapi/strfmt v0.19.0
	github.com/hashicorp/go-multierror v1.0.0
	github.com/marclop/elasticsearch-cli v0.0.0-20190212132418-ee09f7ae57f1
	github.com/pkg/errors v0.8.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	golang.org/x/crypto v0.0.0-20190829043050-9756ffdc2472
	golang.org/x/oauth2 v0.0.0-20190211225200-5f6b76b7c9dd // indirect
	golang.org/x/sync v0.0.0-20190423024810-112230192c58 // indirect
)

replace sourcegraph.com/sourcegraph/go-diff v0.5.1 => github.com/sourcegraph/go-diff v0.5.1
