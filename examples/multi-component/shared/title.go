// Package shared holds cross-entry helpers (ADR-22: each .go is an Sdkcode component).
// Not a schedulable entry — Record Action / Trigger import it at compile time.
package shared

// Demo title written by SetTitle / cleared by ClearTitle.
const DemoTitleValue = "from-sdk"

// EntryTitleValue is written by StampOnEnter (lifecycle entry_action).
const EntryTitleValue = "stamped-on-enter"

// TitleField is the demo object field both actions mutate.
const TitleField = "title__c"

// TitlePatch returns a host-update map for title__c.
func TitlePatch(value string) map[string]any {
	return map[string]any{TitleField: value}
}
