package libcompose

import (
	"errors"

	"context"

	libCompose_options "github.com/docker/libcompose/project/options"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"
)

const (
	// config for up orchestration compose settings
	OPERATION_PROPERTY_LIBCOMPOSE_SETTINGS_UP = "compose.up"
)

/**
 * Operation
 */

// Base Up operation
type BaseLibcomposeOrchestrateUpParametrizedOperation struct{}

// Provide static Properties for the operation
func (base *BaseLibcomposeOrchestrateUpParametrizedOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Add(api_operation.Property(&LibcomposeNoRecreateProperty{}))
	props.Add(api_operation.Property(&LibcomposeForceRecreateProperty{}))
	props.Add(api_operation.Property(&LibcomposeNoBuildProperty{}))
	props.Add(api_operation.Property(&LibcomposeForceRebuildProperty{}))

	return props
}

// LibCompose based up orchestrate operation
type LibcomposeOrchestrateUpOperation struct {
	api_orchestrate.BaseOrchestrationUpOperation
	BaseLibcomposeNameFilesOperation
	BaseLibcomposeOrchestrateUpParametrizedOperation
}

// Validate the libCompose Orchestrate Up operation
func (up *LibcomposeOrchestrateUpOperation) Validate() bool {
	return true
}

// Provide static properties for the operation
func (up *LibcomposeOrchestrateUpOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Merge(up.BaseLibcomposeOrchestrateUpParametrizedOperation.Properties())
	props.Merge(up.BaseLibcomposeNameFilesOperation.Properties())

	return props
}

// Execute the libCompose Orchestrate Up operation
func (up *LibcomposeOrchestrateUpOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	// pass all props to make a project
	project, _ := MakeComposeProject(props)

	// some props we will use locally

	var netContext context.Context
	var upOptions libCompose_options.Up
	// net context
	if netContextProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT); found {
		netContext = netContextProp.Get().(context.Context)
	} else {
		result.AddError(errors.New("Libcompose up operation is missing the context property"))
		result.MarkFailed()
	}

	// up options
	upOptions = libCompose_options.Up{}
	if upOptionsProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_NORECREATE); found {
		upOptions.NoRecreate = upOptionsProp.Get().(bool)
	}
	if upOptionsProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_FORCERECREATE); found {
		upOptions.ForceRecreate = upOptionsProp.Get().(bool)
	}
	if upOptionsProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_NOBUILD); found {
		upOptions.NoBuild = upOptionsProp.Get().(bool)
	}
	if upOptionsProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_FORCEREBUILD); found {
		upOptions.ForceBuild = upOptionsProp.Get().(bool)
	}

	if result.Success() {
		if err := project.APIProject.Up(netContext, upOptions); err != nil {
			result.AddError(err)
			result.MarkFailed()
		} else {
			result.MarkSuccess()
		}
	}

	result.MarkFinished()

	return api_operation.Result(result)
}
