//go:build integration

package example_test

import (
	"testing"

	"github.com/vivarcus/vivarcus-sdk/examples/itest"
)

func TestJobProcessorIntegration(t *testing.T) {
	env := itest.Open(t)
	env.ApplyMDL(t, "mdl/01-object.mdl")
	env.SDKPut(t, "jobs/stamp_records.go")
	env.SDKGetContains(t, "acme.corp.sdkdemo.StampRecords", "stamped__c")
}
