// Package job is the Vivarcus Job Processor API for wasm guests (tinygo / wasip1 safe).
//
// Conceptually aligned with Veeva Vault Java SDK com.veeva.vault.sdk.api.job.Job.
package job

import (
	"fmt"
	"strings"

	"github.com/vivarcus/vivarcus-sdk/wire"
)

// Meta declares component-level static metadata (Java @JobInfo).
type Meta struct {
	// Name is the Sdkjob FQN. Empty → toolchain derives from module + type.
	Name              string
	Label             string
	Idempotent        bool
	Visible           bool
	AdminConfigurable bool
	AdminCancellable  bool
}

// Validate checks required Meta fields.
func (m Meta) Validate() error {
	if strings.TrimSpace(m.Label) == "" {
		return fmt.Errorf("meta label required")
	}
	return nil
}

// InitContext is the runtime context for Job.Init (Java JobInitContext).
type InitContext struct {
	Job     string
	JobID   string
	JobName string
	VaultID string
	Params  map[string]any
}

// ParamString returns a string job parameter.
func (c InitContext) ParamString(name string) string {
	return stringParam(c.Params, name)
}

// ParamStrings returns a list of strings (JSON array or comma-separated).
func (c InitContext) ParamStrings(name string) []string {
	return stringListParam(c.Params, name)
}

// ProcessContext is the runtime context for Job.Process (Java JobProcessContext).
type ProcessContext struct {
	Job     string
	JobID   string
	JobName string
	TaskID  string
	VaultID string
	Items   []Item
}

// Item is one JobItem (Java JobItem key/value bag).
type Item struct {
	ID     string
	Values map[string]any
}

// String returns a string value from the item.
func (i Item) String(name string) string {
	return stringParam(i.Values, name)
}

// Input is returned from Init (Java JobInputSupplier of JobItems).
type Input struct {
	Items []Item
}

// ItemResult is the per-item outcome of Process.
type ItemResult struct {
	ID           string
	Result       string // success | failure | noop
	ErrorCode    string
	ErrorMessage string
}

// ProcessResult is returned from Process. Empty Results with a nil error
// means every item in the current task succeeded.
type ProcessResult struct {
	Results []ItemResult
}

// Job is the required guest interface (Java Job, minus complete* which the
// platform job engine applies from Process results).
type Job interface {
	Meta() Meta
	Init(ctx InitContext) (Input, error)
	Process(ctx ProcessContext) (ProcessResult, error)
}

// ToWireMeta converts Meta to the describe payload.
func ToWireMeta(m Meta) wire.JobMeta {
	return wire.JobMeta{
		ComponentName:     strings.TrimSpace(m.Name),
		Label:             m.Label,
		Idempotent:        m.Idempotent,
		Visible:           m.Visible,
		AdminConfigurable: m.AdminConfigurable,
		AdminCancellable:  m.AdminCancellable,
	}
}

// InitContextFromWire converts a host init invocation into the guest API shape.
func InitContextFromWire(w wire.JobInitContext) InitContext {
	return InitContext{
		Job:     w.Job,
		JobID:   w.JobID,
		JobName: w.JobName,
		VaultID: w.VaultID,
		Params:  w.Params,
	}
}

// ProcessContextFromWire converts a host process invocation into the guest API shape.
func ProcessContextFromWire(w wire.JobProcessContext) ProcessContext {
	items := make([]Item, 0, len(w.Items))
	for _, it := range w.Items {
		items = append(items, Item{ID: it.ID, Values: it.Values})
	}
	return ProcessContext{
		Job:     w.Job,
		JobID:   w.JobID,
		JobName: w.JobName,
		TaskID:  w.TaskID,
		VaultID: w.VaultID,
		Items:   items,
	}
}

func stringParam(m map[string]any, name string) string {
	if len(m) == 0 {
		return ""
	}
	v, ok := m[name]
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func stringListParam(m map[string]any, name string) []string {
	if len(m) == 0 {
		return nil
	}
	v, ok := m[name]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case []string:
		return append([]string(nil), t...)
	case []any:
		out := make([]string, 0, len(t))
		for _, item := range t {
			s, _ := item.(string)
			s = strings.TrimSpace(s)
			if s != "" {
				out = append(out, s)
			}
		}
		return out
	case string:
		s := strings.TrimSpace(t)
		if s == "" {
			return nil
		}
		return []string{s}
	default:
		return nil
	}
}
