package libcompose

import (
	"errors"

	"context"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"
)

// Base Start operation
type BaseLibcomposeOrchestrateStartParametrizedOperation struct{}

// Provide static Properties for the operation
func (base *BaseLibcomposeOrchestrateStartParametrizedOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	return props
}

// LibCompose based start orchestrate operation
type LibcomposeOrchestrateStartOperation struct {
	api_orchestrate.BaseOrchestrationStartOperation
	BaseLibcomposeNameFilesOperation
	BaseLibcomposeOrchestrateStartParametrizedOperation
}

// Validate the libCompose Orchestrate Start operation
func (start *LibcomposeOrchestrateStartOperation) Validate() bool {
	return true
}

// Provide static properties for the operation
func (start *LibcomposeOrchestrateStartOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Merge(start.BaseLibcomposeOrchestrateStartParametrizedOperation.Properties())
	props.Merge(start.BaseLibcomposeNameFilesOperation.Properties())

	return props
}

// Execute the libCompose Orchestrate Start operation
func (start *LibcomposeOrchestrateStartOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	// pass all props to make a project
	project, _ := MakeComposeProject(props)

	// some props we will use locally

	var netContext context.Context

	// net context
	if netContextProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT); found {
		netContext = netContextProp.Get().(context.Context)
	} else {
		result.AddError(errors.New("Libcompose start operation is missing the context property"))
		result.MarkFailed()
	}

	if result.Success() {
		if err := project.APIProject.Start(netContext); err != nil {
			result.MarkFailed()
			result.AddError(err)
		} else {
			result.MarkSuccess()
		}
	}

	result.MarkFinished()

	return api_operation.Result(result)
}
