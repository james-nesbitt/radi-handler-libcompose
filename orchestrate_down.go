package libcompose

import (
	"errors"

	"context"

	libCompose_options "github.com/docker/libcompose/project/options"

	api_operation "github.com/wunderkraut/radi-api/operation"
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
func (base *BaseLibcomposeOrchestrateDownParametrizedOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Add(api_operation.Property(&LibcomposeRemoveVolumesProperty{}))
	props.Add(api_operation.Property(&LibcomposeRemoveImageTypeProperty{}))
	props.Add(api_operation.Property(&LibcomposeRemoveOrphansProperty{}))

	return props
}

// LibCompose based down orchestrate operation
type LibcomposeOrchestrateDownOperation struct {
	api_orchestrate.BaseOrchestrationDownOperation
	BaseLibcomposeNameFilesOperation
	BaseLibcomposeOrchestrateDownParametrizedOperation
}

// Validate the libCompose Orchestrate Down operation
func (down *LibcomposeOrchestrateDownOperation) Validate() bool {
	return true
}

// Provide static properties for the operation
func (down *LibcomposeOrchestrateDownOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Merge(down.BaseLibcomposeOrchestrateDownParametrizedOperation.Properties())
	props.Merge(down.BaseLibcomposeNameFilesOperation.Properties())

	return props
}

// Execute the libCompose Orchestrate Down operation
func (down *LibcomposeOrchestrateDownOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	// pass all props to make a project
	project, _ := MakeComposeProject(props)

	// some props we will use locally
	var netContext context.Context
	var downOptions libCompose_options.Down

	// net context
	if netContextProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT); found {
		netContext = netContextProp.Get().(context.Context)
	} else {
		result.MarkFailed()
		result.AddError(errors.New("Libcompose up operation is missing the context property"))
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
		result.MarkFailed()
		result.AddError(err)
	} else {
		result.MarkSuccess()
	}

	result.MarkFinished()

	return api_operation.Result(result)
}
