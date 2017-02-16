package local

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"os"
	"path"

	api_api "github.com/wunderkraut/radi-api/api"
	api_builder "github.com/wunderkraut/radi-api/builder"
	api_handler "github.com/wunderkraut/radi-api/handler"
	api_result "github.com/wunderkraut/radi-api/result"

	api_command "github.com/wunderkraut/radi-api/operation/command"
	api_config "github.com/wunderkraut/radi-api/operation/config"
	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"
	api_setting "github.com/wunderkraut/radi-api/operation/setting"

	handler_libcompose "github.com/wunderkraut/radi-handler-libcompose"
	handler_local "github.com/wunderkraut/radi-handlers/local"
)

/**
 * Provide a builder for the local API
 *
 * The builder works primarily by creating and initializaing
 * the Handlers that are defined in the other files.
 */

// Provide a handler for building all local operations
type LocalBuilder struct {
	handler_local.LocalBuilder
	settings handler_local.LocalAPISettings

	parent api_api.API

	common_libcompose *handler_libcompose.BaseLibcomposeHandler

	Command     api_command.CommandWrapper
	Orchestrate api_orchestrate.OrchestrateWrapper
}

// Constructor for LocalBuilder
func New_LocalBuilder(settings handler_local.LocalAPISettings) *LocalBuilder {
	return &LocalBuilder{
		LocalBuilder: *handler_local.New_LocalBuilder(settings),
		settings:     settings,
	}
}

// IBuilder ID
func (builder *LocalBuilder) Id() string {
	return "libcompose_local"
}

// Builder Settings
func (builder *LocalBuilder) LocalAPISettings() handler_local.LocalAPISettings {
	return builder.settings
}

// Set the parent API, which may need to build Config and Setting Wrappers
func (builder *LocalBuilder) SetAPI(parent api_api.API) {
	builder.parent = parent
	builder.LocalBuilder.SetAPI(parent)
}

// Initialize the handler for certain implementations
func (builder *LocalBuilder) Activate(implementations api_builder.Implementations, settingsProvider api_builder.SettingsProvider) api_result.Result {
	for _, implementation := range implementations.Order() {
		switch implementation {
		case "project":
			builder.build_Project()
		case "orchestrate":
			builder.build_Orchestrate()
			builder.build_Monitor()
		case "command":
			builder.build_Command()
		default:
			log.WithFields(log.Fields{"implementation": implementation}).Error("Local builder implementation not available")
		}
	}

	return api_result.MakeSuccessfulResult()
}

// Build a Handler base that produces LibCompose projects
func (builder *LocalBuilder) base_libcompose() *handler_libcompose.BaseLibcomposeHandler {
	if builder.common_libcompose == nil {

		log.WithFields(log.Fields{"builder.LocalAPISettings()": builder.LocalAPISettings()}).Debug("Building new Base LibCompose")

		// Set a project name
		projectName := "default"

		settingWrapper := builder.SettingWrapper()

		if settingWrapper == nil {
			log.WithError(errors.New("No setting wrapper avaialble")).Error("Could not set base libCompose project name.")
		} else if settingsProjectName, err := settingWrapper.Get("Project"); err == nil {
			projectName = settingsProjectName
		} else {
			log.WithError(errors.New("Setting value not found in handler config")).Error("Could not set base libCompose project name.")
		}

		// Where to get docker-composer files
		dockerComposeFiles := []string{}
		// add the root composer file
		dockerComposeFiles = append(dockerComposeFiles, path.Join(builder.settings.ProjectRootPath, "docker-compose.yml"))

		// What net context to use
		runContext := builder.settings.Context

		// Output and Error writers
		outputWriter := os.Stdout
		errorWriter := os.Stderr

		// LibComposeHandlerBase
		builder.common_libcompose = handler_libcompose.New_BaseLibcomposeHandler(projectName, dockerComposeFiles, runContext, outputWriter, errorWriter, builder.settings.BytesourceFileSettings)
	}

	return builder.common_libcompose
}

// Add local Handlers for Config and Settings
func (builder *LocalBuilder) ConfigWrapper() api_config.ConfigWrapper {
	// Build a configWrapper if needed
	if builder.Config == nil {
		builder.Config = api_config.New_SimpleConfigWrapper(builder.parent.Operations())
	}
	return builder.Config
}

// Add local Handlers for Setting
func (builder *LocalBuilder) SettingWrapper() api_setting.SettingWrapper {
	// Build a configWrapper if needed
	if builder.Setting == nil {
		builder.Setting = api_setting.New_SimpleSettingWrapper(builder.parent.Operations())
	}
	return builder.Setting
}

// Make a Local based API object for "no existing project found" to allow for project operations
func (builder *LocalBuilder) build_Project() api_result.Result {
	// Build a config wrapper using the base for settings
	local_project := LocalHandler_Project{
		LocalHandler_Base: *builder.Base(),
	}

	res := local_project.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(&local_project))

		log.Debug("LibCompose:localBuilder: Built Project Handler")
	}

	return res
}

// Add local Handlers for Orchestrate operations
func (builder *LocalBuilder) build_Orchestrate() api_result.Result {
	// Build an orchestration handler
	local_orchestration := LocalHandler_Orchestrate{
		LocalHandler_Base:     *builder.Base(),
		BaseLibcomposeHandler: *builder.base_libcompose(),
	}
	local_orchestration.SetSettingWrapper(builder.SettingWrapper())

	res := local_orchestration.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(&local_orchestration))
		// Get an orchestrate wrapper for other handlers
		builder.Orchestrate = local_orchestration.OrchestrateWrapper()

		log.WithFields(log.Fields{"OrchestrateWrapper": builder.Orchestrate}).Debug("LibCompose:localBuilder: Built Orchestrate handler")
	}

	return res
}

// Add local Handlers for Orchestrate operations
func (builder *LocalBuilder) build_Monitor() api_result.Result {
	// Build an orchestration handler
	local_monitor := LocalHandler_Monitor{
		LocalHandler_Base:     *builder.Base(),
		BaseLibcomposeHandler: *builder.base_libcompose(),
	}
	local_monitor.SetSettingWrapper(builder.SettingWrapper())

	res := local_monitor.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(&local_monitor))

		log.Debug("LibCompose:localBuilder: Built Monitor handler")
	}

	return res
}

// Add local Handlers for Command operations
func (builder *LocalBuilder) build_Command() api_result.Result {
	// Build a command Handler
	local_command := LocalHandler_Command{
		LocalHandler_Base:     *builder.Base(),
		BaseLibcomposeHandler: *builder.base_libcompose(),
	}

	local_command.SetConfigWrapper(builder.ConfigWrapper())

	res := local_command.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(&local_command))
		// Get an orchestrate wrapper for other handlers
		builder.Command = local_command.CommandWrapper()

		log.WithFields(log.Fields{"CommandWrapper": builder.Command}).Debug("LibCompose:localBuilder: Built Command Handler")
	}

	return res
}
