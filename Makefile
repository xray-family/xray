PWD=$(shell pwd)

tidy:
	cd ${PWD}/contrib/adapter/gws && go mod tidy
	cd ${PWD}/contrib/adapter/http && go mod tidy
	#cd ${PWD}/examples/gws && go mod tidy
	#cd ${PWD}/examples/http && go mod tidy
	#go mod tidy
