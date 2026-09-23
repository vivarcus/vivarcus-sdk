// Package itest drives one sdk/examples module against a logged-in Vault.
// Integration tests live in each example; this package is only the shared client.
package itest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Env is a Vault session for one example directory (the test's working directory).
type Env struct {
	Dir      string
	Endpoint string
	Token    string
	VaultID  string
	Bin      string
	HTTP     *http.Client
}

// Open loads Vault auth and skips when it is not configured.
func Open(t *testing.T) *Env {
	t.Helper()
	if os.Getenv("VIVARCUS_SKIP_EXAMPLES_INTEGRATION") == "1" {
		t.Skip("VIVARCUS_SKIP_EXAMPLES_INTEGRATION=1")
	}
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	bin := strings.TrimSpace(os.Getenv("VIVARCUS_BIN"))
	if bin == "" {
		bin = "vivarcus"
	}
	endpoint := strings.TrimRight(strings.TrimSpace(os.Getenv("VIVARCUS_ENDPOINT")), "/")
	token := strings.TrimSpace(os.Getenv("VIVARCUS_TOKEN"))
	vault := strings.TrimSpace(os.Getenv("VIVARCUS_VAULT"))
	if endpoint == "" || token == "" || vault == "" {
		out, err := exec.Command(bin, "auth", "status", "--json").Output()
		if err != nil {
			t.Skip("not authenticated — set VIVARCUS_TOKEN, VIVARCUS_ENDPOINT, and VIVARCUS_VAULT, or run vivarcus auth login")
		}
		var status map[string]any
		if err := json.Unmarshal(out, &status); err != nil {
			t.Skip("vivarcus auth status: invalid json")
		}
		if token == "" {
			token = stringField(status, "session_token", "token")
		}
		if endpoint == "" {
			endpoint = strings.TrimRight(stringField(status, "endpoint"), "/")
		}
		if vault == "" {
			vault = stringField(status, "default_vault")
		}
	}
	if endpoint == "" {
		endpoint = strings.TrimRight(configGet(bin, "endpoint"), "/")
	}
	if vault == "" {
		vault = configGet(bin, "default_vault")
	}
	if token == "" {
		token = configGet(bin, "token")
	}
	if endpoint == "" || token == "" || vault == "" {
		t.Skip("need endpoint, session token, and vault id")
	}
	return &Env{
		Dir:      dir,
		Endpoint: endpoint,
		Token:    token,
		VaultID:  vault,
		Bin:      bin,
		HTTP:     &http.Client{Timeout: 15 * time.Minute},
	}
}

func (e *Env) Run(t *testing.T, args ...string) map[string]any {
	t.Helper()
	out, err := e.runOutput(t, args...)
	if err != nil {
		t.Fatalf("vivarcus %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	raw := strings.TrimSpace(string(out))
	if raw == "" {
		return map[string]any{}
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(raw), &body); err != nil {
		t.Fatalf("vivarcus json: %v\n%s", err, raw)
	}
	if status, _ := body["responseStatus"].(string); status != "" && status != "SUCCESS" {
		t.Fatalf("vivarcus %s: status=%s body=%v", strings.Join(args, " "), status, body)
	}
	return body
}

func (e *Env) runOutput(t *testing.T, args ...string) ([]byte, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, e.Bin, args...)
	cmd.Dir = e.Dir
	cmd.Env = append(os.Environ(),
		"VIVARCUS_ENDPOINT="+e.Endpoint,
		"VIVARCUS_TOKEN="+e.Token,
		"VIVARCUS_VAULT="+e.VaultID,
	)
	return cmd.CombinedOutput()
}

func (e *Env) ApplyMDL(t *testing.T, path string) {
	t.Helper()
	e.Run(t, "component", "apply-mdl", "--confirm", "-f", path, "--json")
}

func (e *Env) SDKPut(t *testing.T, path string) {
	t.Helper()
	e.Run(t, "sdk", "put", "-f", path, "--json")
}

func (e *Env) SDKGetContains(t *testing.T, fqn, substr string) {
	t.Helper()
	out := filepath.Join(t.TempDir(), "got.go")
	e.Run(t, "sdk", "get", fqn, "-o", out, "--json")
	b, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), substr) {
		t.Fatalf("sdk get %s missing %q", fqn, substr)
	}
}

func (e *Env) RenderMDL(t *testing.T, src, dst string, env []string) {
	t.Helper()
	script := filepath.Join(e.Dir, "..", "_shared", "scripts", "render-mdl.py")
	cmd := exec.Command("python3", script, src, dst)
	cmd.Dir = e.Dir
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("render %s: %v\n%s", src, err, out)
	}
}

func (e *Env) API(t *testing.T, method, path string, payload any) map[string]any {
	t.Helper()
	var body io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, e.Endpoint+path, body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+e.Token)
	req.Header.Set("X-Vault-Id", e.VaultID)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := e.HTTP.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode >= 400 {
		t.Fatalf("%s %s (%d): %s", method, path, resp.StatusCode, string(raw))
	}
	if len(strings.TrimSpace(string(raw))) == 0 {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("decode %s %s: %v\n%s", method, path, err, raw)
	}
	return out
}

func (e *Env) CreateRecord(t *testing.T, object string, fields map[string]any) string {
	t.Helper()
	body := e.API(t, http.MethodPost, "/api/v1/objects/"+object+"/records", map[string]any{"fields": fields})
	id, _ := body["record_id"].(string)
	if id == "" {
		id, _ = body["id"].(string)
	}
	if id == "" {
		t.Fatalf("create %s: missing id: %v", object, body)
	}
	return id
}

func (e *Env) ExecuteAction(t *testing.T, object, recordID, action string) {
	t.Helper()
	e.API(t, http.MethodPost, fmt.Sprintf("/api/v1/objects/%s/records/%s/actions/execute", object, recordID), map[string]any{"action": action})
}

func (e *Env) Field(t *testing.T, object, recordID, field string) any {
	t.Helper()
	out := e.Run(t, "object", "get", object, recordID, "--json")
	data, _ := out["data"].(map[string]any)
	if data == nil {
		t.Fatalf("object get %s %s: %v", object, recordID, out)
	}
	return data[field]
}

func (e *Env) Transition(t *testing.T, object, recordID, action string) map[string]any {
	t.Helper()
	return e.Run(t, "lifecycle", "transition", object, recordID, "--action", action, "--json")
}

func (e *Env) StartWorkflow(t *testing.T, object, recordID, workflow string) map[string]any {
	t.Helper()
	return e.API(t, http.MethodPost, fmt.Sprintf("/api/v1/objects/%s/records/%s/lifecycle/workflows/start", object, recordID), map[string]any{"workflow": workflow})
}

func (e *Env) CancelWorkflow(t *testing.T, instanceID string) {
	t.Helper()
	e.API(t, http.MethodPost, "/api/v1/workflow-instances/"+instanceID+"/cancel", map[string]any{"comment": "example integration"})
}

func stringField(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok && strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func configGet(bin, key string) string {
	out, err := exec.Command(bin, "config", "get", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
