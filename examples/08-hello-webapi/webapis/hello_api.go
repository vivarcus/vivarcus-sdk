package webapis

import (
	"github.com/vivarcus/vivarcus-sdk/webapi"
)

// HelloAPI is a Custom Web API that echoes a greeting.
type HelloAPI struct{}

func (HelloAPI) Meta() webapi.Meta {
	return webapi.Meta{
		Label:          "Hello World",
		EndpointName:   "hello_world",
		MinimumVersion: "v22.3",
		APIGroup:       "integration__c",
	}
}

func (HelloAPI) Execute(ctx webapi.Context) (webapi.Response, error) {
	name := "world"
	if ctx.JSON != nil {
		if v, ok := ctx.JSON["name"].(string); ok && v != "" {
			name = v
		}
	}
	return webapi.Response{
		Status: webapi.StatusSuccess,
		Data:   map[string]any{"message": "hello " + name},
	}, nil
}

