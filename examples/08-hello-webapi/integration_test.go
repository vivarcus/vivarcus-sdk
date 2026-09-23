//go:build integration

package example_test

import (
	"net/http"
	"testing"

	"github.com/vivarcus/vivarcus-sdk/examples/itest"
)

func TestHelloWebAPIIntegration(t *testing.T) {
	env := itest.Open(t)
	env.ApplyMDL(t, "mdl/01-webapigroup.mdl")
	env.SDKPut(t, "webapis/hello_api.go")

	body := env.API(t, http.MethodPost, "/api/v22.3/custom/hello_world", map[string]any{"name": "Ada"})
	data, _ := body["data"].(map[string]any)
	if data == nil {
		if msg, _ := body["message"].(string); msg != "hello Ada" {
			t.Fatalf("body=%v", body)
		}
		return
	}
	if data["message"] != "hello Ada" {
		t.Fatalf("data=%v", data)
	}
}
