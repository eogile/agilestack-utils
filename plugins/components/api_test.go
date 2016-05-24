package components_test

import (
	"encoding/json"
	"testing"

	"github.com/eogile/agilestack-utils/plugins/components"
	"github.com/stretchr/testify/require"
)

func TestStoreComponents(t *testing.T) {
	deleteAll(t)
	require.Nil(t, components.StoreComponents(&components1))

	pair, _, err := consulClient(t).KV().Get("agilestack/components/", nil)
	require.Nil(t, err)
	require.NotNil(t, pair)

	var result components.Components
	require.Nil(t, json.Unmarshal(pair.Value, &result))
	validateComponents(t, components1, result)
}

func TestStoreComponents_invalid(t *testing.T) {
	deleteAll(t)
	require.NotNil(t, components.StoreComponents(&components.Components{
		PluginName:    "",
		AppComponent:  "App",
		MainComponent: "Main",
	}))

	pair, _, err := consulClient(t).KV().Get("agilestack/components/", nil)
	require.Nil(t, err)
	require.Nil(t, pair)
}


func TestStoreComponents_update(t *testing.T) {
	deleteAll(t)
	require.Nil(t, components.StoreComponents(&components1))

	expected := components.Components{
		PluginName:    components1.PluginName,
		AppComponent:  "App2",
		MainComponent: "Main2",
	}
	require.Nil(t, components.StoreComponents(&expected))

	pair, _, err := consulClient(t).KV().Get("agilestack/components/", nil)
	require.Nil(t, err)
	require.NotNil(t, pair)

	var result components.Components
	require.Nil(t, json.Unmarshal(pair.Value, &result))
	validateComponents(t, expected, result)
}

func TestGetComponents(t *testing.T) {
	deleteAll(t)
	require.Nil(t, components.StoreComponents(&components1))

	result, err := components.GetComponents()
	require.Nil(t, err)
	require.NotNil(t, result)
	validateComponents(t, components1, *result)

}

func TestGetComponents_NoComponents(t *testing.T) {
	deleteAll(t)
	result, err := components.GetComponents()
	require.Nil(t, err)
	require.Nil(t, result)
}

func TestDeleteComponents(t *testing.T) {
	deleteAll(t)

	require.Nil(t, components.StoreComponents(&components1))
	require.Nil(t, components.DeleteComponents())

	result, err := components.GetComponents()
	require.Nil(t, result)
	require.Nil(t, err)
}