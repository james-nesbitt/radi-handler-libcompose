package local

import (
	api_operation "github.com/wunderkraut/radi-api/operation"

	handler_libcompose "github.com/wunderkraut/radi-handler-libcompose"
)

// A handler for local monitoring using libcompose
type LocalHandler_Monitor struct {
	LocalHandler_Base
	LocalHandler_SettingWrapperBase
	handler_libcompose.BaseLibcomposeHandler
}

// [Handler.]Id returns a string ID for the handler
func (handler *LocalHandler_Monitor) Id() string {
	return "libcompose_local.monitor"
}

// [Handler.]Init tells the LocalHandler_Orchestrate to prepare it's operations
func (handler *LocalHandler_Monitor) Operations() api_operation.Operations {
	ops := api_operation.New_SimpleOperations()

	// Use discovered/default settings to build a base operation struct, to be share across orchestration operations
	baseLibcompose := *handler.BaseLibcomposeHandler.BaseLibcomposeNameFilesOperation()

	// Now we can add orchestration operations that use that Base class
	ops.Add(api_operation.Operation(&handler_libcompose.LibcomposeOrchestratePsOperation{BaseLibcomposeNameFilesOperation: baseLibcompose}))

	return ops.Operations()
}
