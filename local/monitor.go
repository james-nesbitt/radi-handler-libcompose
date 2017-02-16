package local

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	handlers_libcompose "github.com/wunderkraut/radi-handlers/libcompose"
)

// A handler for local monitoring using libcompose
type LocalHandler_Monitor struct {
	LocalHandler_Base
	LocalHandler_SettingWrapperBase
	handlers_libcompose.BaseLibcomposeHandler
}

// [Handler.]Id returns a string ID for the handler
func (handler *LocalHandler_Monitor) Id() string {
	return "local.monitor"
}

// [Handler.]Init tells the LocalHandler_Orchestrate to prepare it's operations
func (handler *LocalHandler_Monitor) Init() api_operation.Result {
	result := api_operation.New_StandardResult()

	ops := api_operation.Operations{}

	// Use discovered/default settings to build a base operation struct, to be share across orchestration operations
	baseLibcompose := *handler.BaseLibcomposeHandler.LibComposeBaseOp

	// Now we can add orchestration operations that use that Base class
	ops.Add(api_operation.Operation(&handlers_libcompose.LibcomposeOrchestratePsOperation{BaseLibcomposeNameFilesOperation: baseLibcompose}))

	handler.operations = &ops

	return api_operation.Result(result)
}
