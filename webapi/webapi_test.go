package webapi

import "testing"

func TestMetaValidate(t *testing.T) {
	t.Parallel()
	if err := (Meta{}).Validate(); err == nil {
		t.Fatal("expected required fields")
	}
	ok := Meta{
		Label:          "Hello",
		EndpointName:   "hello_world",
		MinimumVersion: "v22.3",
		APIGroup:       "integration__c",
	}
	if err := ok.Validate(); err != nil {
		t.Fatal(err)
	}
	if ok.NormalizedRunAs() != "RequestOwner" {
		t.Fatalf("default runas=%q", ok.NormalizedRunAs())
	}
	bad := ok
	bad.EndpointName = "HelloWorld"
	if err := bad.Validate(); err == nil {
		t.Fatal("expected snake_case endpoint")
	}
}
