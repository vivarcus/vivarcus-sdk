package main

import "github.com/vivarcus/vivarcus-sdk/job"

// {{TYPE_NAME}} is a Job Processor.
type {{TYPE_NAME}} struct{}

func ({{TYPE_NAME}}) Meta() job.Meta {
	return job.Meta{
		Label:             "{{LABEL}}",
		Idempotent:        true,
		Visible:           true,
		AdminConfigurable: true,
	}
}

func ({{TYPE_NAME}}) Init(ctx job.InitContext) (job.Input, error) {
	ids := ctx.ParamStrings("record_ids")
	items := make([]job.Item, 0, len(ids))
	for _, id := range ids {
		items = append(items, job.Item{ID: id, Values: map[string]any{"id": id}})
	}
	return job.Input{Items: items}, nil
}

func ({{TYPE_NAME}}) Process(ctx job.ProcessContext) (job.ProcessResult, error) {
	return job.ProcessResult{}, nil
}

func main() {}
