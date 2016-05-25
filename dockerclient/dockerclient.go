package dockerclient

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"

	"strings"

	"github.com/eogile/agilestack-utils/slices"
	"github.com/fsouza/go-dockerclient"
	//"github.com/eogile/go-dockerclient"
	"errors"
)

/*DockerClient : the struct supporting the dockerClient
 *
 */
type DockerClient struct {
	*docker.Client
}

//NewClient : create a *DockerClient
func NewClient() *DockerClient {
	endPoint := dockerEndPoint()
	return createDockerClient(endPoint)
}

func dockerEndPoint() string {
	dockerEndpoint := os.Getenv("DOCKER_HOST") //DOCKER_HOST
	if dockerEndpoint == "" && runtime.GOOS == "darwin" {
		//log.Fatal("DOCKER_HOST should be set on a MacOSX\n")
		dockerEndpoint = "unix:///var/run/docker.sock"
	}
	if dockerEndpoint == "" && runtime.GOOS == "linux" {
		/*
		 * We suppose we are inside a container and the host is MacOsx
		 * with boot2docker so we'll try to use tcp://172.17.42.1:2376.
		 */
		dockerEndpoint = "unix:///var/run/docker.sock"
	}
	if dockerEndpoint == "" {
		log.Fatalf("You have to provide a not empty docker endpoint through command ligne")
	}

	return dockerEndpoint
}

func createDockerClient(endPoint string) *DockerClient {
	var client *docker.Client
	var errDocker error

	if runtime.GOOS == "darwin" && endPoint != "unix:///var/run/docker.sock" {
		me := currentUser()
		home := me.HomeDir
		log.Printf("HomeDir = %s", home)
		path := os.Getenv("DOCKER_CERT_PATH")
		log.Printf("dockerCertPath:%s", path)
		ca := fmt.Sprintf("%s/ca.pem", path)
		cert := fmt.Sprintf("%s/cert.pem", path)
		key := fmt.Sprintf("%s/key.pem", path)
		log.Printf("In NewServer, endpoint %v, cert %v, key %v, ca %v",
			endPoint, cert, key, ca)
		client, errDocker = docker.NewTLSClient(endPoint, cert, key, ca)
	} else {
		//we suppose we are inside Linux container
		client, errDocker = docker.NewClient(endPoint)
	}

	if errDocker != nil {
		log.Fatalf("ERROR : unable to create docker client to host endpoint %v : %v",
			endPoint, errDocker)
	}
	return &DockerClient{client}
}

func currentUser() *user.User {
	me, errUser := user.Current()
	if errUser != nil {
		log.Fatalf("ERROR : impossible to get current user %v", errUser)
	}
	return me
}

/*ListMainImages : Returns the list of top-level Docker images.
 */
func (dockerClient *DockerClient) ListMainImages() ([]docker.APIImages, error) {
	listOps := docker.ListImagesOptions{All: false}
	return dockerClient.ListImages(listOps)
}

/*ListAllImages : Returns the list of all Docker images.
 */
func (dockerClient *DockerClient) ListAllImages() ([]docker.APIImages, error) {
	listOps := docker.ListImagesOptions{All: true}
	return dockerClient.ListImages(listOps)
}

/*ListRunningContainers : Returns the list of all running containers.
 */
func (dockerClient *DockerClient) ListRunningContainers() ([]docker.APIContainers, error) {
	listOps := docker.ListContainersOptions{All: false}
	return dockerClient.ListContainers(listOps)
}

/*ListAllContainers : Returns the list of all containers, including those not running
 */
func (dockerClient *DockerClient) ListAllContainers() ([]docker.APIContainers, error) {
	listOps := docker.ListContainersOptions{All: true}
	return dockerClient.ListContainers(listOps)
}

/*IsContainerInstalled : Returns true if a container is installed, given a name
 */
func (dockerClient *DockerClient) IsContainerInstalled(name string) (bool, error) {
	/*
	 * Listing all the containers, even the stopped ones.
	 */
	containers, err := dockerClient.ListAllContainers()

	if err != nil {
		return false, err
	}

	return containersArrayContains(containers, name), nil
}

/*IsContainerRunning : Returns true if a container is running, given a name
 */
func (dockerClient *DockerClient) IsContainerRunning(name string) (bool, error) {
	/*
	 * Listing the running containers
	 */
	containers, err := dockerClient.ListRunningContainers()
	log.Printf("In IsContainerRunning, got %v running containers", len(containers))

	if err != nil {
		return false, err
	}

	return containersArrayContains(containers, name), nil
}

func containersArrayContains(containers []docker.APIContainers, containerName string) bool {
	for _, container := range containers {
		for _, name := range container.Names {
			if ExtractContainerName(name) == containerName {
				return true
			}
		}
	}
	return false
}

// GetImageByName : returns a *docker.APIImages given a name (string)
func (dockerClient *DockerClient) GetImageByName(name string) *docker.APIImages {

	/*
	 * Listing the images
	 */
	images, imageErr := dockerClient.ListMainImages()
	if imageErr != nil {
		log.Printf("Error when listing Docker images : %v", imageErr)
		return nil
	}
	if images == nil || len(images) == 0 {
		log.Print("ListMain Images did not found any images")
		return nil
	}

	for _, image := range images {
		log.Printf("image.RepoTags[0] %s, image.ID=%s", image.RepoTags[0], image.ID)
		for _, repoTag := range image.RepoTags {
			if repoTag == name + ":latest" {
				return &image
			}
		}
	}
	return nil
}

// GetContainerByName : get a container by a name
func (dockerClient *DockerClient) GetContainerByName(name string) *docker.APIContainers {

	/*
	 * Listing the containers
	 */
	containers, containerErr := dockerClient.ListAllContainers()
	if containerErr != nil {
		log.Printf("Error when listing Docker containers : %v", containerErr)
		return nil
	}
	if containers == nil || len(containers) == 0 {
		log.Print("ListAllContainers Images did not found any images")
		return nil
	}

	for _, container := range containers {
		for _, installedName := range container.Names {
			if ExtractContainerName(installedName) == name {
				return &container
			}
		}
	}
	return nil
}

/*LaunchContainer : Launch a container
 * if the container already exists, error
 * if the container is already launched, nothing is done
 */
func (dockerClient *DockerClient) LaunchContainer(options docker.CreateContainerOptions) error {
	containerName := options.Name
	isRunning, errIsRunning := dockerClient.IsContainerRunning(containerName)
	if errIsRunning != nil {
		return errIsRunning
	}
	if isRunning {
		return errors.New(fmt.Sprintf("in LaunchContainer, container %s is already running", containerName))
	}

	isInstalled, errIsInstalled := dockerClient.IsContainerInstalled(containerName)
	if errIsInstalled != nil {
		return errIsInstalled
	}
	if isInstalled {
		return errors.New(fmt.Sprintf("in LaunchContainer, container %s is already installed (probably not running)", containerName))
	}

	imageName := options.Config.Image
	log.Printf("Creating container %s based on image %s", containerName, imageName)

	/*
	 * Finding the Docker image matching the plugin's name
	 */
	image := dockerClient.GetImageByName(imageName)
	if image == nil {
		//create image
		log.Printf("in LaunchContainer, did not found local image %s, pulling from registry", imageName)
		err := dockerClient.PullImage(docker.PullImageOptions{Repository: imageName, Tag: "latest"}, docker.AuthConfiguration{})
		if err != nil {
			log.Printf("Could not pull the image %s, got error %v\n", imageName, err)
		}

	}

	container, err := dockerClient.CreateContainer(options)
	if err != nil {
		log.Printf("Error on createContainer : %v", err)
		return err
	}
	if container == nil {
		return errors.New(fmt.Sprintf("in LaunchContainer, tried to create %s container but got nil container", containerName))
	}
	id := container.ID
	log.Printf("container created with ID: %s ", id)

	err = dockerClient.StartContainer(id, nil)
	if err != nil {
		log.Printf("Error on startContainer : %v", err)
		return err
	}
	attachContainerOptions := docker.AttachToContainerOptions{
		Container: containerName,
	}
	dockerClient.AttachToContainer(attachContainerOptions)

	return nil

}

//LaunchContainerByNames : Launch a container
// if the container already exists, launch it without re-creating it
// if the container is already launched, nothing is done
//DEPRECATED : in favour of LaunchContainer
func (dockerClient *DockerClient) LaunchContainerByNames(imageName string, containerName string, network string, cmd string) error {

	isRunning, errIsRunning := dockerClient.IsContainerRunning(containerName)
	if errIsRunning != nil {
		return errIsRunning
	}
	if isRunning {
		log.Println("in LaunchContainer, container is already running")
		return nil
	}

	isInstalled, errIsInstalled := dockerClient.IsContainerInstalled(containerName)
	if errIsInstalled != nil {
		return errIsInstalled
	}
	var id string
	if isInstalled {
		id = dockerClient.GetContainerByName(containerName).ID
	} else {

		log.Printf("Creating container %s based on image %s", containerName, imageName)

		/*
		 * Finding the Docker image matching the plugin's name
		 */
		image := dockerClient.GetImageByName(imageName)
		if image == nil {
			//create image
			log.Printf("in LaunchContainer, did not found local image %s, pulling from registry", imageName)
			err := dockerClient.PullImage(docker.PullImageOptions{Repository: imageName, Tag: "latest"}, docker.AuthConfiguration{})
			if err != nil {
				log.Printf("Could not pull the image %s, got error %v\n", imageName, err)
			}

		}

		/*
		 * Container configuration
		 */
		containerConfig := docker.Config{
			Image: imageName,
			//AttachStdin: true,
			Tty: true,
		}
		if len(cmd) > 0 {
			containerConfig.Cmd = append(containerConfig.Cmd, cmd)
		}

		/*
		 * Host configuration
		 */
		hostConfig := docker.HostConfig{
			PublishAllPorts: true,
			Binds:           []string{"agilestack-shared:/shared"},
			NetworkMode:     network,
		}

		containerOptions := docker.CreateContainerOptions{
			Name:       containerName,
			Config:     &containerConfig,
			HostConfig: &hostConfig,
		}

		var err error
		container, err := dockerClient.CreateContainer(containerOptions)
		if err != nil {
			log.Printf("Error on createContainer : %v", err)
			return err
		}
		id = container.ID
		log.Printf("container created with ID: %s ", container.ID)
	}
	err := dockerClient.StartContainer(id, nil)
	if err != nil {
		log.Printf("Error on startContainer : %v", err)
		return err
	}
	attachContainerOptions := docker.AttachToContainerOptions{
		Container: containerName,
	}
	dockerClient.AttachToContainer(attachContainerOptions)
	return nil
}

/*CreateNetWorkIfNeeded : Create docker network if not exists
 *
 */
func (dockerClient *DockerClient) CreateNetWorkIfNeeded(networkName string) {
	networks, err := dockerClient.ListNetworks()
	if err != nil {
		log.Println("unable to List docker networks : ", err)
		os.Exit(1)
	}

	options := docker.CreateNetworkOptions{
		Name:   networkName,
		Driver: "bridge",
	}
	if !slices.DockerNetworkInSlice(networkName, networks) {
		//create network
		_, errNet := dockerClient.CreateNetwork(options)
		if errNet != nil {
			log.Printf("Cannot create docker network %v. Got error %v", networkName, errNet)
		}
	}
}

/*ExtractContainerName :extract the name of the container matching the given Docker containers's name.
 *
 * Please notice that in the case where the image is referenced by several
 * repositories, the first one is used to deduce the plugin's name.
 *
 * Examples :
 * - /agilestack-proxy => agilestack-proxy
 */
func ExtractContainerName(imageName string) string {
	if strings.Contains(imageName, "/") {
		items := strings.Split(imageName, "/")
		imageName = items[len(items)-1]
	}

	if strings.Contains(imageName, ":") {
		items := strings.Split(imageName, ":")
		imageName = items[0]
	}

	return imageName
}
