package plugins

import (
	"encoding/json"
	"github.com/hashicorp/consul/api"
	"log"
	"errors"
)

const resourcesPrefix = "agilestack/security/resources/"

//PluginInfoStorageClient is an interface fo operations to get and store informations about plugins
type PluginInfoStorageClient interface {
	//StoreResource insert or modify information about the security Resources provided by a plugin
	StoreResource(resource Resource) error

	//GetResource retrieve a Resource by its name
	GetResource(name string) (*Resource, error)

	//ListResources retrieve a Resource by its name
	ListResources() ([]Resource, error)

	//DeleteResource delete a resource given its name
	DeleteResource(name string) error
}

type ConsulStorageClient struct {
	consulClient *api.Client
}

// NewPluginInfoStorageClient returns a fresh ConsulStorageClient
func NewPluginInfoStorageClient() PluginInfoStorageClient {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Println("Got error when trying to create consulClient", err)
		return nil
	}
	return &ConsulStorageClient{client}
}

//StoreResource store a resource in consul Store
func (c *ConsulStorageClient) StoreResource(resource Resource) error {
	kv := c.consulClient.KV()

	resourceBytes, errJson := json.Marshal(resource)
	if errJson != nil {
		return errJson
	}
	if resource.Key == "" {
		return errors.New("Unable to store a resource with an empty name")
	}
	p := &api.KVPair{Key: resourcesPrefix + resource.Key, Value: resourceBytes}
	_, err := kv.Put(p, &api.WriteOptions{})
	if err != nil {
		return err
	}
	return nil
}

//GetResource retrieve a resource in consul store given a name
func (c *ConsulStorageClient) GetResource(name string) (*Resource, error) {
	kv := c.consulClient.KV()

	pair, _, err := kv.Get(resourcesPrefix+name, nil)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, nil
	}
	resource := &Resource{}
	errJson := json.Unmarshal(pair.Value, resource)
	return resource, errJson
}

// List all the resources
func (c *ConsulStorageClient) ListResources() ([]Resource, error) {
	kv := c.consulClient.KV()

	pairs, _, err := kv.List(resourcesPrefix, nil)
	if err != nil {
		return nil, err
	}
	resources := make([]Resource, len(pairs))
	for i, pair := range pairs {
		resource := Resource{}
		errJson := json.Unmarshal(pair.Value, &resource)
		if errJson != nil {
			return nil, errJson
		}
		resources[i] = resource
	}
	return resources, nil

}

//DeleteResource Delete a resource in consul Store given its name
func (c *ConsulStorageClient) DeleteResource(name string) error {
	kv := c.consulClient.KV()

	_, err := kv.Delete(resourcesPrefix+name, nil)
	return err
}
