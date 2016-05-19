package test

import (
	"log"
	"os"
	"testing"

	"github.com/eogile/agilestack-utils/dockerclient"
	"github.com/eogile/agilestack-utils/test"
	"github.com/hashicorp/consul/api"
)

var dockerClient *dockerclient.DockerClient
var ConsulTestClient *api.Client

const (
	network = "networkPluginsTest"
)

/*
Starts the Consul container and initializes the Consul client.
 */
func setupResources() {
	test.StartConsulContainer(dockerClient)
	ConsulTestClient, _ = test.NewConsulClient()
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
	err := test.CreateNetworkIfNotExists(dockerClient, network)
	if err != nil {
		log.Fatalln("Unable to create a Docker network:", err)
	}

	setupResources()

	exitCode := m.Run()
	tearDownResources()

	os.Exit(exitCode)
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}
