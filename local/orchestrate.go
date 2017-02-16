package local

import (
	api_operation "github.com/wunderkraut/radi-api/operation"
	api_orchestrate "github.com/wunderkraut/radi-api/operation/orchestrate"
	handlers_libcompose "github.com/wunderkraut/radi-handlers/libcompose"
)

// A handler for local orchestration using libcompose
type LocalHandler_Orchestrate struct {
	LocalHandler_Base
	LocalHandler_SettingWrapperBase
	handlers_libcompose.BaseLibcomposeHandler
}

// [Handler.]Id returns a string ID for the handler
func (handler *LocalHandler_Orchestrate) Id() string {
	return "local.orchestrate"
}

// [Handler.]Init tells the LocalHandler_Orchestrate to prepare it's operations
func (handler *LocalHandler_Orchestrate) Init() api_operation.Result {
	result := api_operation.New_StandardResult()

	ops := api_operation.Operations{}

	// Use discovered/default settings to build a base operation struct, to be share across orchestration operations
	baseLibcompose := *handler.BaseLibcomposeHandler.LibComposeBaseOp

	// Now we can add orchestration operations that use that Base class
	ops.Add(api_operation.Operation(&handlers_libcompose.LibcomposeMonitorLogsOperation{BaseLibcomposeNameFilesOperation: baseLibcompose}))
	ops.Add(api_operation.Operation(&handlers_libcompose.LibcomposeOrchestrateUpOperation{BaseLibcomposeNameFilesOperation: baseLibcompose}))
	ops.Add(api_operation.Operation(&handlers_libcompose.LibcomposeOrchestrateDownOperation{BaseLibcomposeNameFilesOperation: baseLibcompose}))
	ops.Add(api_operation.Operation(&handlers_libcompose.LibcomposeOrchestrateStartOperation{BaseLibcomposeNameFilesOperation: baseLibcompose}))
	ops.Add(api_operation.Operation(&handlers_libcompose.LibcomposeOrchestrateStopOperation{BaseLibcomposeNameFilesOperation: baseLibcompose}))

	handler.operations = &ops

	return api_operation.Result(result)
}

// Make OrchestrateWrapper
func (handler *LocalHandler_Orchestrate) OrchestrateWrapper() api_orchestrate.OrchestrateWrapper {
	return api_orchestrate.New_SimpleOrchestrateWrapper(handler.operations)
}
