package libcompose

/**
 * Alter a service, based on the app.
 *
 * This primarily offers a way to use short forms in
 * service definitions in yml files, but is is primarily
 * targets at operations not based on compose, as compose
 * has expected behaviour, and is piped right through
 * the libCompose code.
 */

import (
	"strings"

	libCompose_config "github.com/docker/libcompose/config"
	libCompose_yaml "github.com/docker/libcompose/yaml"
)

/**
 * clean up a service based on this app
 */
func (project *ComposeProject) AlterService(service *libCompose_config.ServiceConfig) {
	project.alterService_ProjectNetwork(service)
	project.alterService_RewriteMappedVolumes(service)
}

// make sure that a service is using the default network [@TODO THIS SHOULD NOT BE NECESSARY]
func (project *ComposeProject) alterService_ProjectNetwork(service *libCompose_config.ServiceConfig) {
	/**
	 * If a service has no network then we create the default network config.
	 *
	 * This is copypasta from github.com/docker/libcompose/project::Project.handleNetworkConfig()
	 * which means that we are duplicating internal functionality that may not be stable.
	 *
	 * This requirement came up after an update to the libcompose upstream library, which broke
	 * the existing missing network setup.  What is happening is that we are alteting a serviceconfig
	 * which will be added to a libcompose.Project::Project struct, and that struct has already
	 * run its initializer, which does the default network configuration.  This means that it is
	 * too late to simply add an empty network.  An alternative is to re-run the initializer, but
	 * as we have no access to that functionality from the Interface, there is little we can do.
	 */

	if service.Networks == nil || len(service.Networks.Networks) == 0 {
		// Add default as network
		service.Networks = &libCompose_yaml.Networks{
			Networks: []*libCompose_yaml.Network{
				{
					Name:     "default",
					RealName: project.composeContext.Context.ProjectName + "_default",
				},
			},
		}
	}
}

// rewrite mapped service volumes to use app points.
func (project *ComposeProject) alterService_RewriteMappedVolumes(service *libCompose_config.ServiceConfig) {

	for index, _ := range service.Volumes.Volumes {
		volume := service.Volumes.Volumes[index]

		switch volume.Source[0] {

		// relate volume to the current user home path
		case []byte("~")[0]:
			homePath := project.pathSettings.UserHomePath
			volume.Source = strings.Replace(volume.Source, "~", homePath, 1)

		// relate volume to project root
		case []byte(".")[0]:
			appPath := project.pathSettings.ProjectRootPath
			volume.Source = strings.Replace(volume.Source, "~", appPath, 1)

		// @TODO this is a stupid special hard-code that we should document somehow
		// @NOTE this is dangerous and will likely only work in cases where PWD is available
		case []byte("!")[0]:
			appPath := project.pathSettings.ExecPath
			volume.Source = strings.Replace(volume.Source, "!", appPath, 1)

		case []byte("@")[0]:
			if aliasPath, found := project.pathSettings.ConfigPaths.Get(volume.Source[1:]); found {
				volume.Source = strings.Replace(volume.Source, volume.Source, aliasPath.PathString(), 1)
			}
		}
	}

}
