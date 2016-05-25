package auth_test

import (
	"log"
	"os"
	"testing"

	"github.com/eogile/agilestack-utils/dockerclient"
	"github.com/eogile/agilestack-utils/test"
)

const (
	network = "network-test"
)

func TestMain(m *testing.M) {
	/*
		Docker client for test utilities
	*/
	dockerClient := dockerclient.NewClient()

	/*
		Creating the Docker network if it does not exist.
	*/
	dockerClient.CreateNetWorkIfNeeded(network)

	/*
	 * Bootstrap the PostgreSQL and the Hydra containers.
	 */
	if err := test.StartHydraContainer(dockerClient, network); err != nil {
		log.Fatalln("Error while starting PostgreSQL and Hydra containers", err)
	}

	exitCode := m.Run()

	if err := test.StopHydraContainer(dockerClient); err != nil {
		log.Fatalln("Error while stopping PostgreSQL and Hydra containers:", err)
	}

	os.Exit(exitCode)
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func SetUp(t *testing.T) {
}
