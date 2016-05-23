package menu

import (
	"encoding/json"
	"log"

	"github.com/hashicorp/consul/api"
)

const menuPrefix = "agilestack/menu/"

type (
	MenuStore interface {
		// Stores the menu and returns the result error
		StoreMenu(menu *Menu) error

		// Lists all the existing menus
		ListMenus() ([]Menu, error)

		// Deletes the menu matching the given plugin name.
		DeleteMenu(pluginName string) error
	}

	ConsulMenuStore struct {
		ConsulClient *api.Client
	}
)

func newMenuStore() (MenuStore, error) {
	config := api.DefaultConfig()
	config.Address = consulAddress
	client, err := api.NewClient(config)
	if err != nil {
		log.Println("Got error when trying to create consulClient", err)
		return nil, err
	}
	return &ConsulMenuStore{
		ConsulClient: client,
	}, nil
}

func (store *ConsulMenuStore) StoreMenu(menu *Menu) error {
	if err:= ValidateMenu(menu); err != nil {
		log.Println("Menu is invalid:", err)
		return err
	}

	kv := store.ConsulClient.KV()

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

func (store *ConsulMenuStore) ListMenus() ([]Menu, error) {
	kv := store.ConsulClient.KV()

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

func (store *ConsulMenuStore) DeleteMenu(pluginName string) error {
	kv := store.ConsulClient.KV()
	_, err := kv.Delete(menuPrefix+pluginName, nil)
	return err
}
