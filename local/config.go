package local

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_config "github.com/wunderkraut/radi-api/operation/config"
	handlers_bytesource "github.com/wunderkraut/radi-handlers/bytesource"
)

// A handler for local config
type LocalHandler_Config struct {
	LocalHandler_Base
}

// Identify the handler
func (handler *LocalHandler_Config) Id() string {
	return "local.config"
}

// [Handler.]Init tells the LocalHandler_Orchestrate to prepare it's operations
func (handler *LocalHandler_Config) Init() api_operation.Result {
	result := api_operation.New_StandardResult()

	ops := api_operation.Operations{}

	// build a ConfigConnector for use with the Config operations.
	connector := handlers_bytesource.New_ConfigConnectYmlFiles(handler.settings.ConfigPaths)

	// Build this base operation to be shared across all of our config operations
	baseConnectorOperation := api_config.New_BaseConfigConnectorOperation(connector)

	// Now we can add config operations that use that Base class
	ops.Add(api_operation.Operation(&api_config.ConfigSimpleConnectorReadersOperation{BaseConfigConnectorOperation: *baseConnectorOperation}))
	ops.Add(api_operation.Operation(&api_config.ConfigSimpleConnectorWritersOperation{BaseConfigConnectorOperation: *baseConnectorOperation}))
	ops.Add(api_operation.Operation(&api_config.ConfigSimpleConnectorListOperation{BaseConfigConnectorOperation: *baseConnectorOperation}))

	handler.operations = &ops

	return api_operation.Result(result)
}

// Make ConfigWrapper
func (handler *LocalHandler_Config) ConfigWrapper() api_config.ConfigWrapper {
	return api_config.ConfigWrapper(api_config.New_SimpleConfigWrapper(handler.operations))
}
