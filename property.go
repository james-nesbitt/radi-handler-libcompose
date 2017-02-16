package libcompose

import (
	// libCompose_options "github.com/docker/libcompose/project/options"

	api_property "github.com/wunderkraut/radi-api/property"
	api_usage "github.com/wunderkraut/radi-api/usage"
)

const (
	/**
	 * Configs used in all operations to build the libCompose project
	 */

	// config for a project name
	OPERATION_PROPERTY_LIBCOMPOSE_PROJECTNAME = "compose.projectname"
	// config for a project yml files
	OPERATION_PROPERTY_LIBCOMPOSE_COMPOSEFILES = "compose.composefiles"

	// Input/Output objects
	OPERATION_PROPERTY_LIBCOMPOSE_OUTPUT = "compose.output"
	OPERATION_PROPERTY_LIBCOMPOSE_ERROR  = "compose.error"

	/**
	 * General Properties for most operations
	 */

	// config for an orchestration context limiter
	OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT = "compose.context"

	/**
	 * Operation specific contexts
	 */

	// Individual possible libcompose properties
	OPERATION_PROPERTY_LIBCOMPOSE_FORCEREMOVE      = "compose.forceremove"
	OPERATION_PROPERTY_LIBCOMPOSE_NOCACHE          = "compose.nocache"
	OPERATION_PROPERTY_LIBCOMPOSE_PULL             = "compose.pull"
	OPERATION_PROPERTY_LIBCOMPOSE_DETACH           = "compose.detach"
	OPERATION_PROPERTY_LIBCOMPOSE_NORECREATE       = "compose.norecreate"
	OPERATION_PROPERTY_LIBCOMPOSE_NOBUILD          = "compose.nobuild"
	OPERATION_PROPERTY_LIBCOMPOSE_FORCERECREATE    = "compose.forcerecreate"
	OPERATION_PROPERTY_LIBCOMPOSE_FORCEREBUILD     = "compose.forcerebuild"
	OPERATION_PROPERTY_LIBCOMPOSE_REMOVEVOLUMES    = "compose.removevolumes"
	OPERATION_PROPERTY_LIBCOMPOSE_REMOVEORPHANS    = "compose.removeorphans"
	OPERATION_PROPERTY_LIBCOMPOSE_REMOVEIMAGETYPES = "compose.removeimagetypes"
	OPERATION_PROPERTY_LIBCOMPOSE_REMOVERUNNING    = "compose.removerunning"
	OPERATION_PROPERTY_LIBCOMPOSE_TIMEOUT          = "compose.timeout"
)

/**
 * Properties which the libCompose handler uses
 */

// Project Name Property for a docker.libCompose project
type LibcomposeProjectnameProperty struct {
	api_property.StringProperty
}

// Id for the Property
func (name *LibcomposeProjectnameProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_PROJECTNAME
}

// Label for the Property
func (name *LibcomposeProjectnameProperty) Label() string {
	return "Project name"
}

// Description for the Property
func (name *LibcomposeProjectnameProperty) Description() string {
	return "Compose project name, which is used in container, volume and network naming."
}

// Is the Property internal only
func (name *LibcomposeProjectnameProperty) Usage() api_usage.Usage {
	return api_property.Usage_Internal()
}

// YAML file list Property for a docker.libCompose project
type LibcomposeComposefilesProperty struct {
	api_property.StringSliceProperty
}

// Id for the Property
func (files *LibcomposeComposefilesProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_COMPOSEFILES
}

// Label for the Property
func (files *LibcomposeComposefilesProperty) Label() string {
	return "docker-compose yml file list"
}

// Description for the Property
func (files *LibcomposeComposefilesProperty) Description() string {
	return "An ordered list of docker-compose yml files, which are passed to libcompose."
}

// Is the Property internal only
func (files *LibcomposeComposefilesProperty) Usage() api_usage.Usage {
	return api_property.Usage_Internal()
}

// A libcompose Property for net context limiting
type LibcomposeContextProperty struct {
	api_property.ContextProperty
}

// Id for the Property
func (contextConf *LibcomposeContextProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT
}

// Label for the Property
func (contextConf *LibcomposeContextProperty) Label() string {
	return "context limiter"
}

// Description for the Property
func (contextConf *LibcomposeContextProperty) Description() string {
	return "A context for controling execution."
}

// Is the Property internal only
func (contextConf *LibcomposeContextProperty) Usage() api_usage.Usage {
	return api_property.Usage_Internal()
}

// Output handler Property for a docker.libCompose project
type LibcomposeOutputProperty struct {
	api_property.WriterProperty
}

// Id for the Property
func (output *LibcomposeOutputProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_OUTPUT
}

// Label for the Property
func (output *LibcomposeOutputProperty) Label() string {
	return "Output writer"
}

// Description for the Property
func (output *LibcomposeOutputProperty) Description() string {
	return "Output io.Writer which will receive compose output from containers."
}

// Is the Property internal only
func (output *LibcomposeOutputProperty) Usage() api_usage.Usage {
	return api_property.Usage_Internal()
}

// Error handler Property for a docker.libCompose project
type LibcomposeErrorProperty struct {
	api_property.WriterProperty
}

// Id for the Property
func (err *LibcomposeErrorProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_ERROR
}

// Label for the Property
func (err *LibcomposeErrorProperty) Label() string {
	return "Error writer"
}

// Description for the Property
func (err *LibcomposeErrorProperty) Description() string {
	return "Error io.Writer which will receive compose output from containers."
}

// Is the Property internal only
func (err *LibcomposeErrorProperty) Usage() api_usage.Usage {
	return api_property.Usage_Internal()
}

/**
 * These Properties are wrappers for the various libCompose options
 * structs in https://github.com/docker/libcompose/blob/master/project/options/types.go
 */

// BUILD : Property for a docker.libCompose project to indicate that a build should ignore cached image layers
type LibcomposeNoCacheProperty struct {
	api_property.BooleanProperty
}

// Id for the Property
func (nocache *LibcomposeNoCacheProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_NOCACHE
}

// Label for the Property
func (nocache *LibcomposeNoCacheProperty) Label() string {
	return "nocache"
}

// Description for the Property
func (nocache *LibcomposeNoCacheProperty) Description() string {
	return "When capturing building, ignore cached docker layers?"
}

// Is the Property internal only
func (nocache *LibcomposeNoCacheProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// Property for a docker.libCompose project to indicate that a process remove ... something
type LibcomposeForceRemoveProperty struct {
	api_property.BooleanProperty
}

// Id for the Property
func (forceremove *LibcomposeForceRemoveProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_FORCEREMOVE
}

// Label for the Property
func (forceremove *LibcomposeForceRemoveProperty) Label() string {
	return "Force remove"
}

// Description for the Property
func (forceremove *LibcomposeForceRemoveProperty) Description() string {
	return "When building, force remove .... something?"
}

// Is the Property internal only
func (forceremove *LibcomposeForceRemoveProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// Property for a docker.libCompose project to indicate that a process hsould stay attached and follow
type LibcomposePullProperty struct {
	api_property.BooleanProperty
}

// Id for the Property
func (pull *LibcomposePullProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_PULL
}

// Label for the Property
func (pull *LibcomposePullProperty) Label() string {
	return "Pull"
}

// Description for the Property
func (pull *LibcomposePullProperty) Description() string {
	return "When building, pull all images before using them?"
}

// Is the Property internal only
func (pull *LibcomposePullProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// Property for a docker.libCompose project to indicate that a process hsould stay attached and follow
type LibcomposeDetachProperty struct {
	api_property.BooleanProperty
}

// Id for the Property
func (detach *LibcomposeDetachProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_DETACH
}

// Label for the Property
func (detach *LibcomposeDetachProperty) Label() string {
	return "Detach"
}

// Description for the Property
func (detach *LibcomposeDetachProperty) Description() string {
	return "When capturing output, detach from the output?"
}

// Is the Property internal only
func (detach *LibcomposeDetachProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// UP : Property for a docker.libCompose project to indicate that a process should not create missing containers
type LibcomposeNoRecreateProperty struct {
	api_property.BooleanProperty
}

// Id for the Property
func (norecreate *LibcomposeNoRecreateProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_NORECREATE
}

// Label for the Property
func (norecreate *LibcomposeNoRecreateProperty) Label() string {
	return "Create"
}

// Description for the Property
func (norecreate *LibcomposeNoRecreateProperty) Description() string {
	return "When starting a container, create it first, if it is missing?"
}

// Is the Property internal only
func (norecreate *LibcomposeNoRecreateProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// UP|RECREATE : Property for a docker.libCompose project to indicate that a process should build containers even if they are found
type LibcomposeForceRecreateProperty struct {
	api_property.BooleanProperty
}

// Id for the Property
func (forcerecreate *LibcomposeForceRecreateProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_FORCERECREATE
}

// Label for the Property
func (forcerecreate *LibcomposeForceRecreateProperty) Label() string {
	return "Force Recreate"
}

// Description for the Property
func (forcerecreate *LibcomposeForceRecreateProperty) Description() string {
	return "Force recreating containers, even if they exist already?"
}

// Is the Property internal only
func (forcerecreate *LibcomposeForceRecreateProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// UP|CREATE : Property for a docker.libCompose project to indicate that a process should not build any containers
type LibcomposeNoBuildProperty struct {
	api_property.BooleanProperty
}

// Id for the Property
func (dontbuild *LibcomposeNoBuildProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_NOBUILD
}

// Label for the Property
func (dontbuild *LibcomposeNoBuildProperty) Label() string {
	return "Don't Build"
}

// Description for the Property
func (dontbuild *LibcomposeNoBuildProperty) Description() string {
	return "Don't build any missing images?"
}

// Is the Property internal only
func (dontbuild *LibcomposeNoBuildProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// UP|CREATE : Property for a docker.libCompose project to indicate that a process should force rebuilding images
type LibcomposeForceRebuildProperty struct {
	api_property.BooleanProperty
}

// Id for the Property
func (forcerebuild *LibcomposeForceRebuildProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_FORCEREBUILD
}

// Label for the Property
func (forcerebuild *LibcomposeForceRebuildProperty) Label() string {
	return "Force rebuild"
}

// Description for the Property
func (forcerebuild *LibcomposeForceRebuildProperty) Description() string {
	return "Force rebuilding any images, even if they exist already?"
}

// Is the Property internal only
func (forcerebuild *LibcomposeForceRebuildProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// DOWN|DELETE : Property for a docker.libCompose project to indicate that a process should remove any volumes
type LibcomposeRemoveVolumesProperty struct {
	api_property.BooleanProperty
}

// Id for the Property
func (removevolumes *LibcomposeRemoveVolumesProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_REMOVEVOLUMES
}

// Label for the Property
func (removevolumes *LibcomposeRemoveVolumesProperty) Label() string {
	return "Remove volumes"
}

// Description for the Property
func (removevolumes *LibcomposeRemoveVolumesProperty) Description() string {
	return "When removing containers, remove any volumes?"
}

// Is the Property internal only
func (removevolumes *LibcomposeRemoveVolumesProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// DOWN : Property for a docker.libCompose project to indicate that a process should remove any orphan containers
type LibcomposeRemoveOrphansProperty struct {
	api_property.BooleanProperty
}

// Id for the Property
func (removeorphans *LibcomposeRemoveOrphansProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_REMOVEORPHANS
}

// Label for the Property
func (removeorphans *LibcomposeRemoveOrphansProperty) Label() string {
	return "Remove orphans"
}

// Description for the Property
func (removeorphans *LibcomposeRemoveOrphansProperty) Description() string {
	return "When removing containers, remove any orphans?"
}

// Is the Property internal only
func (removeorphans *LibcomposeRemoveOrphansProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// DOWN : Property for a docker.libCompose project to indicate that a process should remove images of a certain type
type LibcomposeRemoveImageTypeProperty struct {
	api_property.StringProperty
}

// Id for the Property
func (removeimagetypes *LibcomposeRemoveImageTypeProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_REMOVEIMAGETYPES
}

// Label for the Property
func (removeimagetypes *LibcomposeRemoveImageTypeProperty) Label() string {
	return "Remove image types"
}

// Description for the Property
func (removeimagetypes *LibcomposeRemoveImageTypeProperty) Description() string {
	return "When removing containers, remove either 'none' local' or 'all' images?"
}

// Is the Property internal only
func (removeimagetypes *LibcomposeRemoveImageTypeProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// DELETE : Property for a docker.libCompose project to indicate that a process should delete running containers
type LibcomposeRemoveRunningProperty struct {
	api_property.BooleanProperty
}

// Id for the Property
func (removerunning *LibcomposeRemoveRunningProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_REMOVERUNNING
}

// Label for the Property
func (removerunning *LibcomposeRemoveRunningProperty) Label() string {
	return "Remove running"
}

// Description for the Property
func (removerunning *LibcomposeRemoveRunningProperty) Description() string {
	return "When removing containers, remove running containers?"
}

// Is the Property internal only
func (removerunning *LibcomposeRemoveRunningProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}

// STOP : Property for a docker.libCompose project to indicate that how many seconds a process should run for before timing out
type LibcomposeTimeoutProperty struct {
	api_property.IntProperty
}

// Id for the Property
func (timeout *LibcomposeTimeoutProperty) Id() string {
	return OPERATION_PROPERTY_LIBCOMPOSE_TIMEOUT
}

// Label for the Property
func (timeout *LibcomposeTimeoutProperty) Label() string {
	return "Timeout"
}

// Description for the Property
func (timeout *LibcomposeTimeoutProperty) Description() string {
	return "Timeout in seconds before an operation should force"
}

// Is the Property internal only
func (timeout *LibcomposeTimeoutProperty) Usage() api_usage.Usage {
	return api_property.Usage_Optional()
}
