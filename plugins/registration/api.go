package registration

import (
	"encoding/json"
	"log"

	"github.com/hashicorp/consul/api"
	"net/http"
	"errors"
	"strconv"
)

const consulPrefix = "agilestack/registration/"

var (
	consulAddress = "consul.agilestacknet:8500"
	appBuilderAddress = "http://agilestack-root-app-builder:8080"
)

// Change the address of Consul.
// Default is "consul.agilestacknet:8500"
func SetConsulAddress(consultAddress string) {
	consulAddress = consultAddress
}

// Change the address of the application builder.
// Default is "http://agilestack-root-app-builder:8080".
func SetAppBuilderAddress(address string) {
	appBuilderAddress = address
}

// Stores the routes and the reducers into Consul store.
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

// Loads all the plugins configuration from Consul store.
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

// Deletes from Consul store all the routes and the reducers of the given plugin.
func DeleteRoutesAndReducers(pluginName string) error {
	client, err := store()
	if err != nil {
		return err
	}

	kv := client.KV()
	_, err = kv.Delete(consulPrefix + pluginName, nil)
	return err
}

// Launches the application build.
func LaunchApplicationBuild() error {
	response, err := http.Post(appBuilderAddress + "/plugins", "application/json", nil)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New("Invalid status code: " + strconv.Itoa(response.StatusCode))
	}
	return nil
}

func store() (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = consulAddress
	return api.NewClient(config)
}
