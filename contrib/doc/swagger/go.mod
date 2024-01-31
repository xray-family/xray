module github.com/xray-family/xray/contrib/doc/swagger

go 1.21

replace (
	github.com/xray-family/xray => ../../../
	github.com/xray-family/xray/contrib/adapter/http => ../../../contrib/adapter/http
)

require (
	github.com/stretchr/testify v1.8.4
	github.com/swaggo/files v1.0.0
	github.com/swaggo/swag v1.16.2
	github.com/xray-family/xray v0.0.1
	github.com/xray-family/xray/contrib/adapter/http v0.0.0-20230214121206-a6001cd774ce
	golang.org/x/net v0.20.0
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-openapi/jsonpointer v0.20.2 // indirect
	github.com/go-openapi/jsonreference v0.20.4 // indirect
	github.com/go-openapi/spec v0.20.14 // indirect
	github.com/go-openapi/swag v0.22.9 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/tools v0.17.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
