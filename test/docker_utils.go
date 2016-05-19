package test

import (
	"github.com/eogile/agilestack-utils/dockerclient"
	"github.com/eogile/agilestack-utils/slices"
	"github.com/fsouza/go-dockerclient"
)

/*
Creates a bridge Docker network matching the given name if there is no network
with this name.
*/
func CreateNetworkIfNotExists(dockerClient *dockerclient.DockerClient, networkName string) error {
	/*
	Listing Docker networks to check the presence of the searched network.
	*/
	networks, err := dockerClient.ListNetworks()
	if err != nil {
		return err
	}

	options := docker.CreateNetworkOptions{
		Name:   networkName,
		Driver: "bridge",
	}
	if !slices.DockerNetworkInSlice(networkName, networks) {
		_, errNet := dockerClient.CreateNetwork(options)
		return errNet
	}
	return nil
}
