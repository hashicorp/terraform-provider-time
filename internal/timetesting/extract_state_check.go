package timetesting

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var _ statecheck.StateCheck = &ExtractState{}

type ExtractState struct {
	resourceAddress string
	attributePath   tfjsonpath.Path

	// Value contains the string state value after the check has run.
	Value *string
}

func (e *ExtractState) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	var resource *tfjson.StateResource

	if req.State == nil {
		resp.Error = fmt.Errorf("state is nil")

		return
	}

	if req.State.Values == nil {
		resp.Error = fmt.Errorf("state does not contain any state values")

		return
	}

	if req.State.Values.RootModule == nil {
		resp.Error = fmt.Errorf("state does not contain a root module")

		return
	}

	for _, r := range req.State.Values.RootModule.Resources {
		if e.resourceAddress == r.Address {
			resource = r

			break
		}
	}

	if resource == nil {
		resp.Error = fmt.Errorf("%s - Resource not found in state", e.resourceAddress)

		return
	}

	result, err := tfjsonpath.Traverse(resource.AttributeValues, e.attributePath)

	if err != nil {
		resp.Error = err

		return
	}

	strValue, ok := result.(string)
	if !ok {
		resp.Error = fmt.Errorf("error checking value for attribute at path: %s.%s, expected a string value, receieved %T", e.resourceAddress, e.attributePath.String(), result)

		return
	}

	e.Value = &strValue
}

// NewExtractState returns a state check that will extract a state value into an accessible string pointer `(*ExtractState).Value`.
func NewExtractState(resourceAddress string, attributePath tfjsonpath.Path) *ExtractState {
	return &ExtractState{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
	}
}
