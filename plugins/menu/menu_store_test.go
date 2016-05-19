package menu_test

import (
	"testing"

	"encoding/json"

	"github.com/eogile/agilestack-utils/plugins/menu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreMenu(t *testing.T) {
	deleteAllMenus(t)
	store := getStore()

	err := store.StoreMenu(&menu1)
	assert.Nil(t, err)

	pair, _, err := store.ConsulClient.KV().Get("/agilestack/menu/MyWonderfulPlugin", nil)
	require.Nil(t, err)

	foundMenu := menu.Menu{}
	err = json.Unmarshal(pair.Value, &foundMenu)
	require.Nil(t, err)
	validateMenu(t, menu1, foundMenu)
}

func TestStoreMenu_update(t *testing.T) {
	deleteAllMenus(t)
	store := getStore()

	err := store.StoreMenu(&menu1)
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
	err = store.StoreMenu(&newMenu1)
	assert.Nil(t, err)

	pair, _, err := store.ConsulClient.KV().Get("/agilestack/menu/MyWonderfulPlugin", nil)
	require.Nil(t, err)

	foundMenu := menu.Menu{}
	err = json.Unmarshal(pair.Value, &foundMenu)
	require.Nil(t, err)
	validateMenu(t, newMenu1, foundMenu)
}

func TestListMenus_NoMenu(t *testing.T) {
	deleteAllMenus(t)
	store := getStore()
	menus, err := store.ListMenus()
	assert.Nil(t, err)
	assert.Equal(t, 0, len(menus))
}

func TestListMenus(t *testing.T) {
	deleteAllMenus(t)
	store := getStore()
	err := store.StoreMenu(&menu1)
	require.Nil(t, err)
	err = store.StoreMenu(&menu2)
	require.Nil(t, err)

	menus, err := store.ListMenus()
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

func TestDeleteMenu(t *testing.T) {
	deleteAllMenus(t)
	store := getStore()
	err := store.StoreMenu(&menu1)
	require.Nil(t, err)
	err = store.StoreMenu(&menu2)
	require.Nil(t, err)

	menus, err := store.ListMenus()
	require.Nil(t, err)
	require.Equal(t, 2, len(menus))

	err = store.DeleteMenu(menu1.PluginName)
	require.Nil(t, err)

	menus, err = store.ListMenus()
	require.Nil(t, err)
	require.Equal(t, 1, len(menus))
	require.Equal(t, menu2.PluginName, menus[0].PluginName)
}
