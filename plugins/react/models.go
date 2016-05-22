package react

type (
	Route struct {
		Href          string `json:"href"`
		ComponentName string `json:"componentName"`
	}

	PluginConfiguration struct {
		PluginName string   `json:"pluginName"`
		Reducers   []string `json:"reducers"`
		Routes     []Route  `json:"routes"`
	}
)
