package dockerclient

import (
	"log"
	"os"
	"testing"

	"github.com/fsouza/go-dockerclient"
)

var dockerClient *DockerClient

const (
	network           = "agilestacknet_test_dockerclient"
	testImageName     = "busybox"
	testContainerName = "smallContainerTest"
)

func TestExtractContainerName(t *testing.T) {
	//* Examples :
	//* - /agilestack-proxy => agilestack-proxy
	const name = "/agilestack-proxy"
	const expectedName = "agilestack-proxy"

	if expectedName != ExtractContainerName(name) {
		t.Errorf("Expected %s, got %s", expectedName, name)
	}

}

func TestSmallImagePresent(t *testing.T) {
	err := dockerClient.LaunchContainerByNames(testImageName, testContainerName, network, "")
	if err != nil {
		t.Fatal("got error when launching container, ", err)
	}
	isRunning, errIs := dockerClient.IsContainerRunning(testContainerName)
	if errIs != nil {
		t.Fatal("unable to check id test container is running", errIs)
	}
	if !isRunning {
		t.Errorf("expected container %s running, but it is not", testContainerName)
	}
	container := dockerClient.GetContainerByName(testContainerName)
	if container == nil {
		t.Error("Expected to find by name container: ", testContainerName)
	} else {
		t.Logf("Running container: %v", container)
		errStop := dockerClient.StopContainer(container.ID, 0) //10 seconds
		if errStop != nil {
			dockerClient.Logs(docker.LogsOptions{Container: testContainerName, Stderr: true})
			t.Fatalf("Unable to stop container %s, : %v", testContainerName, errStop)
		}
		//container should be installed but not running
		isInstalled, errInstalled := dockerClient.IsContainerInstalled(testContainerName)
		if errInstalled != nil {
			t.Fatal("Unable to check if container is installed", errInstalled)
		}
		if !isInstalled {
			t.Errorf("%s should have been installed but is not", testContainerName)
		} else {
			//should not be running now
			is2ndRunning, err2ndRunning := dockerClient.IsContainerRunning(testContainerName)
			if err2ndRunning != nil {
				t.Fatal("Got error when checking if running", err2ndRunning)
			}
			if is2ndRunning {
				t.Errorf("%s expected to not be running, but is running", testContainerName)
			}
		}
		errRemove := dockerClient.RemoveContainer(docker.RemoveContainerOptions{container.ID, false, false})
		if errRemove != nil {
			t.Fatal("Error removing the container", errRemove)
		}
		is2ndInstalled, err2ndInstalled := dockerClient.IsContainerInstalled(testContainerName)
		if err2ndInstalled != nil {
			t.Fatal("Unable to check if container is installed", err2ndInstalled)
		}
		if is2ndInstalled {
			t.Errorf("%s should not be installed but is", testContainerName)
		}
	}

}

//Test launch with cmd
//

func TestMain(m *testing.M) {
	log.Println("Launching tests agilestack/utils/dockerClient")
	/*
	 * Docker client for test utilities
	 */
	dockerClient = NewClient()

	dockerClient.CreateNetWorkIfNeeded(network)

	exitCode := m.Run()
	os.Exit(exitCode)
}
