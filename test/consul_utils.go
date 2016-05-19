package test

import (
	"log"
	"sync"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/hashicorp/consul/api"
	"github.com/eogile/agilestack-utils/dockerclient"
)

const (
	consulContainerName = "consulTest"
	consulImageName     = "gliderlabs/consul-server"
	network             = "networkPluginsTest"
	consulTestPort      = "8501"
	consulTestHost      = "127.0.0.1:" + consulTestPort
)

/*
Starts a consul container for tests purposes.
 */
func StartConsulContainer(dockerClient *dockerclient.DockerClient) error {
	/*
	 * Container configuration
	 */
	containerConfig := docker.Config{
		Image: consulImageName,
		Tty:   true,
	}
	containerConfig.Cmd = append(containerConfig.Cmd, "-bootstrap")
	//containerConfig.Cmd = append(containerConfig.Cmd, "-ui")
	containerConfig.ExposedPorts = map[docker.Port]struct{}{
		"8500/tcp": {},
	}

	portBindings := map[docker.Port][]docker.PortBinding{
		"8500/tcp": []docker.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: consulTestPort,
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
		Name:       consulContainerName,
		Config:     &containerConfig,
		HostConfig: &hostConfig,
	}

	err := dockerClient.LaunchContainer(containerOptions)
	if err != nil {
		log.Fatalln("Cannot launch consul : ", err)
	}

	consulTestClient, err := NewConsulClient()
	//test if consul callable
	if err != nil {
		log.Fatalln("Unable to setup Consul client, ", err)
	}
	//wait for condition consul is available
	var mutex sync.Mutex
	conditionIsAvailable := sync.NewCond(&mutex)
	go func() {
		log.Println("In go routine")
		for i := 0; i < 20; i++ {
			log.Println("In go routine loop : ", i)
			health := consulTestClient.Health()
			healthCheck, _, err := health.State(api.HealthPassing, nil)
			if err != nil {
				log.Println("Got error when checking consul health", err)
				//mutex.Lock()
				//conditionIsAvailable.Signal()
				//mutex.Unlock()
				//t.Fatal("cannot check health", err)
			}
			for _, h := range healthCheck {
				status := h.Status
				log.Println("got consul health status : ", status)
				if status == api.HealthPassing {
					mutex.Lock()
					conditionIsAvailable.Signal()
					mutex.Unlock()
				}
			}
			time.Sleep(time.Second)
		}
		mutex.Lock()
		conditionIsAvailable.Signal()
		mutex.Unlock()
		log.Fatalln("Health never came good", err)

	}()

	mutex.Lock()
	conditionIsAvailable.Wait()
	mutex.Unlock()

	return nil
}

/*
Stops and removes the Consul container used in tests.
 */
func RemoveConsulContainer(dockerClient *dockerclient.DockerClient) {
	isRunning, err := dockerClient.IsContainerRunning(consulContainerName)
	if err != nil {
		log.Fatalln("in tearDownConsulContainer, unable to do IsContainerRunning ", err)
	}
	container := dockerClient.GetContainerByName(consulContainerName)
	if container != nil {
		if isRunning {
			err := dockerClient.StopContainer(container.ID, 0)
			if err != nil {
				log.Println("Strange behaviour, says that cannot stop consul, but has been killed : ", err)
			}
		}
		isRunning, err = dockerClient.IsContainerRunning(consulContainerName)
		err = dockerClient.RemoveContainer(docker.RemoveContainerOptions{ID: consulContainerName, RemoveVolumes: false, Force: false})
		if err != nil {
			log.Fatalln("Cannot remove consul : ", err)
		}

	}
}

/*
Returns a Consul client configured to access the Consul container deployed in tests.
 */
func NewConsulClient() (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = consulTestHost
	return api.NewClient(config)
}