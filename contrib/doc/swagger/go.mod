module github.com/lxzan/xray/contrib/doc/swagger

go 1.18

replace (
	github.com/lxzan/xray => ../../../
	github.com/lxzan/xray/contrib/adapter/http => ../../../contrib/adapter/http
)

require (
	github.com/lxzan/xray v1.0.0
	github.com/lxzan/xray/contrib/adapter/http v0.0.0-20230214121206-a6001cd774ce
	github.com/stretchr/testify v1.8.1
	github.com/swaggo/files v1.0.0
	github.com/swaggo/swag v1.8.10
	golang.org/x/net v0.7.0
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	golang.org/x/tools v0.1.12 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
