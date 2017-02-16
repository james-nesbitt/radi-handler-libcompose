package libcompose

import (
	"errors"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_command "github.com/wunderkraut/radi-api/operation/command"
)

/**
 * Implement command containers that mix into
 * libCompose orchestrated containers
 */

const (
	CONFIG_KEY_COMMAND = "commands" // The Config key for settings
)

// A wrapper interface which pulls command information from a config wrapper backend
type CommandConfigWrapper interface {
	List(parent string) ([]string, error)
	Get(key string) (*CommandYmlCommand, error)
}

/**
 * Operations
 */

// LibCompose Command List operation
type LibcomposeCommandListOperation struct {
	api_command.BaseCommandListOperation
	api_command.BaseCommandKeyKeysOperation
	BaseLibcomposeNameFilesOperation

	Wrapper CommandConfigWrapper
}

// Validate the operation
func (list *LibcomposeCommandListOperation) Validate() bool {
	return true
}

// Get properties
func (list *LibcomposeCommandListOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Merge(list.BaseCommandKeyKeysOperation.Properties())
	props.Merge(list.BaseLibcomposeNameFilesOperation.Properties())

	return props
}

// Execute the libCompose Command List operation
func (list *LibcomposeCommandListOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	keyProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_KEY)
	keysProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_KEYS)

	parent := ""
	if key, ok := keyProp.Get().(string); ok && key != "" {
		parent = key
	}

	if keyList, err := list.Wrapper.List(parent); err == nil {
		keysProp.Set(keyList)
		result.MarkSuccess()
	} else {
		result.MarkFailed()
		result.AddError(err)
	}

	result.MarkFinished()

	return api_operation.Result(result)
}

// LibCompose Command Get operation
type LibcomposeCommandGetOperation struct {
	api_command.BaseCommandGetOperation
	api_command.BaseCommandKeyCommandOperation
	BaseLibcomposeNameFilesOperation

	Wrapper CommandConfigWrapper
}

// Validate the operation
func (get *LibcomposeCommandGetOperation) Validate() bool {
	return true
}

// Get properties
func (get *LibcomposeCommandGetOperation) Properties() api_operation.Properties {
	props := api_operation.Properties{}

	props.Merge(get.BaseCommandKeyCommandOperation.Properties())
	props.Merge(get.BaseLibcomposeNameFilesOperation.Properties())

	return props
}

// Execute the libCompose Command Get operation
func (get *LibcomposeCommandGetOperation) Exec(props *api_operation.Properties) api_operation.Result {
	result := api_operation.New_StandardResult()

	keyProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_KEY)
	commandProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_COMMAND)

	if key, ok := keyProp.Get().(string); ok && key != "" {

		if comYml, err := get.Wrapper.Get(key); err == nil {
			// pass all props to make a project
			comProps := get.BaseLibcomposeNameFilesOperation.Properties()
			com := comYml.Command(&comProps)
			commandProp.Set(com)
			result.MarkSuccess()
		} else {
			result.AddError(err)
			result.MarkFailed()
		}

	} else {
		result.AddError(errors.New("No command name provided."))
		result.MarkFailed()
	}

	result.MarkFinished()

	return api_operation.Result(result)
}
