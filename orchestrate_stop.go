package libcompose

import (
	"errors"

	"context"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_property "github.com/wunderkraut/radi-api/property"
	api_result "github.com/wunderkraut/radi-api/result"
	api_usage "github.com/wunderkraut/radi-api/usage"

	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"
)

// Base Stop operation
type BaseLibcomposeOrchestrateStopParametrizedOperation struct{}

// Provide static Properties for the operation
func (base *BaseLibcomposeOrchestrateStopParametrizedOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Add(api_property.Property(&LibcomposeTimeoutProperty{}))

	return props.Properties()
}

// LibCompose based stop orchestrate operation
type LibcomposeOrchestrateStopOperation struct {
	api_orchestrate.BaseOrchestrationStopOperation
	BaseLibcomposeNameFilesOperation
	BaseLibcomposeOrchestrateStopParametrizedOperation
}

// Define the libCompose Orchestrate Down operation usage
func (stop *LibcomposeOrchestrateStopOperation) Id() string {
	return "libcompose." + stop.BaseOrchestrationStopOperation.Id()
}

// Define the libCompose Orchestrate Stop operation usage
func (stop *LibcomposeOrchestrateStopOperation) Usage() api_usage.Usage {
	return api_operation.Usage_External()
}

// Validate the libCompose Orchestrate Stop operation
func (stop *LibcomposeOrchestrateStopOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Provide static properties for the operation
func (stop *LibcomposeOrchestrateStopOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Merge(stop.BaseLibcomposeOrchestrateStopParametrizedOperation.Properties())
	props.Merge(stop.BaseLibcomposeNameFilesOperation.Properties())

	return props.Properties()
}

// Execute the libCompose Orchestrate Stop operation
func (stop *LibcomposeOrchestrateStopOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	// pass all props to make a project
	project, _ := MakeComposeProject(props)

	// some props we will use locally

	var netContext context.Context

	// net context
	if netContextProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT); found {
		netContext = netContextProp.Get().(context.Context)
	} else {
		res.AddError(errors.New("Libcompose stop operation is missing the context property"))
		res.MarkFailed()
	}

	// stop options
	timeout := 10 // timeout in seconds

	if stopOptionsProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_TIMEOUT); found {
		timeout = stopOptionsProp.Get().(int)
	}

	if res.Success() {
		if err := project.APIProject.Stop(netContext, timeout); err != nil {
			res.AddError(err)
			res.MarkFailed()
		} else {
			res.MarkSuccess()
		}
	}

	res.MarkFinished()

	return res.Result()
}
