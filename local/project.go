package local

import (
	"errors"

	log "github.com/Sirupsen/logrus"

	jn_init "github.com/james-nesbitt/init-go"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_property "github.com/wunderkraut/radi-api/property"
	api_result "github.com/wunderkraut/radi-api/result"

	api_project "github.com/wunderkraut/radi-api/operation/project"

	handler_bytesource "github.com/wunderkraut/radi-handlers/bytesource"
	handler_local "github.com/wunderkraut/radi-handlers/local"
)

const (
	LOCAL_PROJECT_CREATE_SOURCE_DEFAULT = "https://raw.githubusercontent.com/wunderkraut/radi-handler-libcompose/master/local/template/minimal-init.yml"
	LOCAL_PROJECT_CREATE_SOURCE_DEMO    = "https://raw.githubusercontent.com/wunderkraut/radi-handler-libcompose/master/local/template/demo-init.yml"
)

/**
 * Local handler for project operations
 */

// A handler for local project handler
type LocalHandler_Project struct {
	handler_local.LocalHandler_Base
}

// [Handler.]Id returns a string ID for the handler
func (handler *LocalHandler_Project) Id() string {
	return "local.project"
}

// [Handler.]Init tells the LocalHandler_Orchestrate to prepare it's operations
func (handler *LocalHandler_Project) Operations() api_operation.Operations {
	ops := api_operation.New_SimpleOperations()

	byteSourceFileSettings := handler.LocalHandler_Base.LocalAPISettings().BytesourceFileSettings

	// Now we can add project operations that use that Base class
	ops.Add(api_operation.Operation(&LocalProjectInitOperation{fileSettings: byteSourceFileSettings}))

	return ops.Operations()
}

/**
 * Operation to initialize the current project as a radi project
 */

type LocalProjectInitOperation struct {
	api_project.ProjectInitOperation
	handler_bytesource.BaseBytesourceFilesettingsOperation

	fileSettings handler_bytesource.BytesourceFileSettings
}

// Id the operation
func (init *LocalProjectInitOperation) Id() string {
	return "libcompose_local." + init.ProjectInitOperation.Id()
}

// Description for the LocalProjectCreateOperation
func (init *LocalProjectInitOperation) Description() string {
	return "Initialize the current project path as a radi project"
}

// Validate the operation
func (init *LocalProjectInitOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Get properties
func (init *LocalProjectInitOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Add(api_property.Property(&api_project.ProjectInitDemoModeProperty{}))

	bytesourceFilesettings := init.BaseBytesourceFilesettingsOperation.Properties()
	if fileSettingsProp, exists := bytesourceFilesettings.Get(handler_bytesource.OPERATION_PROPERTY_BYTESOURCE_FILESETTINGS); exists {
		fileSettingsProp.Set(init.fileSettings)
	}
	props.Merge(bytesourceFilesettings)

	return props
}

// Execute the local project init operation
func (init *LocalProjectInitOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	demoModeProp, _ := props.Get(api_project.OPERATION_PROPERTY_PROJECT_INIT_DEMOMODE)
	settingsProp, _ := props.Get(handler_bytesource.OPERATION_PROPERTY_BYTESOURCE_FILESETTINGS)

	demoMode := demoModeProp.Get().(bool)

	source := LOCAL_PROJECT_CREATE_SOURCE_DEFAULT
	if demoMode {
		source = LOCAL_PROJECT_CREATE_SOURCE_DEMO
	}

	settings := settingsProp.Get().(handler_bytesource.BytesourceFileSettings)

	log.WithFields(log.Fields{"source": source, "root": settings.ProjectRootPath}).Info("Running YML processer")

	tasks := jn_init.InitTasks{}
	tasks.Init(settings.ProjectRootPath)
	if !tasks.Init_Yaml_Run(source) {
		res.MarkFailed()
		res.AddError(errors.New("YML Generator failed"))
	} else {
		tasks.RunTasks()
		res.MarkSuccess()
	}

	res.MarkFinished()

	return res.Result()
}
