package components_test

import (
	"testing"

	"github.com/eogile/agilestack-utils/plugins/components"
	"github.com/stretchr/testify/require"
)

func TestValidate_Nil(t *testing.T) {
	err := components.Validate(nil)
	require.NotNil(t, err)
	require.Equal(t, "Components cannot not be nil", err.Error())
}

func TestValidate_InvalidPluginName(t *testing.T) {
	input := &components.Components{
		AppComponent:  "App",
		MainComponent: "Main",
	}

	for _, pluginName := range []string{"", "	"} {
		input.PluginName = pluginName
		err := components.Validate(input)
		require.NotNil(t, err)
		require.Equal(t, "Plugin name must not be blank", err.Error())
	}
}

func TestValidate_InvalidAppComponentName(t *testing.T) {
	input := &components.Components{
		PluginName:    "agilestack-plugin",
		AppComponent:  "",
		MainComponent: "Main",
	}

	for _, appComponentName := range []string{"", "  ", "ré", "sd dsqd"} {
		input.AppComponent = appComponentName
		err := components.Validate(input)
		require.NotNil(t, err)
		require.Equal(t, "The App component name does not match the pattern \"^[a-zA-Z0-9]+$\": \""+
			appComponentName+
			"\"", err.Error())
	}
}

func TestValidate_InvalidMainComponentName(t *testing.T) {
	input := &components.Components{
		PluginName:    "agilestack-plugin",
		AppComponent:  "App",
		MainComponent: "",
	}

	for _, mainComponentName := range []string{"", "  ", "ré", "sd dsqd"} {
		input.MainComponent = mainComponentName
		err := components.Validate(input)
		require.NotNil(t, err)
		require.Equal(t, "The Main component name does not match the pattern \"^[a-zA-Z0-9]+$\": \""+
		mainComponentName +
		"\"", err.Error())
	}
}