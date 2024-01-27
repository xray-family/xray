PWD=$(shell pwd)

tidy:
	cd ${PWD}/contrib/adapter/gws && go mod tidy
	cd ${PWD}/contrib/adapter/http && go mod tidy
	cd ${PWD}/contrib/adapter/fasthttp && go mod tidy
	cd ${PWD}/contrib/codec/jsoniter && go mod tidy
	cd ${PWD}/contrib/codec/wwwform && go mod tidy
	cd ${PWD}/contrib/doc/swagger && go mod tidy
	cd ${PWD}/contrib/log/zerolog && go mod tidy
	cd ${PWD}/examples/http && go mod tidy
	cd ${PWD}/examples/http3 && go mod tidy
	cd ${PWD}/examples/fasthttp && go mod tidy
	cd ${PWD}/examples/gws && go mod tidy
	go mod tidy

test:
	cd ${PWD}/contrib/adapter/gws && go test --count=1 ./...
	cd ${PWD}/contrib/adapter/http && go test --count=1 ./...
	cd ${PWD}/contrib/codec/jsoniter && go test --count=1 ./...
	cd ${PWD}/contrib/codec/wwwform && go test --count=1 ./...
	cd ${PWD}/contrib/log/zerolog && go test --count=1 ./...
	cd ${PWD}/contrib/doc/swagger && go test --count=1 ./...
	go test --count=1 ./...

cover:
	go test -coverprofile=./bin/cover.out --cover ./...
