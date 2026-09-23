// Package action is the Vivarcus Record Action API for wasm guests (tinygo / wasip1 safe).
package action

import (
	"fmt"
	"strings"
	"time"

	"github.com/vivarcus/vivarcus-sdk/wire"
)

// Usage describes where a Record Action may appear.
type Usage string

const (
	UsageUserAction           Usage = "UserAction"
	UsageUserBulkAction       Usage = "UserBulkAction"
	UsageLifecycleUserAction  Usage = "LifecycleUserAction"
	UsageLifecycleEntryAction Usage = "LifecycleEntryAction"
	UsageEventAction          Usage = "EventAction"
	UsageWorkflowStep         Usage = "WorkflowStep"
	UsageWorkflowCancel       Usage = "WorkflowCancel"
	UsageUnspecified          Usage = "Unspecified"
)

// RunAs selects the execution identity.
type RunAs string

const (
	RunAsSystemUser   RunAs = "SystemUser"
	RunAsRequestOwner RunAs = "RequestOwner"
)

// Meta declares component-level static metadata for one Record Action.
type Meta struct {
	// Name is the Recordaction FQN. Empty → toolchain derives from module + type.
	Name string
	Label               string
	Object              string
	// ObjectAction is the Objectaction api_name (e.g. set_title__c). Empty → no button.
	ObjectAction        string
	Usages              []Usage
	Icon                string
	UserInputObject     string
	UserInputObjectType string
	RunAs               RunAs
}

// Validate checks required Meta fields. USER_BULK_ACTION cannot mix with other usages.
func (m Meta) Validate() error {
	if strings.TrimSpace(m.Label) == "" {
		return fmt.Errorf("meta label required")
	}
	hasBulk := false
	for _, u := range m.Usages {
		if u == UsageUserBulkAction {
			hasBulk = true
		}
	}
	if hasBulk && len(m.Usages) != 1 {
		return fmt.Errorf("UserBulkAction cannot be combined with other usages")
	}
	return nil
}

// SupportsUsage reports whether usages include want. Unspecified matches all
// except UserBulkAction (Java Usage.UNSPECIFIED).
func SupportsUsage(usages []Usage, want Usage) bool {
	if len(usages) == 0 {
		usages = []Usage{UsageUnspecified}
	}
	for _, u := range usages {
		if u == want {
			return true
		}
		if u == UsageUnspecified && want != UsageUserBulkAction {
			return true
		}
	}
	return false
}

// RecordActionContext is the runtime context for Record Action execution.
type RecordActionContext struct {
	Action             string
	Records            []Record
	UserInputRecord    *Record
	Configuration      map[string]string
	PreExecutionInputs map[string]string
	Usage              Usage
	UsageEntryPoint    string
	UsageLabel         string
	VaultID            string
	CurrentUserID      string
	InitiatingUserID   string
}

// Record is the SDK view of an object record.
type Record struct {
	ID     string
	Object string
	fields map[string]any
	host   FieldWriter
}

// FieldWriter receives staged field mutations.
type FieldWriter interface {
	SetField(recordID, field string, value any)
}

// AttachFieldWriter connects record mutations to the guest reactor.
func AttachFieldWriter(r *Record, w FieldWriter) {
	if r != nil {
		r.host = w
	}
}

// NewRecord builds a Record from typed wire fields.
func NewRecord(id, object string, fields []wire.Field) Record {
	m := make(map[string]any, len(fields))
	for _, f := range fields {
		m[f.Name] = f.Any()
	}
	return Record{ID: id, Object: object, fields: m}
}

func (r Record) Has(field string) bool {
	_, ok := r.fields[field]
	return ok
}

func (r Record) Lookup(field string) (any, bool) {
	v, ok := r.fields[field]
	return v, ok
}

func (r Record) GetString(field string) string {
	v, ok := r.Lookup(field)
	if !ok || v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

func (r Record) GetInt(field string) int64 {
	v, ok := r.Lookup(field)
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case int64:
		return n
	case int:
		return int64(n)
	case float64:
		return int64(n)
	default:
		return 0
	}
}

func (r Record) GetBool(field string) bool {
	v, ok := r.Lookup(field)
	if !ok || v == nil {
		return false
	}
	b, _ := v.(bool)
	return b
}

// GetTime returns a time field. Wire kind=time stores RFC3339 strings.
func (r Record) GetTime(field string) time.Time {
	v, ok := r.Lookup(field)
	if !ok || v == nil {
		return time.Time{}
	}
	switch t := v.(type) {
	case time.Time:
		return t
	case string:
		parsed, err := time.Parse(time.RFC3339, t)
		if err != nil {
			return time.Time{}
		}
		return parsed
	default:
		return time.Time{}
	}
}

// SetValue stages a field mutation for host persistence after Execute.
func (r *Record) SetValue(field string, value any) {
	if r.fields == nil {
		r.fields = map[string]any{}
	}
	r.fields[field] = value
	if r.host != nil {
		r.host.SetField(r.ID, field, value)
	}
}

// FieldsMap returns a copy of field values.
func (r Record) FieldsMap() map[string]any {
	out := make(map[string]any, len(r.fields))
	for k, v := range r.fields {
		out[k] = v
	}
	return out
}

// ExecuteResult is returned from Execute (Phase 1 empty beyond mutations via SetValue).
type ExecuteResult struct{}

// PreExecuteResult is Java onPreExecute / UserActionExecutionDialog.
type PreExecuteResult struct {
	Title          string
	ConfirmMessage string
}

// PostExecuteResult is Java onPostExecute / UserActionExecutionBanner.
type PostExecuteResult struct {
	Message string
}

// RecordAction is the required guest interface (Java RecordAction).
type RecordAction interface {
	Meta() Meta
	IsExecutable(ctx RecordActionContext) bool
	Execute(ctx RecordActionContext) (ExecuteResult, error)
}

// PreExecuteHook is optional (Java default onPreExecute).
type PreExecuteHook interface {
	OnPreExecute(ctx RecordActionContext) (PreExecuteResult, error)
}

// PostExecuteHook is optional (Java default onPostExecute).
type PostExecuteHook interface {
	OnPostExecute(ctx RecordActionContext) (PostExecuteResult, error)
}

// ToWireMeta converts Meta to the wire shape.
func ToWireMeta(m Meta) wire.Meta {
	usages := make([]string, 0, len(m.Usages))
	for _, u := range m.Usages {
		usages = append(usages, string(u))
	}
	runAs := string(m.RunAs)
	if runAs == "" {
		runAs = string(RunAsSystemUser)
	}
	return wire.Meta{
		ComponentName:       strings.TrimSpace(m.Name),
		Label:               m.Label,
		Object:              m.Object,
		ObjectAction:        strings.TrimSpace(m.ObjectAction),
		Usages:              usages,
		Icon:                m.Icon,
		UserInputObject:     m.UserInputObject,
		UserInputObjectType: m.UserInputObjectType,
		RunAs:               runAs,
	}
}

// ContextFromWire converts wire context into the guest API shape.
func ContextFromWire(w wire.Context, host FieldWriter) RecordActionContext {
	recs := make([]Record, 0, len(w.Records))
	for _, wr := range w.Records {
		rec := NewRecord(wr.ID, wr.Object, wr.Fields)
		AttachFieldWriter(&rec, host)
		recs = append(recs, rec)
	}
	var userIn *Record
	if w.UserInputRecord != nil {
		r := NewRecord(w.UserInputRecord.ID, w.UserInputRecord.Object, w.UserInputRecord.Fields)
		AttachFieldWriter(&r, host)
		userIn = &r
	}
	return RecordActionContext{
		Action:             w.Action,
		Records:            recs,
		UserInputRecord:    userIn,
		Configuration:      w.Configuration,
		PreExecutionInputs: w.PreExecutionInputs,
		Usage:              Usage(w.Usage),
		UsageEntryPoint:    w.UsageEntryPoint,
		UsageLabel:         w.UsageLabel,
		VaultID:            w.VaultID,
		CurrentUserID:      w.CurrentUserID,
		InitiatingUserID:   w.InitiatingUserID,
	}
}
