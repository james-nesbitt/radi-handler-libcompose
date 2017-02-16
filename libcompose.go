package libcompose

import (
	"context"
	"io"

	log "github.com/Sirupsen/logrus"

	libCompose_docker "github.com/docker/libcompose/docker"
	libCompose_dockerctx "github.com/docker/libcompose/docker/ctx"
	libCompose_project "github.com/docker/libcompose/project"

	api_operation "github.com/wunderkraut/radi-api/operation"
	handlers_bytesource "github.com/wunderkraut/radi-handlers/bytesource"
)

/**
 * A wrapper for the various libCompose objects that
 * are needed to handler orchestration through libcompose
 */

func MakeComposeProject(properties *api_operation.Properties) (*ComposeProject, bool) {

	projectNameProp, _ := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_PROJECTNAME)
	composeProjectName := projectNameProp.Get().(string)

	bytesourceFilesettingsProp, _ := properties.Get(handlers_bytesource.OPERATION_PROPERTY_BYTESOURCE_FILESETTINGS)
	pathSettings := bytesourceFilesettingsProp.Get().(handlers_bytesource.BytesourceFileSettings)

	projectFilesProp, _ := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_COMPOSEFILES)
	composeFiles := projectFilesProp.Get().([]string)

	contextProp, _ := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_CONTEXT)
	netContext := contextProp.Get().(context.Context)

	outputProp, _ := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_OUTPUT)
	outputWriter := outputProp.Get().(io.Writer)

	errProp, _ := properties.Get(OPERATION_PROPERTY_LIBCOMPOSE_ERROR)
	errorWriter := errProp.Get().(io.Writer)

	loggerFactory := NewLibcomposeLoggerFactory(outputWriter, errorWriter)

	composeContext := &libCompose_dockerctx.Context{
		Context: libCompose_project.Context{
			ComposeFiles:  composeFiles,
			ProjectName:   composeProjectName,
			LoggerFactory: loggerFactory,
		},
	}

	project, err := libCompose_docker.NewProject(composeContext, nil)

	if err != nil {
		log.WithError(err).Fatal("Could not make the docker-compose project.")
		return nil, false
	}

	composeProject := ComposeProject{
		netContext:     netContext,
		composeContext: composeContext,
		APIProject:     project,
		pathSettings:   pathSettings,
	}

	return &composeProject, true
}

// A wundertools wrapper for the APIProject class
type ComposeProject struct {
	libCompose_project.APIProject
	netContext     context.Context
	composeContext *libCompose_dockerctx.Context
	pathSettings   handlers_bytesource.BytesourceFileSettings
}

// get a specific service
func (project *ComposeProject) Service(name string) (libCompose_project.Service, error) {
	return project.composeContext.Project.CreateService(name)
}

// List all the service names
func (project *ComposeProject) serviceNames() []string {
	return project.composeContext.Project.ServiceConfigs.Keys()
}

func (project *ComposeProject) ServicePS(names ...string) (libCompose_project.InfoSet, error) {
	return project.composeContext.Project.Ps(context.Background(), names...)
}

func (project *ComposeProject) Context() context.Context {
	return project.netContext
}
