package local

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"os"
	"path"

	api_api "github.com/wunderkraut/radi-api/api"
	api_builder "github.com/wunderkraut/radi-api/builder"
	api_handler "github.com/wunderkraut/radi-api/handler"
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_command "github.com/wunderkraut/radi-api/operation/command"
	api_config "github.com/wunderkraut/radi-api/operation/config"
	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"
	api_security "github.com/wunderkraut/radi-api/operation/security"
	api_setting "github.com/wunderkraut/radi-api/operation/setting"
	api_result "github.com/wunderkraut/radi-api/result"
	handler_libcompose "github.com/wunderkraut/radi-handler-libcompose"
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

	common_libcompose *handler_libcompose.BaseLibcomposeHandler

	Command     api_command.CommandWrapper
	Orchestrate api_orchestrate.OrchestrateWrapper
}

// Constructor for LocalBuilder
func New_LocalBuilder_FromSettings(settings LocalAPISettings) *LocalBuilder {
	return &LocalBuilder{
		LocalBuilder: api_local.New_LocalBuilder(settings),
	}
}

// IBuilder ID
func (builder *LocalBuilder) Id() string {
	return "libcompose_local"
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

// Create a shareable common base
func (builder *LocalBuilder) base() *LocalHandler_Base {
	if builder.common_base == nil {

		log.Debug("Building new base handler")

		// Create a base handler, which just wraps the settings
		builder.common_base = New_LocalHandler_Base(&builder.settings)
	}
	return builder.common_base
}

// Build a Handler base that produces LibCompose projects
func (builder *LocalBuilder) base_libcompose() *handler_libcompose.BaseLibcomposeHandler {
	if builder.common_libcompose == nil {

		log.WithFields(log.Fields{"builder.Setting": builder.Setting}).Debug("Building new Base LibCompose")

		// Set a project name
		projectName := "default"

		if builder.Setting == nil {
			log.WithError(errors.New("No setting wrapper avaialble")).Error("Could not set base libCompose project name.")
		} else if settingsProjectName, err := builder.Setting.Get("Project"); err == nil {
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
func (builder *LocalBuilder) build_Config() api_result.Result {
	// Build a config whandler
	local_config := LocalHandler_Config{
		LocalHandler_Base: *builder.base(),
	}

	res := local_config.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(&local_config))
		// Get a config wrapper for other handlers
		builder.Config = local_config.ConfigWrapper()

		log.WithFields(log.Fields{"ConfigWrapper": builder.Config}).Debug("localBuilder: Built Config Handler")
	}

	return res
}

// Add local Handlers for Setting
func (builder *LocalBuilder) build_Setting() api_result.Result {

	// Build a settings handler which uses the configwrapper and the base
	local_setting := LocalHandler_Setting{
		LocalHandler_Base: *builder.base(),
	}
	local_setting.SetConfigWrapper(builder.Config)

	res := local_setting.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(&local_setting))
		// Get a settings wrapper for other handlers
		builder.Setting = local_setting.SettingWrapper()

		log.WithFields(log.Fields{"SettingWrapper": builder.Setting}).Debug("localBuilder: Built Setting Handler")
	}

	return res
}

// Make a Local based API object for "no existing project found" to allow for project operations
func (builder *LocalBuilder) build_Project() api_result.Result {
	// Build a config wrapper using the base for settings
	local_project := LocalHandler_Project{
		LocalHandler_Base: *builder.base(),
	}

	res := local_project.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(&local_project))

		log.Debug("localBuilder: Built Project Handler")
	}

	return res
}

// Add local Handlers for Orchestrate operations
func (builder *LocalBuilder) build_Orchestrate() api_result.Result {
	// Build an orchestration handler
	local_orchestration := LocalHandler_Orchestrate{
		LocalHandler_Base:     *builder.base(),
		BaseLibcomposeHandler: *builder.base_libcompose(),
	}
	local_orchestration.SetSettingWrapper(builder.Setting)

	res := local_orchestration.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(&local_orchestration))
		// Get an orchestrate wrapper for other handlers
		builder.Orchestrate = local_orchestration.OrchestrateWrapper()

		log.WithFields(log.Fields{"OrchestrateWrapper": builder.Orchestrate}).Debug("localBuilder: Built Orchestrate handler")
	}

	return res
}

// Add local Handlers for Orchestrate operations
func (builder *LocalBuilder) build_Monitor() api_result.Result {
	// Build an orchestration handler
	local_monitor := LocalHandler_Monitor{
		LocalHandler_Base:     *builder.base(),
		BaseLibcomposeHandler: *builder.base_libcompose(),
	}
	local_monitor.SetSettingWrapper(builder.Setting)

	res := local_monitor.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(&local_monitor))

		log.Debug("localBuilder: Built Monitor handler")
	}

	return res
}

// Add local Handlers for Command operations
func (builder *LocalBuilder) build_Command() api_result.Result {
	// Build a command Handler
	local_command := LocalHandler_Command{
		LocalHandler_Base:     *builder.base(),
		BaseLibcomposeHandler: *builder.base_libcompose(),
	}
	local_command.SetConfigWrapper(builder.Config)

	res := local_command.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(&local_command))
		// Get an orchestrate wrapper for other handlers
		builder.Command = local_command.CommandWrapper()

		log.WithFields(log.Fields{"CommandWrapper": builder.Command}).Debug("localBuilder: Built Command Handler")
	}

	return res
}

// Add local Handlers for Security operations
func (builder *LocalBuilder) build_Security() api_result.Result {
	// Build a command Handler
	local_security := LocalHandler_Security{
		LocalHandler_Base: *builder.base(),
	}
	local_security.SetConfigWrapper(builder.Config)

	res := local_security.Validate()
	<-res.Finished()

	if res.Success() {
		builder.AddHandler(api_handler.Handler(&local_security))
		// Get an orchestrate wrapper for other handlers
		builder.Security = local_security.SecurityWrapper()

		log.WithFields(log.Fields{"CSecurityWrapper": builder.Command}).Debug("localBuilder: Built Security Handler")
	}

	return res
}
