package libcompose

import (
	"context"
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"

	libCompose_config "github.com/docker/libcompose/config"
	libCompose_project_options "github.com/docker/libcompose/project/options"

	api_operation "github.com/wunderkraut/radi-api/operation"
	api_property "github.com/wunderkraut/radi-api/property"
	api_result "github.com/wunderkraut/radi-api/result"
	api_usage "github.com/wunderkraut/radi-api/usage"

	api_command "github.com/wunderkraut/radi-api/operation/command"
	api_config "github.com/wunderkraut/radi-api/operation/config"
)

// Constructor for BaseCommandConfigWrapperYmlOperation
func New_BaseCommandConfigWrapperYmlOperation(configWrapper api_config.ConfigWrapper) *BaseCommandConfigWrapperYmlOperation {
	return &BaseCommandConfigWrapperYmlOperation{
		wrapper:  configWrapper,
		commands: &CommandYmlCommands{},
	}
}

// Command config wrapper that reads YML commands
type BaseCommandConfigWrapperYmlOperation struct {
	wrapper  api_config.ConfigWrapper
	commands *CommandYmlCommands
}

// Load all commands from config
func (commands *BaseCommandConfigWrapperYmlOperation) load() error {
	commands.commands = &CommandYmlCommands{} // reset the command list

	if sources, err := commands.wrapper.Get(CONFIG_KEY_COMMAND); err == nil {
		for _, scope := range sources.Order() {

			scopedSource, _ := sources.Get(scope)

			scopedConfig := CommandYmlConfig{}
			if err := yaml.Unmarshal(scopedSource, &scopedConfig); err == nil {
				scopedCommands := scopedConfig.CommandYmlCommands()
				commands.commands.Merge(scopedCommands)
				log.WithFields(log.Fields{"scope": scope, "merged": commands.commands.Order(), "new": scopedCommands.Order()}).Debug("Commands:Config->Load()")
			} else {
				log.WithError(err).WithFields(log.Fields{"scope": scope}).Error("Commands: Couldn't marshall yml scope")
			}
		}
		return nil
	} else {
		log.WithError(err).Error("Error loading config for " + CONFIG_KEY_COMMAND)
		return err
	}
} // Save Commands
func (commands *BaseCommandConfigWrapperYmlOperation) save() error {
	return errors.New("Not allowed to save commands yet")
}

// CommandSource interface List implementation
func (commands *BaseCommandConfigWrapperYmlOperation) Get(key string) (*CommandYmlCommand, error) {
	if commands.commands.Empty() {
		commands.load()
	}
	if comm, err := commands.commands.Get(key); err == nil {
		log.WithFields(log.Fields{"key": key, "comm": comm}).WithError(err).Debug("Commands:Config->Get()")
		return comm, nil
	} else {
		return nil, err
	}
}

// CommandSource interface List implementation
func (commands *BaseCommandConfigWrapperYmlOperation) Set(key string, comm *CommandYmlCommand) error {
	return errors.New("Not yet able to set Commands")
}

// CommandSource interface List implementation
func (commands *BaseCommandConfigWrapperYmlOperation) List(parent string) ([]string, error) {
	if commands.commands.Empty() {
		commands.load()
	}

	keys := []string{}
	for _, key := range commands.commands.Order() {
		if parent == "" || (key != parent && strings.HasPrefix(key, parent)) {
			keys = append(keys, key)
		}
	}
	return keys, nil
}

// Complete Config container for data from command config sources
//  - Currently only holds commands
type CommandYmlConfig struct {
	Commands CommandYmlCommands `yaml:"Commands"`
}

// Return just the yml commands
func (commConfig *CommandYmlConfig) CommandYmlCommands() CommandYmlCommands {
	return commConfig.Commands
}

// YML Struct for Commands
type CommandYmlCommands struct {
	comms map[string]*CommandYmlCommand
	order []string
}

// Yaml UnMarshaller
func (comms *CommandYmlCommands) UnmarshalYAML(unmarshal func(interface{}) error) error {
	holder := map[string]*CommandYmlCommand{}
	if error := unmarshal(&holder); error == nil {
		for key, comm := range holder {
			comm.setId(key)
			comms.Set(key, comm)
		}
		return nil
	} else {
		return error
	}
}

// Safe lazy constructor
func (comms *CommandYmlCommands) safe() {
	if comms.comms == nil {
		comms.comms = map[string]*CommandYmlCommand{}
		comms.order = []string{}
	}
}

// Safe lazy constructor
func (comms *CommandYmlCommands) Empty() bool {
	return (&comms.comms == nil) || (len(comms.comms) == 0)
}

// Add a command
func (comms *CommandYmlCommands) Set(key string, comm *CommandYmlCommand) error {
	comms.safe()
	if _, exists := comms.comms[key]; !exists {
		comms.order = append(comms.order, key)
	}
	comms.comms[key] = comm
	return nil
}

// Get a comm
func (comms *CommandYmlCommands) Get(key string) (*CommandYmlCommand, error) {
	comms.safe()
	if com, found := comms.comms[key]; found {
		return com, nil
	} else {
		return com, errors.New("Command not found")
	}
}

// Comm order
func (comms *CommandYmlCommands) Order() []string {
	comms.safe()
	return comms.order
}

// Comm merge
func (comms *CommandYmlCommands) Merge(merge CommandYmlCommands) error {
	comms.safe()
	for _, key := range merge.Order() {
		if _, err := comms.Get(key); err != nil {
			mergeComm, _ := merge.Get(key)
			comms.Set(key, mergeComm)
		}
	}
	return nil
}

// A  struct to hold yml commands, which can be used to create command.Command objects
type CommandYmlCommand struct {
	scope string
	id    string

	label       string
	description string
	help        string

	active     bool
	persistant bool
	internal   bool

	project           *ComposeProject
	projectProperties api_property.Properties
	serviceConfig     libCompose_config.ServiceConfig
}

// Yaml UnMarshaller
func (comm *CommandYmlCommand) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var holder struct {
		Scope       string `yaml:"Scope,omitempty"`
		Id          string `yaml:"Id,omitempty"`
		Label       string `yaml:"Label,omitempty"`
		Description string `yaml:"Description,omitempty"`
		Help        string `yaml:"Man,omitempty"`

		Persistant bool `yaml:"Persistant,omitempty"`
		Internal   bool `yaml:"Internal,omitempty"`
		Disabled   bool `yaml:"Disabled,omitempty"`
	}
	if error := unmarshal(&holder); error == nil {
		comm.id = holder.Id
		comm.description = holder.Description
		comm.scope = holder.Scope
		comm.persistant = holder.Persistant
		comm.internal = holder.Internal
		if holder.Disabled {
			comm.active = false
			comm.internal = true
		}
	}

	var serviceHolder libCompose_config.ServiceConfig
	if error := unmarshal(&serviceHolder); error == nil {
		comm.serviceConfig = serviceHolder
	}

	// @TODO can we get properties from the commands yml?

	return nil
}

// Turn this CommandYmlCommand into a command.Command
func (ymlCommand *CommandYmlCommand) Command(projectProps api_property.Properties) api_command.Command {
	ymlCommand.projectProperties = projectProps
	return api_command.Command(ymlCommand)
}

// Return string Id
func (ymlCommand *CommandYmlCommand) setId(id string) {
	ymlCommand.id = id
}

// Return string Scope
func (ymlCommand *CommandYmlCommand) Scope() string {
	return ymlCommand.scope
}

/**
 * Command interace
 */

func (ymlCommand *CommandYmlCommand) Validate() api_result.Result {
	return api_result.New_StandardResult()
}

func (ymlCommand *CommandYmlCommand) Usage() api_usage.Usage {
	return api_operation.Usage_External()
}

// Return string Id
func (ymlCommand *CommandYmlCommand) Id() string {
	return ymlCommand.id
}

// Return string Label
func (ymlCommand *CommandYmlCommand) Label() string {
	return ymlCommand.label
}

// Return string Description
func (ymlCommand *CommandYmlCommand) Description() string {
	return ymlCommand.description
}

// Return string man page
func (ymlCommand *CommandYmlCommand) Help() string {
	return ymlCommand.help
}

// Return string Description
func (ymlCommand *CommandYmlCommand) Properties() api_property.Properties {
	props := api_property.New_SimplePropertiesEmpty()

	// @TODO find a way to add more dynamic properties from YAML

	props.Add(api_property.Property(&api_command.CommandFlagsProperty{}))

	return props.Properties()
}

func (ymlCommand *CommandYmlCommand) Exec(props api_property.Properties) api_result.Result {
	res := api_result.New_StandardResult()

	flags := []string{}
	if propFlags, found := props.Get(api_command.OPERATION_PROPERTY_COMMAND_FLAGS); found {
		flags = propFlags.Get().([]string)
	}

	// @TODO GET this from a property
	runContext := context.Background()

	runOptions := libCompose_project_options.Run{
		Detached: false,
	}

	// get the service for the command
	service := ymlCommand.serviceConfig

	// create a libcompose project
	project, _ := MakeComposeProject(ymlCommand.projectProperties)

	// allow our app to alter the service, to do some string replacements etc
	project.AlterService(&service)

	project.AddConfig(ymlCommand.Id(), &service)
	project.Run(runContext, ymlCommand.Id(), flags, runOptions)

	if !ymlCommand.persistant {
		deleteOptions := libCompose_project_options.Delete{
			RemoveVolume: true,
		}
		project.Delete(runContext, deleteOptions, ymlCommand.Id())
	}

	res.MarkSuccess()
	res.MarkFinished()

	return res.Result()
}
