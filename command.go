package libcompose

import (
	"errors"

	api_command "github.com/wunderkraut/radi-api/operation/command"
	api_property "github.com/wunderkraut/radi-api/property"
	api_result "github.com/wunderkraut/radi-api/result"
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
func (list *LibcomposeCommandListOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Get properties
func (list *LibcomposeCommandListOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Merge(list.BaseCommandKeyKeysOperation.Properties())
	props.Merge(list.BaseLibcomposeNameFilesOperation.Properties())

	return props.Properties()
}

// Execute the libCompose Command List operation
func (list *LibcomposeCommandListOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	keyProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_KEY)
	keysProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_KEYS)

	parent := ""
	if key, ok := keyProp.Get().(string); ok && key != "" {
		parent = key
	}

	if keyList, err := list.Wrapper.List(parent); err == nil {
		keysProp.Set(keyList)
		res.MarkSuccess()
	} else {
		res.MarkFailed()
		res.AddError(err)
	}

	res.MarkFinished()

	return res.Result()
}

// LibCompose Command Get operation
type LibcomposeCommandGetOperation struct {
	api_command.BaseCommandGetOperation
	api_command.BaseCommandKeyCommandOperation
	BaseLibcomposeNameFilesOperation

	Wrapper CommandConfigWrapper
}

// Validate the operation
func (get *LibcomposeCommandGetOperation) Validate() api_result.Result {
	return api_result.MakeSuccessfulResult()
}

// Get properties
func (get *LibcomposeCommandGetOperation) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	props.Merge(get.BaseCommandKeyCommandOperation.Properties())
	props.Merge(get.BaseLibcomposeNameFilesOperation.Properties())

	return props.Properties()
}

// Execute the libCompose Command Get operation
func (get *LibcomposeCommandGetOperation) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	keyProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_KEY)
	commandProp, _ := props.Get(api_command.OPERATION_PROPERTY_COMMAND_COMMAND)

	if key, ok := keyProp.Get().(string); ok && key != "" {

		if comYml, err := get.Wrapper.Get(key); err == nil {
			// pass all props to make a project
			comProps := get.BaseLibcomposeNameFilesOperation.Properties()
			com := comYml.Command(comProps)
			commandProp.Set(com)
			res.MarkSuccess()
		} else {
			res.AddError(err)
			res.MarkFailed()
		}

	} else {
		res.AddError(errors.New("No command name provided."))
		res.MarkFailed()
	}

	res.MarkFinished()

	return res.Result()
}
