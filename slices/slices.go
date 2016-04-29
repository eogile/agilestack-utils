package slices

import (
	//"github.com/eogile/go-dockerclient"
	"github.com/fsouza/go-dockerclient"
)

//
//function to check if in string is in a slice of string
//
func StringInSlice(a string, list []string) bool {
	if len(list) == 0 {
		return false
	}
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//
//function to check if a docker netWork is in a slice of docker networks
//Based on network name
//
func DockerNetworkInSlice(a string, list []docker.Network) bool {
	if a == "" {
		return false
	}
	if len(list) == 0 {
		return false
	}
	for _, b := range list {
		if b.Name == a {
			return true
		}
	}
	return false
}
