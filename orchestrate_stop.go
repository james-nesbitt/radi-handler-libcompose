package libcompose

import (
	"errors"

	"context"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"
)

// Base Stop operation
type BaseLibcomposeOrchestrateStopParametrizedOperation struct{}

// Provide static Properties for the operation
func (base *BaseLibcomposeOrchestrateStopParametrizedOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Add(api_operation.Property(&LibcomposeTimeoutProperty{}))

	return props
}

// LibCompose based stop orchestrate operation
type LibcomposeOrchestrateStopOperation struct {
	api_orchestrate.BaseOrchestrationStopOperation
	BaseLibcomposeNameFilesOperation
	BaseLibcomposeOrchestrateStopParametrizedOperation

	properties *api_operation.Properties
}

// Validate the libCompose Orchestrate Stop operation
func (stop *LibcomposeOrchestrateStopOperation) Validate() bool {
	return true
}

// Provide static properties for the operation
func (stop *LibcomposeOrchestrateStopOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Merge(stop.BaseLibcomposeOrchestrateStopParametrizedOperation.Properties())
	props.Merge(stop.BaseLibcomposeNameFilesOperation.Properties())

	return props
}

// Execute the libCompose Orchestrate Stop operation
func (stop *LibcomposeOrchestrateStopOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	// pass all props to make a project
	project, _ := MakeComposeProject(props)

	// some props we will use locally

	var netContext context.Context

	// net context
	if netContextProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT); found {
		netContext = netContextProp.Get().(context.Context)
	} else {
		result.AddError(errors.New("Libcompose stop operation is missing the context property"))
		result.MarkFailed()
	}

	// stop options
	timeout := 10 // timeout in seconds

	if stopOptionsProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_TIMEOUT); found {
		timeout = stopOptionsProp.Get().(int)
	}

	if result.Success() {
		if err := project.APIProject.Stop(netContext, timeout); err != nil {
			result.AddError(err)
			result.MarkFailed()
		} else {
			result.MarkSuccess()
		}
	}

	result.MarkFinished()

	return api_operation.Result(result)
}
