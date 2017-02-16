package local

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_config "github.com/wunderkraut/radi-api/operation/config"
	api_security "github.com/wunderkraut/radi-api/operation/security"
	handlers_configwrapper "github.com/wunderkraut/radi-handlers/configwrapper"
)

/**
 * Local implementation of security operations
 *
 * This security implementation bases authentication
 * on local user assessment, and authorization on a
 * configwrapper source, based on a local file.
 *
 * The local file is not considered safe, but rather
 * is considered an initial implementation, which
 * can also be considered as a project based import
 * source to populate external authorization tools.
 */

// A handler for local security
type LocalHandler_Security struct {
	LocalHandler_Base
	LocalHandler_ConfigWrapperBase
}

// Identify the handler
func (handler *LocalHandler_Security) Id() string {
	return "local.security"
}

// [Handler.]Init tells the LocalHandler_Orchestrate to prepare it's operations
func (handler *LocalHandler_Security) Init() api_operation.Result {
	result := api_operation.New_StandardResult()

	ops := api_operation.Operations{}

	// Make a SecurityWrapper Base operation
	securityWrapper := handlers_configwrapper.New_SecurityConfigWrapperYml(handler.ConfigWrapper()).SecurityConfigWrapper()
	base := handlers_configwrapper.New_SecurityWrapperBaseOperation(securityWrapper)

	// Add operations from using the base
	ops.Add(api_operation.Operation(&handlers_configwrapper.SecurityConfigWrapperUserOperation{SecurityWrapperBaseOperation: *base}))
	ops.Add(api_operation.Operation(&handlers_configwrapper.SecurityConfigWrapperAuthorizeOperation{SecurityWrapperBaseOperation: *base}))

	handler.operations = &ops

	return api_operation.Result(result)
}

// Make ConfigWrapper
func (handler *LocalHandler_Security) SecurityWrapper() api_security.SecurityWrapper {
	return api_security.New_SimpleSecurityWrapper(handler.operations).SecurityWrapper()
}

/**
 * A local config based CurrentUser operation
 */

// Local Current user
type LocalCurrentUserOperation struct {
	settings      *LocalAPISettings
	configWrapper api_config.ConfigWrapper
}
