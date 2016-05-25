package components

import (
	"encoding/json"
	"log"

	"github.com/hashicorp/consul/api"
)

const consulPath = "agilestack/components"

var (
	consulAddress = "consul.agilestacknet:8500"
)

// Change the address of Consul.
// Default is "consul.agilestacknet:8500"
func SetConsulAddress(consultAddress string) {
	consulAddress = consultAddress
}

func StoreComponents(components *Components) error {
	if err := Validate(components); err != nil {
		return err
	}

	client, err := consulClient()
	if err != nil {
		return err
	}

	kv := client.KV()

	bytes, err := json.Marshal(components)
	if err != nil {
		log.Println("Error while marshalling the components:", err)
		return err
	}

	pair := &api.KVPair{
		Key:   consulPath,
		Value: bytes,
	}
	_, err = kv.Put(pair, &api.WriteOptions{})
	return err
}

func GetComponents() (*Components, error) {
	client, err := consulClient()
	if err != nil {
		return nil, err
	}

	pair, _, err := client.KV().Get(consulPath, nil)
	if err != nil {
		log.Println("Error while loading components from Consul:", err)
		return nil, err
	}

	if pair == nil {
		return nil, nil
	}

	var components Components
	err = json.Unmarshal(pair.Value, &components)
	if err != nil {
		return nil, err
	}
	return &components, nil
}

func DeleteComponents() error {
	client, err := consulClient()
	if err != nil {
		return err
	}

	_, err = client.KV().Delete(consulPath, nil)
	return err
}

func consulClient() (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = consulAddress
	return api.NewClient(config)
}
