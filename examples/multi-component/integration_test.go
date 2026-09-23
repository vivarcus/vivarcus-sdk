//go:build integration

package example_test

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vivarcus/vivarcus-sdk/examples/itest"
)

func TestMultiComponentIntegration(t *testing.T) {
	env := itest.Open(t)
	const (
		object = "multi_demo__c"
		field  = "title__c"
		fqn    = "acme.corp.sdkdemo.StampDemoOnEnter"
	)
	dir := t.TempDir()
	env.RenderMDL(t, "mdl/02-lifecycle.mdl", filepath.Join(dir, "02-lifecycle.mdl"), []string{"ACTION_FQN=" + fqn})
	env.ApplyMDL(t, filepath.Join(dir, "02-lifecycle.mdl"))

	objectMDL := filepath.Join(dir, "01-object.mdl")
	env.RenderMDL(t, "mdl/01-object.mdl", objectMDL, nil)
	env.ApplyMDL(t, objectMDL)
	env.ApplyMDL(t, "mdl/03-bind-object-lifecycle.mdl")
	vpk := filepath.Join(dir, "action.vpk")
	pkg := exec.Command("bash", filepath.Join("..", "_shared", "scripts", "package-vpk.sh"), ".", vpk,
		"--component", "10:Object:"+object+":"+objectMDL)
	pkg.Dir = env.Dir
	if out, err := pkg.CombinedOutput(); err != nil {
		t.Fatalf("package-vpk: %v\n%s", err, out)
	}

	imported := env.Run(t, "package", "import", "-f", vpk, "--json")
	pkgID, _ := imported["package_id"].(string)
	if pkgID == "" {
		pkgID, _ = imported["id"].(string)
	}
	if pkgID == "" {
		t.Fatalf("package import: %v", imported)
	}
	env.Run(t, "package", "validate", pkgID, "--json")
	env.Run(t, "package", "deploy", pkgID, "--confirm", "--json")

	rec := env.CreateRecord(t, object, map[string]any{"name__v": "multi-component", field: "initial"})
	name, _ := env.Field(t, object, rec, "name__v").(string)
	if !strings.HasSuffix(name, "-trig") {
		t.Fatalf("name__v=%q", name)
	}
	env.ExecuteAction(t, object, rec, "set_title_shared__c")
	if got := env.Field(t, object, rec, field); got != "from-sdk" {
		t.Fatalf("set_title=%v", got)
	}
	env.ExecuteAction(t, object, rec, "clear_title__c")
	if got := env.Field(t, object, rec, field); got != nil && got != "" {
		t.Fatalf("clear_title=%v", got)
	}

	entry := env.CreateRecord(t, object, map[string]any{"name__v": "multi-entry", field: "initial"})
	env.Transition(t, object, entry, "submit__c")
	if got := env.Field(t, object, entry, field); got != "stamped-on-enter" {
		t.Fatalf("entry_action=%v", got)
	}
}
