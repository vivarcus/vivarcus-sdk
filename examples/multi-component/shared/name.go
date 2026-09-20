package shared

// NameField is stamped by the demo Record Trigger on create.
const NameField = "name__v"

const nameTrigSuffix = "-trig"

// StampName appends the demo trigger suffix (BEFORE_INSERT sample).
func StampName(name string) string {
	return name + nameTrigSuffix
}
