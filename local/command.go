package local

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_command "github.com/wunderkraut/radi-api/operation/command"
	handlers_libcompose "github.com/wunderkraut/radi-handlers/libcompose"
)

/**
 * Command operations for local projects
 */

// A handler for local command
type LocalHandler_Command struct {
	LocalHandler_Base
	LocalHandler_ConfigWrapperBase
	handlers_libcompose.BaseLibcomposeHandler
}

// Identify the handler
func (handler *LocalHandler_Command) Id() string {
	return "local.command"
}

// [Handler.]Init tells the LocalHandler_Orchestrate to prepare it's operations
func (handler *LocalHandler_Command) Init() api_operation.Result {
	result := api_operation.New_StandardResult()

	ops := api_operation.Operations{}

	// Get shared base operation from the base handler
	baseLibcompose := *handler.BaseLibcomposeHandler.LibComposeBaseOp

	// Make a wrapper for the Command Config interpretation, based on itnerpreting YML settings
	wrapper := handlers_libcompose.CommandConfigWrapper(handlers_libcompose.New_BaseCommandConfigWrapperYmlOperation(handler.ConfigWrapper()))

	ops.Add(api_operation.Operation(&handlers_libcompose.LibcomposeCommandListOperation{BaseLibcomposeNameFilesOperation: baseLibcompose, Wrapper: wrapper}))
	ops.Add(api_operation.Operation(&handlers_libcompose.LibcomposeCommandGetOperation{BaseLibcomposeNameFilesOperation: baseLibcompose, Wrapper: wrapper}))

	handler.operations = &ops

	return api_operation.Result(result)
}

// Make OrchestrateWrapper
func (handler *LocalHandler_Command) CommandWrapper() api_command.CommandWrapper {
	return api_command.New_SimpleCommandWrapper(handler.operations)
}
