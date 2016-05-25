package test

import (
	"errors"
	"log"
	"time"

	"github.com/eogile/agilestack-utils/dockerclient"
	"github.com/fsouza/go-dockerclient"
)

const (
	postgreSQLPort          = "5432"
	postgreSQLContainerName = "postgresql-test"
	hydraPort               = "9090"
	hydraContainerName      = "hydra-test"
)

func StartHydraContainer(dockerClient *dockerclient.DockerClient, network string) error {
	if err := startPostgreSQLContainer(dockerClient, network); err != nil {
		return err
	}
	return startHydra(dockerClient, network, 0)
}

func StopHydraContainer(dockerClient *dockerclient.DockerClient) error {
	if err := stopContainer(dockerClient, postgreSQLContainerName); err != nil {
		return err
	}
	return stopContainer(dockerClient, hydraContainerName)
}

// Starts a PostgreSQL container.
func startPostgreSQLContainer(dockerClient *dockerclient.DockerClient, network string) error {
	containerConfig := docker.Config{
		Image: "postgres:9.5",
		Tty:   true,
		Env: []string{
			"POSTGRES_USER=hydra",
			"POSTGRES_PASSWORD=hydra_agilestack_test",
			"POSTGRES_DB=hydra",
		},
		ExposedPorts: map[docker.Port]struct{}{
			"5432/tcp": {},
		},
	}

	portBindings := map[docker.Port][]docker.PortBinding{
		"5432/tcp": []docker.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: postgreSQLPort,
			},
		},
	}

	/*
	 * Host configuration
	 */
	hostConfig := docker.HostConfig{
		PublishAllPorts: true,
		NetworkMode:     network,
		PortBindings:    portBindings,
	}

	containerOptions := docker.CreateContainerOptions{
		Name:       postgreSQLContainerName,
		Config:     &containerConfig,
		HostConfig: &hostConfig,
	}

	return dockerClient.LaunchContainer(containerOptions)
}

// Starts an "eogile/agilestack-hydra-host" container.
// As the PostgreSQL container may take some time to initialize itself, several
// attempts to start the Hydra container may be required (the container fails to
// start if it cannot connect the PostgreSQL database).
func startHydra(dockerClient *dockerclient.DockerClient, network string, counter int) error {
	log.Println("Hydra container launch attempt:", counter)
	if counter >= 15 {
		return errors.New("Max attempts threshold achieved")
	}

	containerConfig := docker.Config{
		Image: "eogile/agilestack-hydra-host",
		Tty:   true,
		Env: []string{
			"CLIENT_ID=superapp2",
			"CLIENT_SECRET=supersecret2",
			"SUPERACCOUNT_USERNAME=superadmin@eogile.com",
			"SUPERACCOUNT_SECRET=supersecret",
			"DATABASE_URL=postgres://hydra:hydra_agilestack_test@" +
				postgreSQLContainerName +
				"." +
				network +
				":5432/hydra?sslmode=disable",
		},
		ExposedPorts: map[docker.Port]struct{}{
			"9090/tcp": {},
		},
	}

	portBindings := map[docker.Port][]docker.PortBinding{
		"9090/tcp": []docker.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: hydraPort,
			},
		},
	}

	/*
	 * Host configuration
	 */
	hostConfig := docker.HostConfig{
		PublishAllPorts: true,
		NetworkMode:     network,
		PortBindings:    portBindings,
	}

	containerOptions := docker.CreateContainerOptions{
		Name:       hydraContainerName,
		Config:     &containerConfig,
		HostConfig: &hostConfig,
	}

	time.Sleep(500000000)
	if err := dockerClient.LaunchContainer(containerOptions); err != nil {
		return err
	}
	time.Sleep(500000000)
	if isRunning, err := dockerClient.IsContainerRunning(hydraContainerName); err != nil || isRunning {
		return err
	}

	if err := stopContainer(dockerClient, hydraContainerName); err != nil {
		return err
	}

	return startHydra(dockerClient, network, counter+1)
}

// Stops the container matching the given name.
// If there is no such container, nothing is done.
func stopContainer(dockerClient *dockerclient.DockerClient, containerName string) error {

	isRunning, err := dockerClient.IsContainerRunning(containerName)
	if err != nil {
		return err
	}

	container := dockerClient.GetContainerByName(containerName)
	if container != nil {
		if isRunning {
			err := dockerClient.StopContainer(container.ID, 0)
			if err != nil {
				log.Println("Strange behaviour, says that cannot stop container, but has been killed : ", err)
				return err
			}
		}

		return dockerClient.RemoveContainer(docker.RemoveContainerOptions{
			ID:            containerName,
			RemoveVolumes: true,
			Force:         false,
		})
	}
	return nil
}
