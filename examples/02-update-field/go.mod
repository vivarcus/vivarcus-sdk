module github.com/acme.corp.sdkdemo

go 1.26.2

require (
	github.com/vivarcus/vivarcus-sdk v1.26.3-3.13317
	github.com/vivarcus/vivarcus-sdk/examples/itest v0.0.0
)

replace github.com/vivarcus/vivarcus-sdk => ../..

replace github.com/vivarcus/vivarcus-sdk/examples/itest => ../_shared/itest
