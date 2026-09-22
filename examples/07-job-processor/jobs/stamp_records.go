package jobs

import (
	"github.com/vivarcus/vivarcus-sdk/job"
	"github.com/vivarcus/vivarcus-sdk/platform"
)

// StampRecords is a Job Processor that writes stamped__c=true on each item.
type StampRecords struct{}

func (StampRecords) Meta() job.Meta {
	return job.Meta{
		Label:             "Stamp Records",
		Idempotent:        true,
		Visible:           true,
		AdminConfigurable: true,
	}
}

func (StampRecords) Init(ctx job.InitContext) (job.Input, error) {
	ids := ctx.ParamStrings("record_ids")
	items := make([]job.Item, 0, len(ids))
	for _, id := range ids {
		items = append(items, job.Item{ID: id, Values: map[string]any{"id": id}})
	}
	return job.Input{Items: items}, nil
}

func (StampRecords) Process(ctx job.ProcessContext) (job.ProcessResult, error) {
	for _, item := range ctx.Items {
		rec, err := platform.Get("sdk_demo__c", item.ID)
		if err != nil {
			return job.ProcessResult{}, err
		}
		if err := platform.Update(rec.Object, rec.ID, map[string]any{
			"stamped__c": true,
		}); err != nil {
			return job.ProcessResult{}, err
		}
	}
	return job.ProcessResult{}, nil
}

