package shared

import "strings"

// NameField is stamped by the demo Record Trigger on create.
const NameField = "name__v"

const nameTrigSuffix = "-trig"

// StampName appends the demo trigger suffix once (BEFORE_INSERT sample).
// Idempotent so this helper and examples/05-stamp-trigger can both be deployed.
func StampName(name string) string {
	if strings.HasSuffix(name, nameTrigSuffix) {
		return name
	}
	return name + nameTrigSuffix
}
