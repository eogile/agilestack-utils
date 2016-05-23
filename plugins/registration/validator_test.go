package registration_test

import (
	"testing"

	"github.com/eogile/agilestack-utils/plugins/registration"
	"github.com/stretchr/testify/require"
)

func TestValidate_nil(t *testing.T) {
	err := registration.Validate(nil)
	require.NotNil(t, err)
	require.Equal(t, "The configuration cannot not be nil", err.Error())
}

func TestValidate_PluginNameEmpty(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "",
		Reducers:   []string{},
		Routes:     []registration.Route{},
	}
	err := registration.Validate(&config)
	require.NotNil(t, err)
	require.Equal(t, "Plugin name must not be blank", err.Error())
}

func TestValidate_PluginNameBlank(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "	",
		Reducers: []string{},
		Routes:   []registration.Route{},
	}
	err := registration.Validate(&config)
	require.NotNil(t, err)
	require.Equal(t, "Plugin name must not be blank", err.Error())
}

func TestValidate_ReducersNil(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Routes:     []registration.Route{},
	}
	err := registration.Validate(&config)
	require.NotNil(t, err)
	require.Equal(t, "The reducers slice cannot be nil", err.Error())
}

func TestValidate_InvalidReducerName(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers: []string{
			"reducer1",
			"",
		},
		Routes: []registration.Route{},
	}

	for _, reducerName := range []string{"", "	", "réact", "registration-reducer", "registration route"} {
		config.Reducers[1] = reducerName
		err := registration.Validate(&config)
		require.NotNil(t, err)
		require.Equal(t, "The name of a reducer does not match the pattern \"^[a-zA-Z0-9]+$\": \""+
			reducerName+"\"", err.Error())
	}
}

func TestValidate_RoutesNil(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
	}
	err := registration.Validate(&config)
	require.NotNil(t, err)
	require.Equal(t, "The routes slice cannot be nil", err.Error())
}

func TestValidate_InvalidRouteName(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
		Routes: []registration.Route{
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
			},
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-2",
			},
		},
	}

	for _, routeName := range []string{"", "	", "réact", "registration-component", "registration component"} {
		config.Routes[1].ComponentName = routeName
		err := registration.Validate(&config)
		require.NotNil(t, err)
		require.Equal(t, "The component name of a route does not match the pattern \"^[a-zA-Z0-9]+$\": \""+
			routeName+"\"", err.Error())
	}
}

func TestValidate_InvalidRoute(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
		Routes: []registration.Route{
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
			},
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-2",
			},
		},
	}

	for _, routeLink := range []string{"", "/	", "/réact", "registration-component", "/registration component"} {
		config.Routes[1].Href = routeLink
		err := registration.Validate(&config)
		require.NotNil(t, err)
		require.Equal(t, "The link of a route does not match the pattern \"^/[a-z0-9\\-_/]+$\": \""+
		routeLink +"\"", err.Error())
	}
}

func TestValidate(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
		Routes: []registration.Route{
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
			},
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-2_1",
			},
		},
	}
	require.Nil(t, registration.Validate(&config))
}