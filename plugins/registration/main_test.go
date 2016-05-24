package registration_test

import (
	"log"
	"testing"

	"github.com/eogile/agilestack-utils/plugins/registration"
	"github.com/eogile/agilestack-utils/plugins/test"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/require"
)

var (
	config1 = registration.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{
			"reducer1",
			"reducer2",
		},
		Routes: []registration.Route{
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-1",
				Routes: []registration.SubRoute{},
				Type:"content-route",
			},
			registration.Route{
				ComponentName: "Component1",
				Href:          "/route-2_1",
				Routes: []registration.SubRoute{},
				Type:"content-route",
			},
		},
	}

	config2 = registration.PluginConfiguration{
		PluginName: "Plugin 2",
		Reducers:   []string{},
		Routes: []registration.Route{
			registration.Route{
				ComponentName: "SomeBusinessComponent",
				Href:          "/route-10",
				Routes: []registration.SubRoute{},
				Type:"content-route",
			},
		},
	}

	config3 = registration.PluginConfiguration{
		PluginName: "Plugin 3 éé",
		Reducers:   []string{
			"reducer3",
		},
		Routes: []registration.Route{},
	}
)

func TestMain(m *testing.M) {
	log.Println("Launching tests agilestack/utils/plugins/menu")
	registration.SetConsulAddress("127.0.0.1:8501")
	test.DoTestMain(m)
}

func consulClient(t *testing.T) *api.Client {
	config := api.DefaultConfig()
	config.Address = "localhost:8501"
	client, err := api.NewClient(config)
	require.Nil(t, err)
	return client
}

func validateConfig(t *testing.T, expectedConfig , resultConfig *registration.PluginConfiguration) {
	require.Equal(t, expectedConfig.PluginName, resultConfig.PluginName)
	require.Equal(t, expectedConfig.Reducers, resultConfig.Reducers)
	require.Equal(t, len(expectedConfig.Routes), len(resultConfig.Routes))

	for i, expectedRoute := range expectedConfig.Routes {
		resultRoute := resultConfig.Routes[i]
		require.Equal(t, expectedRoute.ComponentName, resultRoute.ComponentName)
		require.Equal(t, expectedRoute.Href, resultRoute.Href)
	}
}

func deleteAll(t *testing.T) {
	_, err := consulClient(t).KV().DeleteTree("agilestack/", &api.WriteOptions{})
	require.Nil(t, err)
}