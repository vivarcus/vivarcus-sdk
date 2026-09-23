//go:build integration

package example_test

import (
	"testing"

	"github.com/vivarcus/vivarcus-sdk/examples/itest"
)

func TestHelloActionIntegration(t *testing.T) {
	env := itest.Open(t)
	env.SDKPut(t, "actions/noop_action.go")
	env.SDKGetContains(t, "acme.corp.sdkdemo.NoopAction", "NoopAction")
}
