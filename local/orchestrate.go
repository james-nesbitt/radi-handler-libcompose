package local

import (
	api_operation "github.com/wunderkraut/radi-api/operation"

	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"

	handler_libcompose "github.com/wunderkraut/radi-handler-libcompose"
	handler_local "github.com/wunderkraut/radi-handlers/local"
)

// A handler for local orchestration using libcompose
type LocalHandler_Orchestrate struct {
	handler_local.LocalHandler_Base
	handler_local.LocalHandler_SettingWrapperBase
	handler_libcompose.BaseLibcomposeHandler
}

// [Handler.]Id returns a string ID for the handler
func (handler *LocalHandler_Orchestrate) Id() string {
	return "libcompose_local.orchestrate"
}

// [Handler.]Init tells the LocalHandler_Orchestrate to prepare it's operations
func (handler *LocalHandler_Orchestrate) Operations() api_operation.Operations {
	ops := api_operation.New_SimpleOperations()

	// Use discovered/default settings to build a base operation struct, to be share across orchestration operations
	baseLibcompose := *handler.BaseLibcomposeHandler.BaseLibcomposeNameFilesOperation()

	// Now we can add orchestration operations that use that Base class
	ops.Add(api_operation.Operation(&handler_libcompose.LibcomposeMonitorLogsOperation{BaseLibcomposeNameFilesOperation: baseLibcompose}))
	ops.Add(api_operation.Operation(&handler_libcompose.LibcomposeOrchestrateUpOperation{BaseLibcomposeNameFilesOperation: baseLibcompose}))
	ops.Add(api_operation.Operation(&handler_libcompose.LibcomposeOrchestrateDownOperation{BaseLibcomposeNameFilesOperation: baseLibcompose}))
	ops.Add(api_operation.Operation(&handler_libcompose.LibcomposeOrchestrateStartOperation{BaseLibcomposeNameFilesOperation: baseLibcompose}))
	ops.Add(api_operation.Operation(&handler_libcompose.LibcomposeOrchestrateStopOperation{BaseLibcomposeNameFilesOperation: baseLibcompose}))

	return ops.Operations()
}

// Make OrchestrateWrapper
func (handler *LocalHandler_Orchestrate) OrchestrateWrapper() api_orchestrate.OrchestrateWrapper {
	return api_orchestrate.New_SimpleOrchestrateWrapper(handler.Operations())
}
