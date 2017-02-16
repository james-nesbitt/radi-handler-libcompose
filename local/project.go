package local

import (
	"errors"
	"io"

	log "github.com/Sirupsen/logrus"

	jn_init "github.com/james-nesbitt/init-go"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_project "github.com/wunderkraut/radi-api/operation/project"
	handlers_bytesource "github.com/wunderkraut/radi-handlers/bytesource"
)

/**
 * Local handler for project operations
 */

// A handler for local project handler
type LocalHandler_Project struct {
	LocalHandler_Base
}

// [Handler.]Id returns a string ID for the handler
func (handler *LocalHandler_Project) Id() string {
	return "local.project"
}

// [Handler.]Init tells the LocalHandler_Orchestrate to prepare it's operations
func (handler *LocalHandler_Project) Init() api_operation.Result {
	result := api_operation.New_StandardResult()

	ops := api_operation.Operations{}

	// Now we can add project operations that use that Base class
	ops.Add(api_operation.Operation(&LocalProjectInitOperation{fileSettings: handler.LocalHandler_Base.settings.BytesourceFileSettings}))
	ops.Add(api_operation.Operation(&LocalProjectCreateOperation{fileSettings: handler.LocalHandler_Base.settings.BytesourceFileSettings}))
	ops.Add(api_operation.Operation(&LocalProjectGenerateOperation{fileSettings: handler.LocalHandler_Base.settings.BytesourceFileSettings}))

	handler.operations = &ops

	return api_operation.Result(result)
}

/**
 * Operation to initialize the current project as a radi project
 */

type LocalProjectInitOperation struct {
	api_project.ProjectInitOperation
	handlers_bytesource.BaseBytesourceFilesettingsOperation

	fileSettings handlers_bytesource.BytesourceFileSettings
}

// Id the operation
func (init *LocalProjectInitOperation) Id() string {
	return "local." + init.ProjectInitOperation.Id()
}

// Description for the LocalProjectCreateOperation
func (init *LocalProjectInitOperation) Description() string {
	return "Initialize the current project path as a radi project"
}

// Validate the operation
func (init *LocalProjectInitOperation) Validate() bool {
	return true
}

// Get properties
func (init *LocalProjectInitOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Add(api_operation.Property(&api_project.ProjectInitDemoModeProperty{}))

	bytesourceFilesettings := init.BaseBytesourceFilesettingsOperation.Properties()
	if fileSettingsProp, exists := bytesourceFilesettings.Get(handlers_bytesource.OPERATION_PROPERTY_BYTESOURCE_FILESETTINGS); exists {
		fileSettingsProp.Set(init.fileSettings)
	}
	props.Merge(bytesourceFilesettings)

	return props
}

// Execute the local project init operation
func (init *LocalProjectInitOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	demoModeProp, _ := props.Get(api_project.OPERATION_PROPERTY_PROJECT_INIT_DEMOMODE)
	settingsProp, _ := props.Get(handlers_bytesource.OPERATION_PROPERTY_BYTESOURCE_FILESETTINGS)

	demoMode := demoModeProp.Get().(bool)

	source := "https://raw.githubusercontent.com/wunderkraut/radi-handlers/master/local/template/minimal-init.yml"
	if demoMode {
		source = "https://raw.githubusercontent.com/wunderkraut/radi-handlers/master/local/template/demo-init.yml"
	}

	settings := settingsProp.Get().(handlers_bytesource.BytesourceFileSettings)

	log.WithFields(log.Fields{"source": source, "root": settings.ProjectRootPath}).Info("Running YML processer")

	tasks := jn_init.InitTasks{}
	tasks.Init(settings.ProjectRootPath)
	if !tasks.Init_Yaml_Run(source) {
		result.MarkFailed()
		result.AddError(errors.New("YML Generator failed"))
	} else {
		tasks.RunTasks()
		result.MarkSuccess()
	}

	result.MarkFinished()

	return api_operation.Result(result)
}

/**
 * Operation to create a local project from a bytesource
 */

type LocalProjectCreateOperation struct {
	api_project.ProjectCreateOperation
	handlers_bytesource.BaseBytesourceFilesettingsOperation

	fileSettings handlers_bytesource.BytesourceFileSettings
}

// Id the operation
func (create *LocalProjectCreateOperation) Id() string {
	return "local." + create.ProjectCreateOperation.Id()
}

// Description for the LocalProjectCreateOperation
func (create *LocalProjectCreateOperation) Description() string {
	return "Create a new local project from a yml templating source."
}

// Validate the operation
func (create *LocalProjectCreateOperation) Validate() bool {
	return true
}

// Get properties
func (create *LocalProjectCreateOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	//create.properties.Add(api_operation.Property(&api_project.ProjectCreateTypeProperty{}))
	props.Add(api_operation.Property(&api_project.ProjectCreateSourceProperty{}))

	bytesourceFilesettings := create.BaseBytesourceFilesettingsOperation.Properties()
	if fileSettingsProp, exists := bytesourceFilesettings.Get(handlers_bytesource.OPERATION_PROPERTY_BYTESOURCE_FILESETTINGS); exists {
		fileSettingsProp.Set(create.fileSettings)
	}
	props.Merge(bytesourceFilesettings)

	return props
}

// Execute the local project init operation
func (create *LocalProjectCreateOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	//typeProp, _ := props.Get(api_project.OPERATION_PROPERTY_PROJECT_CREATE_TYPE)
	sourceProp, _ := props.Get(api_project.OPERATION_PROPERTY_PROJECT_CREATE_SOURCE)
	settingsProp, _ := props.Get(handlers_bytesource.OPERATION_PROPERTY_BYTESOURCE_FILESETTINGS)

	source := sourceProp.Get().(string)
	settings := settingsProp.Get().(handlers_bytesource.BytesourceFileSettings)

	log.WithFields(log.Fields{"source": source, "root": settings.ProjectRootPath}).Info("Running YML processer")

	tasks := jn_init.InitTasks{}
	tasks.Init(settings.ProjectRootPath)
	if !tasks.Init_Yaml_Run(source) {
		result.MarkFailed()
		result.AddError(errors.New("YML Generator failed"))
	} else {
		tasks.RunTasks()

		// @TODO Get some err from the tasks run ?

		result.MarkSuccess()
	}

	result.MarkFinished()

	return api_operation.Result(result)
}

/**
 * Operation to create a template from the local project
 */

type LocalProjectGenerateOperation struct {
	api_project.ProjectGenerateOperation
	handlers_bytesource.BaseBytesourceFilesettingsOperation

	fileSettings handlers_bytesource.BytesourceFileSettings
}

// Id the operation
func (generate *LocalProjectGenerateOperation) Id() string {
	return "local." + generate.ProjectGenerateOperation.Id()
}

// Description for the LocalProjectCreateOperation
func (generate *LocalProjectGenerateOperation) Description() string {
	return "Create a yml template from the current project"
}

// Validate the operation
func (generate *LocalProjectGenerateOperation) Validate() bool {
	return true
}

// Get properties
func (generate *LocalProjectGenerateOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	//generate.properties.Add(api_operation.Property(&api_project.ProjectCreateTypeProperty{}))
	props.Merge(generate.BaseBytesourceFilesettingsOperation.Properties())

	bytesourceFilesettings := generate.BaseBytesourceFilesettingsOperation.Properties()
	if fileSettingsProp, exists := bytesourceFilesettings.Get(handlers_bytesource.OPERATION_PROPERTY_BYTESOURCE_FILESETTINGS); exists {
		fileSettingsProp.Set(generate.fileSettings)
	}
	props.Merge(bytesourceFilesettings)

	return props
}

// Execute the local project init operation
func (generate *LocalProjectGenerateOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	//typeProp, _ := props.Get(api_project.OPERATION_PROPERTY_PROJECT_CREATE_TYPE)
	settingsProp, _ := props.Get(handlers_bytesource.OPERATION_PROPERTY_BYTESOURCE_FILESETTINGS)

	settings := settingsProp.Get().(handlers_bytesource.BytesourceFileSettings)

	var method string = "yaml"
	var writer io.Writer

	skip := []string{}

	if method == "test" {
		log.WithFields(log.Fields{"root": settings.ProjectRootPath}).Info("Running TEST YML generator")

		logger := log.StandardLogger().Writer()
		defer logger.Close()
		writer = io.Writer(logger)
	} else {
		projectPath, _ := settings.ConfigPaths.Get("project")
		destination := projectPath.FullPath("init.yml")

		log.WithFields(log.Fields{"root": settings.ProjectRootPath, "path": destination}).Info("Running YML generator")

		/** @TODO REMOVE THIS HARDCODED PATH : make skip allow full paths*/
		skip = append(skip, "radi/init.yml")

		if fileWriter, err := destination.Writer(); err != nil {
			log.WithError(err).Error("Failed to create template file")
			writer = fileWriter
		} else {
			writer = fileWriter
		}

	}

	if settings.ProjectDoesntExist {
		result.MarkFailed()
		result.AddError(errors.New("No project root path has been defined, so no project can be generated."))
	} else if !jn_init.Init_Generate(method, settings.ProjectRootPath, skip, 1024*1024, writer) {
		result.MarkFailed()
		result.AddError(errors.New("YML Generator failed"))
	} else {
		result.MarkSuccess()
	}

	result.MarkFinished()

	return api_operation.Result(result)
}
