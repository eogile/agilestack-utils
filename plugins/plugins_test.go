package plugins

import (
	"github.com/eogile/agilestack-utils/dockerclient"
	"github.com/eogile/agilestack-utils/slices"
	"github.com/fsouza/go-dockerclient"
	"github.com/hashicorp/consul/api"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

var dockerClient *dockerclient.DockerClient
var consulTestClient *api.Client
var consulTestStorageClient *ConsulStorageClient

const (
	consulContainerName = "consulTest"
	consulImageName     = "gliderlabs/consul-server"
	network             = "networkPluginsTest"
	consulTestPort      = "8501"
	consulTestHost      = "127.0.0.1:" + consulTestPort
)

func setupConsulContainer() {

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
	//setup consul ENV : CONSUL_HTTP_ADDR in form of 127.0.0.1:8501
	//needed later for consul Client
	os.Setenv("CONSUL_HTTP_ADDR", consulTestHost)
	consulTestClient, err = api.NewClient(api.DefaultConfig())
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

	consulTestStorageClient = &ConsulStorageClient{consulTestClient}
}

func tearDownConsulContainer() {
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

func setupResources() {
	setupConsulContainer()

}

func tearDownResources() {
	tearDownConsulContainer()
}

func TestMain(m *testing.M) {
	log.Println("Launching tests agilestack/utils/plugins")

	/*
	 * Docker client for test utilities
	 */
	dockerClient = dockerclient.NewClient()

	/*
	 * Create agilestacknet docker network if not exists
	 */
	networks, err := dockerClient.ListNetworks()
	if err != nil {
		log.Println("unable to List docker networks : ", err)
		os.Exit(1)
	}

	options := docker.CreateNetworkOptions{
		Name:   network,
		Driver: "bridge",
	}
	if !slices.DockerNetworkInSlice(network, networks) {
		//create network
		_, errNet := dockerClient.CreateNetwork(options)
		if errNet != nil {
			log.Printf("Cannot create docker network %v. Got error %v", network, errNet)
		}
	}
	setupResources()

	exitCode := m.Run()
	tearDownResources()

	os.Exit(exitCode)
}
