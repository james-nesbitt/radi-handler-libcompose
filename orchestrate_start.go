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

// Base Start operation
type BaseLibcomposeOrchestrateStartParametrizedOperation struct{}

// Provide static Properties for the operation
func (base *BaseLibcomposeOrchestrateStartParametrizedOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	return props.Properties()
}

// LibCompose based start orchestrate operation
type LibcomposeOrchestrateStartOperation struct {
	api_orchestrate.BaseOrchestrationStartOperation
	BaseLibcomposeNameFilesOperation
	BaseLibcomposeOrchestrateStartParametrizedOperation
}

// Define the libCompose Orchestrate Down operation usage
func (start *LibcomposeOrchestrateStartOperation) Id() string {
	return "libcompose."down.BaseOrchestrationStartOperation.Id()
}

// Define the libCompose Orchestrate Start operation usage
func (start *LibcomposeOrchestrateStartOperation) Usage() api_usage.Usage {
	return api_operation.Usage_External()
}

// Validate the libCompose Orchestrate Start operation
func (start *LibcomposeOrchestrateStartOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Provide static properties for the operation
func (start *LibcomposeOrchestrateStartOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Merge(start.BaseLibcomposeOrchestrateStartParametrizedOperation.Properties())
	props.Merge(start.BaseLibcomposeNameFilesOperation.Properties())

	return props.Properties()
}

// Execute the libCompose Orchestrate Start operation
func (start *LibcomposeOrchestrateStartOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	// pass all props to make a project
	project, _ := MakeComposeProject(props)

	// some props we will use locally

	var netContext context.Context

	// net context
	if netContextProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT); found {
		netContext = netContextProp.Get().(context.Context)
	} else {
		res.AddError(errors.New("Libcompose start operation is missing the context property"))
		res.MarkFailed()
	}

	if res.Success() {
		if err := project.APIProject.Start(netContext); err != nil {
			res.MarkFailed()
			res.AddError(err)
		} else {
			res.MarkSuccess()
		}
	}

	res.MarkFinished()

	return res.Result()
}
