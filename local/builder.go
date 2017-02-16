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
	handlers_libcompose "github.com/wunderkraut/radi-handlers/libcompose"
)

/**
 * Provide a builder for the local API
 *
 * The builder works primarily by creating and initializaing
 * the Handlers that are defined in the other files.
 */

// Constructor for LocalBuilder
func New_LocalBuilder(settings LocalAPISettings) *LocalBuilder {
	return &LocalBuilder{
		settings: settings,
	}
}

// Provide a handler for building all local operations
type LocalBuilder struct {
	settings LocalAPISettings

	parent   api_api.API
	handlers api_handler.Handlers

	common_base       *LocalHandler_Base
	common_libcompose *handlers_libcompose.BaseLibcomposeHandler

	Command     api_command.CommandWrapper
	Config      api_config.ConfigWrapper
	Setting     api_setting.SettingWrapper
	Orchestrate api_orchestrate.OrchestrateWrapper
	Security    api_security.SecurityWrapper
}

// IBuilder ID
func (builder *LocalBuilder) Id() string {
	return "local"
}

// Set the parent API, which may need to build Config and Setting Wrappers
func (builder *LocalBuilder) SetAPI(parent api_api.API) {
	builder.parent = parent
}

// Initialize the handler for certain implementations
func (builder *LocalBuilder) Activate(implementations api_builder.Implementations, settingsProvider api_builder.SettingsProvider) error {
	for _, implementation := range implementations.Order() {
		switch implementation {
		case "config":
			builder.build_Config()
		case "setting":
			builder.build_Setting()
		case "project":
			builder.build_Project()
		case "orchestrate":
			builder.build_Orchestrate()
			builder.build_Monitor()
		case "command":
			builder.build_Command()
		case "security":
			builder.build_Security()

		default:
			log.WithFields(log.Fields{"implementation": implementation}).Error("Local builder implementation not available")
		}
	}

	return nil
}

// Get the list of operations for the implementations
func (builder *LocalBuilder) Operations() *api_operation.Operations {
	ops := builder.handlers.Operations()
	return &ops
}

// Get the list of operations for the implementations
func (builder *LocalBuilder) AddHandler(hand api_handler.Handler) {
	builder.handlers.Add(hand)
}

// Create a shareable common base
func (builder *LocalBuilder) base() *LocalHandler_Base {
	if builder.common_base == nil {

		log.Debug("Building new base handler")

		// Create a base handler, which just wraps the settings
		builder.common_base = &LocalHandler_Base{
			settings:   &builder.settings,
			operations: &api_operation.Operations{},
		}
	}
	return builder.common_base
}

// Build a Handler base that produces LibCompose projects
func (builder *LocalBuilder) base_libcompose() *handlers_libcompose.BaseLibcomposeHandler {
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
		builder.common_libcompose = handlers_libcompose.New_BaseLibcomposeHandler(projectName, dockerComposeFiles, runContext, outputWriter, errorWriter, builder.settings.BytesourceFileSettings)
	}

	return builder.common_libcompose
}

// Add local Handlers for Config and Settings
func (builder *LocalBuilder) build_Config() error {
	// Build a config whandler
	local_config := LocalHandler_Config{
		LocalHandler_Base: *builder.base(),
	}
	local_config.Init()
	builder.AddHandler(api_handler.Handler(&local_config))
	// Get a config wrapper for other handlers
	builder.Config = local_config.ConfigWrapper()

	log.WithFields(log.Fields{"ConfigWrapper": builder.Config}).Debug("localBuilder: Built Config Handler")

	return nil
}

// Add local Handlers for Setting
func (builder *LocalBuilder) build_Setting() error {

	// Build a settings handler which uses the configwrapper and the base
	local_setting := LocalHandler_Setting{
		LocalHandler_Base: *builder.base(),
	}
	local_setting.SetConfigWrapper(builder.Config)
	local_setting.Init()
	builder.AddHandler(api_handler.Handler(&local_setting))
	// Get a settings wrapper for other handlers
	builder.Setting = local_setting.SettingWrapper()

	log.WithFields(log.Fields{"SettingWrapper": builder.Setting}).Debug("localBuilder: Built Setting Handler")

	return nil
}

// Make a Local based API object for "no existing project found" to allow for project operations
func (builder *LocalBuilder) build_Project() error {
	// Build a config wrapper using the base for settings
	local_project := LocalHandler_Project{
		LocalHandler_Base: *builder.base(),
	}
	local_project.Init()
	builder.AddHandler(api_handler.Handler(&local_project))

	log.Debug("localBuilder: Built Project Handler")

	return nil
}

// Add local Handlers for Orchestrate operations
func (builder *LocalBuilder) build_Orchestrate() error {
	// Build an orchestration handler
	local_orchestration := LocalHandler_Orchestrate{
		LocalHandler_Base:     *builder.base(),
		BaseLibcomposeHandler: *builder.base_libcompose(),
	}
	local_orchestration.SetSettingWrapper(builder.Setting)
	local_orchestration.Init()
	builder.AddHandler(api_handler.Handler(&local_orchestration))
	// Get an orchestrate wrapper for other handlers
	builder.Orchestrate = local_orchestration.OrchestrateWrapper()

	log.WithFields(log.Fields{"OrchestrateWrapper": builder.Orchestrate}).Debug("localBuilder: Built Orchestrate handler")

	return nil
}

// Add local Handlers for Orchestrate operations
func (builder *LocalBuilder) build_Monitor() error {
	// Build an orchestration handler
	local_monitor := LocalHandler_Monitor{
		LocalHandler_Base:     *builder.base(),
		BaseLibcomposeHandler: *builder.base_libcompose(),
	}
	local_monitor.SetSettingWrapper(builder.Setting)
	local_monitor.Init()
	builder.AddHandler(api_handler.Handler(&local_monitor))

	log.Debug("localBuilder: Built Monitor handler")

	return nil
}

// Add local Handlers for Command operations
func (builder *LocalBuilder) build_Command() error {
	// Build a command Handler
	local_command := LocalHandler_Command{
		LocalHandler_Base:     *builder.base(),
		BaseLibcomposeHandler: *builder.base_libcompose(),
	}
	local_command.SetConfigWrapper(builder.Config)
	local_command.Init()
	builder.AddHandler(api_handler.Handler(&local_command))
	// Get an orchestrate wrapper for other handlers
	builder.Command = local_command.CommandWrapper()

	log.WithFields(log.Fields{"CommandWrapper": builder.Command}).Debug("localBuilder: Built Command Handler")

	return nil
}

// Add local Handlers for Security operations
func (builder *LocalBuilder) build_Security() error {
	// Build a command Handler
	local_security := LocalHandler_Security{
		LocalHandler_Base: *builder.base(),
	}
	local_security.SetConfigWrapper(builder.Config)
	local_security.Init()
	builder.AddHandler(api_handler.Handler(&local_security))
	// Get an orchestrate wrapper for other handlers
	builder.Security = local_security.SecurityWrapper()

	log.WithFields(log.Fields{"CSecurityWrapper": builder.Command}).Debug("localBuilder: Built Security Handler")

	return nil
}
