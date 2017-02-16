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
	// config for up orchestration compose settings
	OPERATION_PROPERTY_LIBCOMPOSE_SETTINGS_UP = "compose.up"
)

/**
 * Operation
 */

// Base Up operation
type BaseLibcomposeOrchestrateUpParametrizedOperation struct{}

// Provide static Properties for the operation
func (base *BaseLibcomposeOrchestrateUpParametrizedOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Add(api_property.Property(&LibcomposeNoRecreateProperty{}))
	props.Add(api_property.Property(&LibcomposeForceRecreateProperty{}))
	props.Add(api_property.Property(&LibcomposeNoBuildProperty{}))
	props.Add(api_property.Property(&LibcomposeForceRebuildProperty{}))

	return props.Properties()
}

// LibCompose based up orchestrate operation
type LibcomposeOrchestrateUpOperation struct {
	api_orchestrate.BaseOrchestrationUpOperation
	BaseLibcomposeNameFilesOperation
	BaseLibcomposeOrchestrateUpParametrizedOperation
}

// Define the libCompose Orchestrate Down operation usage
func (up *LibcomposeOrchestrateUpOperation) Id() string {
	return "libcompose." + up.BaseOrchestrationUpOperation.Id()
}

// Define the libCompose Orchestrate Up operation usage
func (up *LibcomposeOrchestrateUpOperation) Usage() api_usage.Usage {
	return api_operation.Usage_External()
}

// Validate the libCompose Orchestrate Up operation
func (up *LibcomposeOrchestrateUpOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Provide static properties for the operation
func (up *LibcomposeOrchestrateUpOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Merge(up.BaseLibcomposeOrchestrateUpParametrizedOperation.Properties())
	props.Merge(up.BaseLibcomposeNameFilesOperation.Properties())

	return props.Properties()
}

// Execute the libCompose Orchestrate Up operation
func (up *LibcomposeOrchestrateUpOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	// pass all props to make a project
	project, _ := MakeComposeProject(props)

	// some props we will use locally

	var netContext context.Context
	var upOptions libCompose_options.Up
	// net context
	if netContextProp, found := props.Get(OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT); found {
		netContext = netContextProp.Get().(context.Context)
	} else {
		res.AddError(errors.New("Libcompose up operation is missing the context property"))
		res.MarkFailed()
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

	if res.Success() {
		if err := project.APIProject.Up(netContext, upOptions); err != nil {
			res.AddError(err)
			res.MarkFailed()
		} else {
			res.MarkSuccess()
		}
	}

	res.MarkFinished()

	return res.Result()
}
