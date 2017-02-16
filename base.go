package libcompose

import (
	"io"

	"context"

	api_property "github.com/wunderkraut/radi-api/property"
	handler_bytesource "github.com/wunderkraut/radi-handlers/bytesource"
)

/**
 * Some usefull Base classes used by other libcompose operations
 * and Properties
 */

/**
 * Handlers
 */

// A libcompose handler base that can produce a BaseLibcomposeNameFilesOperation for ops base
type BaseLibcomposeHandler struct {
	LibComposeBaseOp *BaseLibcomposeNameFilesOperation
}

// Constructor for BaseLibcomposeHandler
func New_BaseLibcomposeHandler(projectName string, dockerComposeFiles []string, runContext context.Context, outputWriter io.Writer, errorWriter io.Writer, filesettings handler_bytesource.BytesourceFileSettings) *BaseLibcomposeHandler {
	baseLibcomposeOp := New_BaseLibcomposeNameFilesOperation(projectName, dockerComposeFiles, runContext, outputWriter, errorWriter, filesettings)
	return &BaseLibcomposeHandler{LibComposeBaseOp: baseLibcomposeOp}
}

// Retrieve the BaseLibCompoer Operation
func (base *BaseLibcomposeHandler) BaseLibcomposeNameFilesOperation() *BaseLibcomposeNameFilesOperation {
	return base.LibComposeBaseOp
}

/**
 * Operations
 */

// A base libcompose operation with Properties for project-name, and yml files
type BaseLibcomposeNameFilesOperation struct {
	projectName        string
	dockerComposeFiles []string
	runContext         context.Context
	outputWriter       io.Writer
	errorWriter        io.Writer
	filesettings       handler_bytesource.BytesourceFileSettings
}

// Constructor for BaseLibcomposeNameFilesOperation
func New_BaseLibcomposeNameFilesOperation(projectName string, dockerComposeFiles []string, runContext context.Context, outputWriter io.Writer, errorWriter io.Writer, filesettings handler_bytesource.BytesourceFileSettings) *BaseLibcomposeNameFilesOperation {
	return &BaseLibcomposeNameFilesOperation{
		projectName:        projectName,
		dockerComposeFiles: dockerComposeFiles,
		runContext:         runContext,
		outputWriter:       outputWriter,
		errorWriter:        errorWriter,
		filesettings:       filesettings,
	}
}

// Provide static Properties for the operation - set values from the default
func (base *BaseLibcomposeNameFilesOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	projectName := LibcomposeProjectnameProperty{}
	projectName.Set(base.projectName)
	props.Add(api_property.Property(&projectName))

	filesettings := handler_bytesource.BytesourceFilesettingsProperty{}
	filesettings.Set(base.filesettings)
	props.Add(api_property.Property(&filesettings))

	composeFiles := LibcomposeComposefilesProperty{}
	composeFiles.Set(base.dockerComposeFiles)
	props.Add(api_property.Property(&composeFiles))

	runContext := LibcomposeContextProperty{}
	runContext.Set(base.runContext)
	props.Add(api_property.Property(&runContext))

	output := LibcomposeOutputProperty{}
	output.Set(base.outputWriter)
	props.Add(api_property.Property(&output))
	err := LibcomposeErrorProperty{}
	err.Set(base.errorWriter)
	props.Add(api_property.Property(&err))

	return props.Properties()
}
