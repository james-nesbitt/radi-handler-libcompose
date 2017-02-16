package libcompose

import (
	"errors"

	"context"

	libCompose_options "github.com/docker/libcompose/project/options"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_property "github.com/wunderkraut/radi-api/property"
	api_result "github.com/wunderkraut/radi-api/result"
	api_usage "github.com/wunderkraut/radi-api/usage"

	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"
)

const (
	// config for down orchestration compose settings
	OPERATION_PROPERTY_LIBCOMPOSE_SETTINGS_DOWN = "compose.down"
)

/**
 * Operations
 */

// Base Down operation
type BaseLibcomposeOrchestrateDownParametrizedOperation struct{}

// Provide static Properties for the operation
func (base *BaseLibcomposeOrchestrateDownParametrizedOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Add(api_property.Property(&LibcomposeRemoveVolumesProperty{}))
	props.Add(api_property.Property(&LibcomposeRemoveImageTypeProperty{}))
	props.Add(api_property.Property(&LibcomposeRemoveOrphansProperty{}))

	return props.Properties()
}

// LibCompose based down orchestrate operation
type LibcomposeOrchestrateDownOperation struct {
	api_orchestrate.BaseOrchestrationDownOperation
	BaseLibcomposeNameFilesOperation
	BaseLibcomposeOrchestrateDownParametrizedOperation
}

// Define the libCompose Orchestrate Down operation usage
func (down *LibcomposeOrchestrateDownOperation) Id() string {
	return "libcompose."down.BaseOrchestrationDownOperation.Id()
}

// Define the libCompose Orchestrate Down operation usage
func (down *LibcomposeOrchestrateDownOperation) Usage() api_usage.Usage {
	return api_operation.Usage_External()
}

// Validate the libCompose Orchestrate Down operation
func (down *LibcomposeOrchestrateDownOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Provide static properties for the operation
func (down *LibcomposeOrchestrateDownOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Merge(down.BaseLibcomposeOrchestrateDownParametrizedOperation.Properties())
	props.Merge(down.BaseLibcomposeNameFilesOperation.Properties())

	return props.Properties()
}

// Execute the libCompose Orchestrate Down operation
func (down *LibcomposeOrchestrateDownOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	// pass all props to make a project
	project, _ := MakeComposeProject(props)

	// some props we will use locally
	var netContext context.Context
	var downOptions libCompose_options.Down

	// net context
	if netContextProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT); found {
		netContext = netContextProp.Get().(context.Context)
	} else {
		res.MarkFailed()
		res.AddError(errors.New("Libcompose up operation is missing the context property"))
	}

	// up options
	downOptions = libCompose_options.Down{}
	if downOptionsProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_REMOVEVOLUMES); found {
		downOptions.RemoveVolume = downOptionsProp.Get().(bool)
	}
	if downOptionsProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_REMOVEIMAGETYPES); found {
		downOptions.RemoveImages = libCompose_options.ImageType(downOptionsProp.Get().(string))
	}
	if downOptionsProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_REMOVEORPHANS); found {
		downOptions.RemoveOrphans = downOptionsProp.Get().(bool)
	}

	if err := project.APIProject.Down(netContext, downOptions); err != nil {
		res.MarkFailed()
		res.AddError(err)
	} else {
		res.MarkSuccess()
	}

	res.MarkFinished()

	return res.Result()
}
