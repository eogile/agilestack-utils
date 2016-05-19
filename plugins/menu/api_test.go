package menu_test

import (
	"testing"
	"encoding/json"

	"github.com/eogile/agilestack-utils/plugins/menu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreMenu_api(t *testing.T) {
	deleteAllMenus(t)
	store := getStore()

	err := menu.StoreMenu(&menu1)
	assert.Nil(t, err)

	pair, _, err := store.ConsulClient.KV().Get("/agilestack/menu/MyWonderfulPlugin", nil)
	require.Nil(t, err)

	foundMenu := menu.Menu{}
	err = json.Unmarshal(pair.Value, &foundMenu)
	require.Nil(t, err)
	validateMenu(t, menu1, foundMenu)
}

func TestStoreMenu_update_api(t *testing.T) {
	deleteAllMenus(t)
	store := getStore()

	err := menu.StoreMenu(&menu1)
	assert.Nil(t, err)

	newMenu1 := menu.Menu{
		PluginName: menu1.PluginName,
		Entries: []menu.MenuEntry{
			menu.MenuEntry{
				Name:   "new menu entry",
				Route:  "/new-route",
				Weight: 77,
				Entries: []menu.MenuEntry{},
			},
		},
	}
	err = menu.StoreMenu(&newMenu1)
	assert.Nil(t, err)

	pair, _, err := store.ConsulClient.KV().Get("/agilestack/menu/MyWonderfulPlugin", nil)
	require.Nil(t, err)

	foundMenu := menu.Menu{}
	err = json.Unmarshal(pair.Value, &foundMenu)
	require.Nil(t, err)
	validateMenu(t, newMenu1, foundMenu)
}

func TestStoreMenu_api_Invalid(t *testing.T) {
	deleteAllMenus(t)
	store := getStore()

	err := menu.StoreMenu(&menu.Menu{
		PluginName:"AgileStack rocks",
	})
	assert.NotNil(t, err)
	require.Equal(t, "The menu entries slice must not be nil", err.Error())

	pair, _, err := store.ConsulClient.KV().Get("/agilestack/menu/MyWonderfulPlugin", nil)
	require.Nil(t, err)
	require.Nil(t, pair)
}

func TestStoreMenu_api_NameWithSpace(t *testing.T) {
	deleteAllMenus(t)
	store := getStore()
	inputMenu := &menu.Menu{
		PluginName:"AgileStack rocks",
		Entries: []menu.MenuEntry{
			menu.MenuEntry{
				Name:"entry 1",
				Route:"/route1",
				Weight:10,
				Entries: []menu.MenuEntry{},
			},
			menu.MenuEntry{
				Name:"entry2",
				Route:"/route1",
				Weight:10,
				Entries: []menu.MenuEntry{},
			},
		},
	}

	err := menu.StoreMenu(inputMenu)
	assert.Nil(t, err)

	menus, err := menu.ListMenus()
	require.Nil(t, err)
	       require.Equal(t, 1, len(menus))
	validateMenu(t, *inputMenu, menus[0])

	// manually
	pair, _, err := store.ConsulClient.KV().Get("/agilestack/menu/AgileStack rocks", nil)
	require.Nil(t, err)
	require.NotNil(t, pair)

	foundMenu := menu.Menu{}
	err = json.Unmarshal(pair.Value, &foundMenu)
	require.Nil(t, err)
	validateMenu(t, *inputMenu, foundMenu)
}

func TestStoreMenu_api_NameWithAccent(t *testing.T) {
	deleteAllMenus(t)
	store := getStore()
	inputMenu := &menu.Menu{
		PluginName:"AgileStack éé",
		Entries: []menu.MenuEntry{
			menu.MenuEntry{
				Name:"entry 1",
				Route:"/route1",
				Weight:10,
				Entries: []menu.MenuEntry{},
			},
			menu.MenuEntry{
				Name:"entry2",
				Route:"/route1",
				Weight:10,
				Entries: []menu.MenuEntry{},
			},
		},
	}

	err := menu.StoreMenu(inputMenu)
	assert.Nil(t, err)

	menus, err := menu.ListMenus()
	require.Nil(t, err)
	require.Equal(t, 1, len(menus))
	validateMenu(t, *inputMenu, menus[0])

	// manually
	pair, _, err := store.ConsulClient.KV().Get("/agilestack/menu/AgileStack éé", nil)
	require.Nil(t, err)
	require.NotNil(t, pair)

	foundMenu := menu.Menu{}
	err = json.Unmarshal(pair.Value, &foundMenu)
	require.Nil(t, err)
	validateMenu(t, *inputMenu, foundMenu)
}

func TestListMenus_NoMenu_api(t *testing.T) {
	deleteAllMenus(t)
	menus, err := menu.ListMenus()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(menus))
}

func TestListMenus_api(t *testing.T) {
	deleteAllMenus(t)
	err := menu.StoreMenu(&menu1)
	require.Nil(t, err)
	err = menu.StoreMenu(&menu2)
	require.Nil(t, err)

	menus, err := menu.ListMenus()
	require.Nil(t, err)

	require.Equal(t, 2, len(menus))

	if menus[0].PluginName == menu1.PluginName {
		validateMenu(t, menu1, menus[0])
		validateMenu(t, menu2, menus[1])
	} else {
		validateMenu(t, menu2, menus[0])
		validateMenu(t, menu1, menus[1])
	}
}

func TestDeleteMenu_api(t *testing.T) {
	deleteAllMenus(t)
	err := menu.StoreMenu(&menu1)
	require.Nil(t, err)
	err = menu.StoreMenu(&menu2)
	require.Nil(t, err)

	menus, err := menu.ListMenus()
	require.Nil(t, err)
	require.Equal(t, 2, len(menus))

	err = menu.DeleteMenu(menu1.PluginName)
	require.Nil(t, err)

	menus, err = menu.ListMenus()
	require.Nil(t, err)
	require.Equal(t, 1, len(menus))
	require.Equal(t, menu2.PluginName, menus[0].PluginName)
}
