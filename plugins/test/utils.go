package test

import (
	"log"
	"os"
	"testing"

	"github.com/eogile/agilestack-utils/dockerclient"
	"github.com/eogile/agilestack-utils/test"
)

var dockerClient *dockerclient.DockerClient

const (
	network = "networkPluginsTest"
)

/*
Starts the Consul container and initializes the Consul client.
 */
func setupResources() {
	test.StartConsulContainer(dockerClient, network)
}

/*
Stops and removes the Consul container.
 */
func tearDownResources() {
	test.RemoveConsulContainer(dockerClient)
}

func DoTestMain(m *testing.M) {
	/*
	Docker client for test utilities
	 */
	dockerClient = dockerclient.NewClient()

	/*
	Creating the Docker network if it does not exist.
	 */
	dockerClient.CreateNetWorkIfNeeded(network)

	setupResources()

	exitCode := m.Run()
	tearDownResources()

	os.Exit(exitCode)
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}
