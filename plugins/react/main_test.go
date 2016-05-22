package react_test

import (
	"log"
	"testing"

	"github.com/eogile/agilestack-utils/plugins/react"
	"github.com/eogile/agilestack-utils/plugins/test"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/require"
)

var (
	config1 = react.PluginConfiguration{
		PluginName: "My wonderful plugin",
		Reducers:   []string{
			"reducer1",
			"reducer2",
		},
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

	config2 = react.PluginConfiguration{
		PluginName: "Plugin 2",
		Reducers:   []string{},
		Routes: []react.Route{
			react.Route{
				ComponentName: "SomeBusinessComponent",
				Href:          "/route-10",
			},
		},
	}

	config3 = react.PluginConfiguration{
		PluginName: "Plugin 3 éé",
		Reducers:   []string{
			"reducer3",
		},
		Routes: []react.Route{},
	}
)

func TestMain(m *testing.M) {
	log.Println("Launching tests agilestack/utils/plugins/menu")
	react.ConsulAddress = "127.0.0.1:8501"
	test.DoTestMain(m)
}

func consulClient(t *testing.T) *api.Client {
	config := api.DefaultConfig()
	config.Address = "localhost:8501"
	client, err := api.NewClient(config)
	require.Nil(t, err)
	return client
}

func validateConfig(t *testing.T, expectedConfig , resultConfig *react.PluginConfiguration) {
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