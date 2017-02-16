package local

import (
	api_handler "github.com/wunderkraut/radi-api/handler"
	api_operation "github.com/wunderkraut/radi-api/operation"

	api_command "github.com/wunderkraut/radi-api/operation/command"
	handler_libcompose "github.com/wunderkraut/radi-handler-libcompose"
)

/**
 * Command operations for local projects
 */

// A handler for local command
type LocalHandler_Command struct {
	LocalHandler_Base
	LocalHandler_ConfigWrapperBase
	handler_libcompose.BaseLibcomposeHandler
}

// Identify the handler
func (commHandler *LocalHandler_Command) Handler() api_handler.Handler {
	return api_handler.Handler(commHandler)
}

// Identify the handler
func (commHandler *LocalHandler_Command) Id() string {
	return "libcompose_local.command"
}

// [Handler.]Init tells the LocalHandler_Orchestrate to prepare it's operations
func (commHandler *LocalHandler_Command) Operations() api_operation.Operations {
	ops := api_operation.New_SimpleOperations()

	// Get shared base operation from the base handler
	baseLibcompose := *commHandler.BaseLibcomposeHandler.BaseLibcomposeNameFilesOperation()

	// Make a wrapper for the Command Config interpretation, based on itnerpreting YML settings
	wrapper := handler_libcompose.CommandConfigWrapper(handler_libcompose.New_BaseCommandConfigWrapperYmlOperation(commHandler.ConfigWrapper()))

	ops.Add(api_operation.Operation(&handler_libcompose.LibcomposeCommandListOperation{BaseLibcomposeNameFilesOperation: baseLibcompose, Wrapper: wrapper}))
	ops.Add(api_operation.Operation(&handler_libcompose.LibcomposeCommandGetOperation{BaseLibcomposeNameFilesOperation: baseLibcompose, Wrapper: wrapper}))

	return ops.Operations()
}

// Make OrchestrateWrapper
func (commHandler *LocalHandler_Command) CommandWrapper() api_command.CommandWrapper {
	return api_command.New_SimpleCommandWrapper(commHandler.Operations())
}
