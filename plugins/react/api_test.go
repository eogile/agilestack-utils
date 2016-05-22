package react_test

import (
	"testing"
	"github.com/eogile/agilestack-utils/plugins/react"
	"encoding/json"
	"github.com/stretchr/testify/require"
)

func TestStoreRoutesAndReducers_PluginNameWithSpaces(t *testing.T) {
	deleteAll(t)

	err := react.StoreRoutesAndReducers(&config1)
	require.Nil(t, err)

	pair, _, err := consulClient(t).KV().Get("/agilestack/react/My wonderful plugin", nil)
	require.Nil(t, err)

	foundConfig := react.PluginConfiguration{}
	err = json.Unmarshal(pair.Value, &foundConfig)
	require.Nil(t, err)
	validateConfig(t, &config1, &foundConfig)
}

func TestStoreRoutesAndReducers_Update(t *testing.T) {
	deleteAll(t)

	err := react.StoreRoutesAndReducers(&config1)
	require.Nil(t, err)

	// The new version of config1
	newConfig1 := react.PluginConfiguration{
		PluginName: config1.PluginName,
		Reducers:   []string{
			"reducer10",
		},
		Routes: []react.Route{
			react.Route{
				ComponentName: "Component1",
				Href:          "/route-10",
			},
		},
	}
	err = react.StoreRoutesAndReducers(&newConfig1)
	require.Nil(t, err)

	pair, _, err := consulClient(t).KV().Get("/agilestack/react/My wonderful plugin", nil)
	require.Nil(t, err)

	foundConfig := react.PluginConfiguration{}
	err = json.Unmarshal(pair.Value, &foundConfig)
	require.Nil(t, err)
	validateConfig(t, &newConfig1, &foundConfig)
}

func TestStoreRoutesAndReducers_Invalid(t *testing.T) {
	deleteAll(t)
	err := react.StoreRoutesAndReducers(&react.PluginConfiguration{
		PluginName:"SomePlugin",
		Reducers:[]string{},
	})
	require.NotNil(t, err)
	require.Equal(t, "The routes slice cannot be nil", err.Error())

	pair, _, err := consulClient(t).KV().Get("/agilestack/menu/SomePlugin", nil)
	require.Nil(t, err)
	require.Nil(t, pair)
}

func TestStoreRoutesAndReducers_NameWithAccent(t *testing.T) {
	deleteAll(t)

	err := react.StoreRoutesAndReducers(&config3)
	require.Nil(t, err)

	pair, _, err := consulClient(t).KV().Get("/agilestack/react/Plugin 3 éé", nil)
	require.Nil(t, err)

	foundConfig := react.PluginConfiguration{}
	err = json.Unmarshal(pair.Value, &foundConfig)
	require.Nil(t, err)
	validateConfig(t, &config3, &foundConfig)
}

func TestStoreRoutesAndReducers_DoNotReplaceOtherConfigurations(t *testing.T) {
	deleteAll(t)

	require.Nil(t, react.StoreRoutesAndReducers(&config1))
	require.Nil(t, react.StoreRoutesAndReducers(&config3))

	/*
	 * Config1 still exists
	 */
	pair1, _, err := consulClient(t).KV().Get("/agilestack/react/My wonderful plugin", nil)
	require.Nil(t, err)

	foundConfig1 := react.PluginConfiguration{}
	err = json.Unmarshal(pair1.Value, &foundConfig1)
	require.Nil(t, err)
	validateConfig(t, &config1, &foundConfig1)


	/*
	 * Config 3 also exists
	 */
	pair3, _, err := consulClient(t).KV().Get("/agilestack/react/Plugin 3 éé", nil)
	require.Nil(t, err)

	foundConfig3 := react.PluginConfiguration{}
	err = json.Unmarshal(pair3.Value, &foundConfig3)
	require.Nil(t, err)
	validateConfig(t, &config3, &foundConfig3)
}

func TestListRoutesAndReducers_Empty(t *testing.T) {
	deleteAll(t)

	configurations, err := react.ListRoutesAndReducers()
	require.Nil(t, err)
	require.NotNil(t, configurations)
	require.Equal(t, 0, len(configurations))
}

func TestListRoutesAndReducers(t *testing.T) {
	deleteAll(t)

	require.Nil(t, react.StoreRoutesAndReducers(&config1))
	require.Nil(t, react.StoreRoutesAndReducers(&config2))

	configurations, err := react.ListRoutesAndReducers()
	require.Nil(t, err)
	require.NotNil(t, configurations)
	require.Equal(t, 2, len(configurations))

	if configurations[0].PluginName == config1.PluginName {
		validateConfig(t, &config1, &configurations[0])
		validateConfig(t, &config2, &configurations[1])
	} else {
		validateConfig(t, &config2, &configurations[0])
		validateConfig(t, &config1, &configurations[1])
	}
}

func TestDeleteRoutesAndReducers(t *testing.T) {
	deleteAll(t)

	require.Nil(t, react.StoreRoutesAndReducers(&config1))
	require.Nil(t, react.StoreRoutesAndReducers(&config3))
	require.Nil(t, react.DeleteRoutesAndReducers(config3.PluginName))

	configurations, err := react.ListRoutesAndReducers()
	require.Nil(t, err)
	require.NotNil(t, configurations)

	// The only existing configuration is config1
	require.Equal(t, 1, len(configurations))
	validateConfig(t, &config1, &configurations[0])
}