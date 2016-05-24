package components_test

import (
	"log"
	"testing"

	"github.com/eogile/agilestack-utils/plugins/components"
	"github.com/eogile/agilestack-utils/plugins/test"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/require"
)

var components1 = components.Components{
	PluginName:    "my-plugin",
	AppComponent:  "App",
	MainComponent: "Main",
}

func TestMain(m *testing.M) {
	log.Println("Launching tests agilestack/utils/plugins/menu")
	components.SetConsulAddress("127.0.0.1:8501")
	test.DoTestMain(m)
}

func consulClient(t *testing.T) *api.Client {
	config := api.DefaultConfig()
	config.Address = "localhost:8501"
	client, err := api.NewClient(config)
	require.Nil(t, err)
	return client
}

func deleteAll(t *testing.T) {
	_, err := consulClient(t).KV().DeleteTree("agilestack/", &api.WriteOptions{})
	require.Nil(t, err)
}

func validateComponents(t *testing.T, expected, result components.Components) {
	require.Equal(t, expected.PluginName, result.PluginName)
	require.Equal(t, expected.AppComponent, result.AppComponent)
	require.Equal(t, expected.MainComponent, result.MainComponent)
}
