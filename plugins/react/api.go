package react

import (
	"encoding/json"
	"log"

	"github.com/hashicorp/consul/api"
)

const consulPrefix = "agilestack/react/"

var consulAddress = "consul.agilestacknet:8500"

func SetConsulAddress(consultAddress string) {
	consulAddress = consultAddress
}

func StoreRoutesAndReducers(config *PluginConfiguration) error {
	client, err := store()
	if err != nil {
		return err
	}

	err = Validate(config)
	if err != nil {
		return err
	}

	kv := client.KV()

	bytes, err := json.Marshal(config)
	if err != nil {
		log.Println("Error while marshalling the menu:", err)
	}

	pair := &api.KVPair{
		Key:   consulPrefix + config.PluginName,
		Value: bytes,
	}
	_, err = kv.Put(pair, &api.WriteOptions{})
	return err
}

func ListRoutesAndReducers() ([]PluginConfiguration, error) {
	client, err := store()
	if err != nil {
		return nil, err
	}

	pairs, _, err := client.KV().List(consulPrefix, nil)
	if err != nil {
		log.Println("Error while listing plugins configurations in Consul:", err)
		return nil, err
	}

	configurations := make([]PluginConfiguration, 0)
	for _, pair := range pairs {
		config := PluginConfiguration{}
		err = json.Unmarshal(pair.Value, &config)
		if err != nil {
			log.Printf("Error while unmarshalling the configuration '%s': %v", pair.Key, err)
			return nil, err
		}
		configurations = append(configurations, config)
	}
	return configurations, nil
}

func DeleteRoutesAndReducers(pluginName string) error {
	client, err := store()
	if err != nil {
		return err
	}

	kv := client.KV()
	_, err = kv.Delete(consulPrefix+pluginName, nil)
	return err
}

func store() (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = consulAddress
	return api.NewClient(config)
}
