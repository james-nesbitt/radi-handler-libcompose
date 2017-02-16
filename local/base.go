package local

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_config "github.com/wunderkraut/radi-api/operation/config"
	api_setting "github.com/wunderkraut/radi-api/operation/setting"
)

/**
 * Local handlers provides operations based entirely on the
 * Local environment, primarily based on config files in
 * a project, based on the current path'
 */

// Constructor for a localHandlerBase
func New_LocalHandler_Base(settings *LocalAPISettings) *LocalHandler_Base {
	return &LocalHandler_Base{
		settings:   settings,
		operations: &api_operation.Operations{},
	}
}

// A handler for base local handlers
type LocalHandler_Base struct {
	settings   *LocalAPISettings
	operations *api_operation.Operations
}

// Return the stored operatons
func (base *LocalHandler_Base) Operations() *api_operation.Operations {
	return base.operations
}

// A handler for base local handlers that use a config source (like a yml file)
type LocalHandler_ConfigWrapperBase struct {
	configWrapper api_config.ConfigWrapper
}

// An accessor for the ConfigBase ConfigWrapper
func (base *LocalHandler_ConfigWrapperBase) ConfigWrapper() api_config.ConfigWrapper {
	return base.configWrapper
}

// An accessor to set the ConfigBase ConfigWrapper
func (base *LocalHandler_ConfigWrapperBase) SetConfigWrapper(configWrapper api_config.ConfigWrapper) {
	base.configWrapper = configWrapper
}

// A handler for local settings
type LocalHandler_SettingWrapperBase struct {
	settingWrapper api_setting.SettingWrapper
}

// An accessor for the SettingBase SettingWrapper
func (base *LocalHandler_SettingWrapperBase) SettingWrapper() api_setting.SettingWrapper {
	return base.settingWrapper
}

// An accessor to set the SettingsBase SettingWrapper
func (base *LocalHandler_SettingWrapperBase) SetSettingWrapper(settingWrapper api_setting.SettingWrapper) {
	base.settingWrapper = settingWrapper
}
