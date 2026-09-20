package platform_test

import (
	"testing"

	"github.com/vivarcus/vivarcus-sdk/platform"
	"github.com/vivarcus/vivarcus-sdk/wire"
)

func TestQueryUnavailable(t *testing.T) {
	platform.QueryFunc = nil
	_, err := platform.Query("select id from demo__c")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestQueryDecode(t *testing.T) {
	want := wire.QueryResult{
		Records: []wire.Record{{
			ID:     "r1",
			Object: "demo__c",
			Fields: []wire.Field{{Name: "name__v", Kind: "string", S: "x"}},
		}},
		Total: 1,
		Size:  1,
	}
	raw, err := wire.Encode(want)
	if err != nil {
		t.Fatal(err)
	}
	platform.QueryFunc = func(vql string) ([]byte, int32) {
		if vql != "select id from demo__c" {
			t.Fatalf("vql=%q", vql)
		}
		return raw, 0
	}
	got, err := platform.Query("select id from demo__c")
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Records) != 1 || got.Records[0].ID != "r1" {
		t.Fatalf("got=%+v", got)
	}
}
