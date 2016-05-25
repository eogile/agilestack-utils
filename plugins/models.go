package plugins

import (
	"github.com/eogile/agilestack-utils/plugins/components"
	"github.com/eogile/agilestack-utils/plugins/menu"
	"github.com/eogile/agilestack-utils/plugins/registration"
)

type FullRegistration struct {
	PluginName  string
	SourcesPath string
	Menu        *menu.Menu
	Components  *components.Components
	Config      *registration.PluginConfiguration
}
