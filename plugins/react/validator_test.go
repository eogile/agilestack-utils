package react_test

import (
	"testing"

	"github.com/eogile/agilestack-utils/plugins/react"
	"github.com/stretchr/testify/require"
)

func TestValidate_nil(t *testing.T) {
	err := react.Validate(nil)
	require.NotNil(t, err)
	require.Equal(t, "The configuration cannot not be nil", err.Error())
}

func TestValidate_PluginNameEmpty(t *testing.T) {
	config := react.PluginConfiguration{
		PluginName: "",
		Reducers:   []string{},
		Routes:     []react.Route{},
	}
	err := react.Validate(&config)
	require.NotNil(t, err)
	require.Equal(t, "Plugin name must not be blank", err.Error())
}

func TestValidate_PluginNameBlank(t *testing.T) {
	config := react.PluginConfiguration{
		PluginName: "	",
		Reducers: []string{},
		Routes:   []react.Route{},
	}
	err := react.Validate(&config)
	require.NotNil(t, err)
	require.Equal(t, "Plugin name must not be blank", err.Error())
}

func TestValidate_ReducersNil(t *testing.T) {
	config := react.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Routes:     []react.Route{},
	}
	err := react.Validate(&config)
	require.NotNil(t, err)
	require.Equal(t, "The reducers slice cannot be nil", err.Error())
}

func TestValidate_InvalidReducerName(t *testing.T) {
	config := react.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers: []string{
			"reducer1",
			"",
		},
		Routes: []react.Route{},
	}

	for _, reducerName := range []string{"", "	", "réact", "react-reducer", "react route"} {
		config.Reducers[1] = reducerName
		err := react.Validate(&config)
		require.NotNil(t, err)
		require.Equal(t, "The name of a reducer does not match the pattern \"^[a-zA-Z0-9]+$\": \""+
			reducerName+"\"", err.Error())
	}
}

func TestValidate_RoutesNil(t *testing.T) {
	config := react.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
	}
	err := react.Validate(&config)
	require.NotNil(t, err)
	require.Equal(t, "The routes slice cannot be nil", err.Error())
}

func TestValidate_InvalidRouteName(t *testing.T) {
	config := react.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
		Routes: []react.Route{
			react.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
			},
			react.Route{
				ComponentName: "Component1",
				Href:          "/route-2",
			},
		},
	}

	for _, routeName := range []string{"", "	", "réact", "react-component", "react component"} {
		config.Routes[1].ComponentName = routeName
		err := react.Validate(&config)
		require.NotNil(t, err)
		require.Equal(t, "The component name of a route does not match the pattern \"^[a-zA-Z0-9]+$\": \""+
			routeName+"\"", err.Error())
	}
}

func TestValidate_InvalidRoute(t *testing.T) {
	config := react.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
		Routes: []react.Route{
			react.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
			},
			react.Route{
				ComponentName: "Component1",
				Href:          "/route-2",
			},
		},
	}

	for _, routeLink := range []string{"", "/	", "/réact", "react-component", "/react component"} {
		config.Routes[1].Href = routeLink
		err := react.Validate(&config)
		require.NotNil(t, err)
		require.Equal(t, "The link of a route does not match the pattern \"^/[a-z0-9\\-_/]+$\": \""+
		routeLink +"\"", err.Error())
	}
}

func TestValidate(t *testing.T) {
	config := react.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
		Routes: []react.Route{
			react.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
			},
			react.Route{
				ComponentName: "Component1",
				Href:          "/route-2_1",
			},
		},
	}
	require.Nil(t, react.Validate(&config))
}