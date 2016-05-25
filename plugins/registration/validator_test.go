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
				Routes:        []registration.SubRoute{},
				Type:          "content-route",
			},
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-2",
				Routes:        []registration.SubRoute{},
				Type:          "content-route",
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
				Routes:        []registration.SubRoute{},
				Type:          "content-route",
			},
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-2",
				Routes:        []registration.SubRoute{},
				Type:          "content-route",
			},
		},
	}

	for _, routeLink := range []string{"", "/	", "/réact", "registration-component", "/registration component"} {
		config.Routes[1].Href = routeLink
		err := registration.Validate(&config)
		require.NotNil(t, err)
		require.Equal(t, "The link of a route does not match the pattern \"^/[a-z0-9\\-_/:]*$\": \""+
			routeLink+"\"", err.Error())
	}
}

func TestValidate_SubRoutesNil(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
		Routes: []registration.Route{
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
				Routes:        []registration.SubRoute{},
				Type:          "content-route",
			},
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-2",
				Type:          "content-route",
			},
		},
	}
	err := registration.Validate(&config)
	require.NotNil(t, err)
	require.Equal(t, "The routes slice cannot be nil", err.Error())
}

func TestValidate_SubRoutesLinkPattern(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
		Routes: []registration.Route{
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
				Routes:        []registration.SubRoute{},
				Type:          "content-route",
			},
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-2",
				Routes: []registration.SubRoute{
					registration.SubRoute{
						Href:          "/route100",
						ComponentName: "component42",
						Routes:        []registration.SubRoute{},
					},
					registration.SubRoute{
						Href:          "/route 101",
						ComponentName: "component42",
						Routes:        []registration.SubRoute{},
					},
				},
				Type: "content-route",
			},
		},
	}
	err := registration.Validate(&config)
	require.NotNil(t, err)
	require.Equal(t, "The link of a route does not match the pattern \"^/[a-z0-9\\-_/:]*$\": \""+
		config.Routes[1].Routes[1].Href+"\"", err.Error())
}

func TestValidate_SubSubRoutesNil(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
		Routes: []registration.Route{
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
				Routes:        []registration.SubRoute{},
				Type:          "content-route",
			},
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-2",
				Routes: []registration.SubRoute{
					registration.SubRoute{
						Href:          "/route100",
						ComponentName: "component42",
						Routes:        []registration.SubRoute{},
					},
					registration.SubRoute{
						Href:          "/route101",
						ComponentName: "component42",
					},
				},
				Type: "content-route",
			},
		},
	}
	err := registration.Validate(&config)
	require.NotNil(t, err)
	require.Equal(t, "The routes slice cannot be nil", err.Error())
}

func TestValidate_SubRoutesComponentName(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
		Routes: []registration.Route{
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
				Routes:        []registration.SubRoute{},
				Type:          "content-route",
			},
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-2",
				Routes: []registration.SubRoute{
					registration.SubRoute{
						Href:          "/route100",
						ComponentName: "component42",
						Routes:        []registration.SubRoute{},
					},
					registration.SubRoute{
						Href:          "/route101",
						ComponentName: "component 42",
						Routes:        []registration.SubRoute{},
					},
				},
				Type: "content-route",
			},
		},
	}
	err := registration.Validate(&config)
	require.NotNil(t, err)
	require.Equal(t, "The component name of a route does not match the pattern \"^[a-zA-Z0-9]+$\": \""+
		config.Routes[1].Routes[1].ComponentName+"\"", err.Error())
}

func TestValidate_NoType(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
		Routes: []registration.Route{
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
				Routes:        []registration.SubRoute{},
				Type:          "content-route",
			},
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-2",
				Routes:        []registration.SubRoute{},
			},
		},
	}
	err := registration.Validate(&config)
	require.NotNil(t, err)
	require.Equal(t, "Invalid route type: \"\"", err.Error())
}

func TestValidate_InvalidType(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
		Routes: []registration.Route{
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
				Routes:        []registration.SubRoute{},
				Type:          "content-route",
			},
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-2",
				Routes:        []registration.SubRoute{},
				Type:          "content-route2",
			},
		},
	}
	err := registration.Validate(&config)
	require.NotNil(t, err)
	require.Equal(t, "Invalid route type: \"content-route2\"", err.Error())
}

func TestValidate_IndexRouteWithPath(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
		Routes: []registration.Route{
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
				Routes:        []registration.SubRoute{},
				Type:          "content-route",
				IsIndex:       false,
			},
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-2",
				Routes:        []registration.SubRoute{},
				Type:          "content-route",
				IsIndex:       true,
			},
		},
	}
	err := registration.Validate(&config)
	require.NotNil(t, err)
	require.Equal(t, "An index route cannot have a path", err.Error())
}

// A route with sub routes may not be have a path
func TestValidate_RouteWithoutLink(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
		Routes: []registration.Route{
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
				Routes:        []registration.SubRoute{},
				Type:          "content-route",
			},
			registration.Route{
				ComponentName: "Component1",
				Routes: []registration.SubRoute{
					registration.SubRoute{
						Href:          "/route100",
						ComponentName: "component42",
						Routes:        []registration.SubRoute{},
					},
					registration.SubRoute{
						Href:          "/route101",
						ComponentName: "component42",
						Routes:        []registration.SubRoute{},
					},
				},
				Type: "full-screen-route",
			},
		},
	}
	require.Nil(t, registration.Validate(&config))
}

func TestValidate(t *testing.T) {
	config := registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{},
		Routes: []registration.Route{
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
				Routes:        []registration.SubRoute{},
				Type:          "full-screen-route",
			},
			registration.Route{
				ComponentName: "Component10",
				Href:          "/",
				Routes:        []registration.SubRoute{},
				Type:          "full-screen-route",
			},
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-2_1/:param",
				Routes:        []registration.SubRoute{},
				Type:          "full-screen-route",
			},
		},
	}
	require.Nil(t, registration.Validate(&config))
}
