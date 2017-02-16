package local

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_setting "github.com/wunderkraut/radi-api/operation/setting"
	handlers_configwrapper "github.com/wunderkraut/radi-handlers/configwrapper"
)

// A handler for local settings
type LocalHandler_Setting struct {
	LocalHandler_Base
	LocalHandler_ConfigWrapperBase
}

// Identify the handler
func (handler *LocalHandler_Setting) Id() string {
	return "local.setting"
}

// [Handler.]Init tells the LocalHandler_Orchestrate to prepare it's operations
func (handler *LocalHandler_Setting) Init() api_operation.Result {
	result := api_operation.New_StandardResult()

	ops := api_operation.Operations{}

	// Make a wrapper for the Settings Config interpretation, based on itnerpreting YML settings
	wrapper := handlers_configwrapper.SettingsConfigWrapper(handlers_configwrapper.New_BaseSettingConfigWrapperYmlOperation(handler.ConfigWrapper()))

	// Now we can add config operations that use that Base class
	ops.Add(api_operation.Operation(&handlers_configwrapper.SettingConfigWrapperGetOperation{Wrapper: wrapper}))
	ops.Add(api_operation.Operation(&handlers_configwrapper.SettingConfigWrapperSetOperation{Wrapper: wrapper}))
	ops.Add(api_operation.Operation(&handlers_configwrapper.SettingConfigWrapperListOperation{Wrapper: wrapper}))

	handler.operations = &ops

	return api_operation.Result(result)
}

// Make ConfigWrapper
func (handler *LocalHandler_Setting) SettingWrapper() api_setting.SettingWrapper {
	return api_setting.New_SimpleSettingWrapper(handler.operations)
}
