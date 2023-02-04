PWD=$(shell pwd)

tidy:
	cd ${PWD}/contrib/adapter/gws && go mod tidy
	cd ${PWD}/contrib/adapter/http && go mod tidy
	go mod tidy

test:
	go test ./...

cover:
	go test --cover ./...