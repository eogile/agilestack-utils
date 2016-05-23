package menu

import (
	"github.com/hashicorp/consul/api"
	"log"
	"encoding/json"
)

var consulAddress = "consul.agilestacknet:8500"
const menuPrefix = "agilestack/menu/"

func SetConsulAddress(consultAddress string) {
	consulAddress = consultAddress
}

// Stores the menu and returns the result error.
func StoreMenu(menu *Menu) error {
	client, err := client()
	if err != nil {
		return err
	}

	if err:= ValidateMenu(menu); err != nil {
		log.Println("Menu is invalid:", err)
		return err
	}

	kv := client.KV()

	bytes, err := json.Marshal(menu)
	if err != nil {
		log.Println("Error while marshalling the menu:", err)
	}

	pair := &api.KVPair{
		Key:   menuPrefix + menu.PluginName,
		Value: bytes,
	}
	_, err = kv.Put(pair, &api.WriteOptions{})
	return err
}

// Lists all the existing menus.
func  ListMenus() ([]Menu, error) {
	client, err := client()
	if err != nil {
		return nil, err
	}

	kv := client.KV()

	pairs, _, err := kv.List(menuPrefix, nil)
	if err != nil {
		log.Println("Error while listing menu entries in Consul:", err)
		return nil, err
	}

	menus := make([]Menu, 0)
	for _, menuPair := range pairs {
		menu := Menu{}
		err = json.Unmarshal(menuPair.Value, &menu)
		if err != nil {
			log.Printf("Error while unmarshalling menu '%s': %v", menuPair.Key, err)
			return nil, err
		}
		menus = append(menus, menu)
	}
	return menus, nil
}

// Deletes the menu matching the given plugin name.
func DeleteMenu(pluginName string) error {
	client, err := client()
	if err != nil {
		return err
	}

	kv := client.KV()
	_, err = kv.Delete(menuPrefix+pluginName, nil)
	return err
}

func client() (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = consulAddress
	return api.NewClient(config)
}