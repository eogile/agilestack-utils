package menu_test

import (
	"testing"
	"github.com/eogile/agilestack-utils/plugins/menu"
	"github.com/stretchr/testify/require"
)

func TestValidate_nil(t *testing.T) {
	err := menu.ValidateMenu(nil)
	require.NotNil(t, err)
	require.Equal(t, "Menu must not be nil", err.Error())
}

func TestValidate_pluginNameEmpty(t *testing.T) {
	inputMenu := &menu.Menu{
		PluginName:"",
	}
	err := menu.ValidateMenu(inputMenu)
	require.NotNil(t, err)
	require.Equal(t, "Plugin name must not be blank", err.Error())
}

func TestValidate_pluginNameBlank(t *testing.T) {
	inputMenu := &menu.Menu{
		PluginName:"	",
	}
	err := menu.ValidateMenu(inputMenu)
	require.NotNil(t, err)
	require.Equal(t, "Plugin name must not be blank", err.Error())
}

func TestValidate_MenuEntries(t *testing.T) {
	inputMenu := &menu.Menu{
		PluginName:"AgileStack rocks",
	}
	err := menu.ValidateMenu(inputMenu)
	require.NotNil(t, err)
	require.Equal(t, "The menu entries slice must not be nil", err.Error())
}

func TestValidate_MenuEntryNameEmpty(t *testing.T) {
	inputMenu := &menu.Menu{
		PluginName:"AgileStack rocks",
		Entries: []menu.MenuEntry {
			menu.MenuEntry{
				Name:"entry1",
				Route:"/route1",
				Weight:10,
				Entries: []menu.MenuEntry {},
			},
			menu.MenuEntry{
				Name:"",
				Route:"/route1",
				Weight:10,
				Entries: []menu.MenuEntry {},
			},
		},
	}
	err := menu.ValidateMenu(inputMenu)
	require.NotNil(t, err)
	require.Equal(t, "Menu entry name must be not blank", err.Error())
}

func TestValidate_MenuEntryRoutePatternSpace(t *testing.T) {
	inputMenu := &menu.Menu{
		PluginName:"AgileStack rocks",
		Entries: []menu.MenuEntry {
			menu.MenuEntry{
				Name:"entry1",
				Route:"/route1",
				Weight:10,
				Entries: []menu.MenuEntry {},
			},
			menu.MenuEntry{
				Name:"entry2",
				Route:"/route 1",
				Weight:10,
				Entries: []menu.MenuEntry {},
			},
		},
	}
	err := menu.ValidateMenu(inputMenu)
	require.NotNil(t, err)
	require.Equal(t, "Menu entry route does not match the pattern \"^[a-z0-9\\-_/]+$\": \"/route 1\"", err.Error())
}

func TestValidate_MenuEntryRoutePatternAccent(t *testing.T) {
	inputMenu := &menu.Menu{
		PluginName:"AgileStack rocks",
		Entries: []menu.MenuEntry {
			menu.MenuEntry{
				Name:"entry1",
				Route:"/route1",
				Weight:10,
				Entries: []menu.MenuEntry {},
			},
			menu.MenuEntry{
				Name:"entry2",
				Route:"/routé1",
				Weight:10,
				Entries: []menu.MenuEntry {},
			},
		},
	}
	err := menu.ValidateMenu(inputMenu)
	require.NotNil(t, err)
	require.Equal(t, "Menu entry route does not match the pattern \"^[a-z0-9\\-_/]+$\": \"/routé1\"", err.Error())
}

func TestValidate_MenuEntrySubEntriesNil(t *testing.T) {
	inputMenu := &menu.Menu{
		PluginName:"AgileStack rocks",
		Entries: []menu.MenuEntry {
			menu.MenuEntry{
				Name:"entry1",
				Route:"/route1",
				Weight:10,
				Entries: []menu.MenuEntry {},
			},
			menu.MenuEntry{
				Name:"entry2",
				Route:"/route1",
				Weight:10,
			},
		},
	}
	err := menu.ValidateMenu(inputMenu)
	require.NotNil(t, err)
	require.Equal(t, "The menu entries slice must not be nil", err.Error())
}

func TestValidate_MenuEntrySubEntryInvalid(t *testing.T) {
	inputMenu := &menu.Menu{
		PluginName:"AgileStack rocks",
		Entries: []menu.MenuEntry {
			menu.MenuEntry{
				Name:"entry1",
				Route:"/route1",
				Weight:10,
				Entries: []menu.MenuEntry {},
			},
			menu.MenuEntry{
				Name:"entry2",
				Route:"/route1",
				Weight:10,
				Entries: []menu.MenuEntry {
					menu.MenuEntry{
						Name:"entry 3",
						Route:"",
						Weight:10,
						Entries: []menu.MenuEntry {},
					},
				},
			},
		},
	}
	err := menu.ValidateMenu(inputMenu)
	require.NotNil(t, err)
	require.Equal(t, "Menu entry route does not match the pattern \"^[a-z0-9\\-_/]+$\": \"\"", err.Error())
}

func TestValidate(t *testing.T) {
	inputMenu := &menu.Menu{
		PluginName:"AgileStack rocks",
		Entries: []menu.MenuEntry {
			menu.MenuEntry{
				Name:"entry1",
				Route:"/route1",
				Weight:10,
				Entries: []menu.MenuEntry {},
			},
			menu.MenuEntry{
				Name:"entry2",
				Route:"/route1",
				Weight:10,
				Entries: []menu.MenuEntry {},
			},
		},
	}
	err := menu.ValidateMenu(inputMenu)
	require.Nil(t, err)
}