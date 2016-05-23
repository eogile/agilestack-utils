package menu_test

import (
	"log"
	"testing"

	"github.com/eogile/agilestack-utils/plugins/menu"
	"github.com/eogile/agilestack-utils/plugins/test"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	menu1 = menu.Menu{
		PluginName: "MyWonderfulPlugin",
		Entries: []menu.MenuEntry{
			menu.MenuEntry{
				Name:   "Entry 1",
				Route:  "/entry-1",
				Weight: 10,
				Entries: []menu.MenuEntry{},
			},
			menu.MenuEntry{
				Name:   "Entry 2",
				Route:  "/entry-2",
				Weight: 11,
				Entries: []menu.MenuEntry{
					menu.MenuEntry{
						Name:   "Entry 2.1",
						Route:  "/entry-2-1",
						Weight: 10,
						Entries: []menu.MenuEntry{},
					},
				},
			},
		},
	}

	menu2 = menu.Menu{
		PluginName: "plugin_2",
		Entries: []menu.MenuEntry{
			menu.MenuEntry{
				Name:   "Entry 4",
				Route:  "/entry-4",
				Weight: 1,
				Entries: []menu.MenuEntry{},
			},
		},
	}
)

func TestMain(m *testing.M) {
	log.Println("Launching tests agilestack/utils/plugins/menu")
	menu.SetConsulAddress("127.0.0.1:8501")
	test.DoTestMain(m)
}

func consulClient(t *testing.T) *api.Client {
	config := api.DefaultConfig()
	config.Address = "localhost:8501"
	client, err := api.NewClient(config)
	require.Nil(t, err)
	return client
}


func validateEntry(t *testing.T, expectedEntry menu.MenuEntry, resultEntry menu.MenuEntry) {
	assert.Equal(t, expectedEntry.Name, resultEntry.Name)
	assert.Equal(t, expectedEntry.Route, resultEntry.Route)
	assert.Equal(t, expectedEntry.Weight, resultEntry.Weight)
	require.Equal(t, len(expectedEntry.Entries), len(resultEntry.Entries))

	for i, expectedSubEntry := range expectedEntry.Entries {
		resultSubEntry := resultEntry.Entries[i]
		validateEntry(t, expectedSubEntry, resultSubEntry)
	}
}

func validateMenu(t *testing.T, expectedMenu menu.Menu, resultMenu menu.Menu) {
	require.Equal(t, expectedMenu.PluginName, resultMenu.PluginName)
	require.Equal(t, len(expectedMenu.Entries), len(resultMenu.Entries))

	// First menu entry
	for i, inputEntry := range expectedMenu.Entries {
		validateEntry(t, inputEntry, resultMenu.Entries[i])
	}
}

func deleteAllMenus(t *testing.T) {
	_, err := consulClient(t).KV().DeleteTree("agilestack/", &api.WriteOptions{})
	require.Nil(t, err)
}
