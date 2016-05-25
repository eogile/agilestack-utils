package registration

type (
	SubRoute struct {
		Href          string     `json:"href"`
		ComponentName string     `json:"componentName"`
		Routes        []SubRoute `json:"routes"`
	}

	Route struct {
		Href          string     `json:"href"`
		ComponentName string     `json:"componentName"`
		Routes        []SubRoute `json:"routes"`
		Type          string     `json:"type"`
		IsIndex       bool       `json:"isIndex"`
	}

	PluginConfiguration struct {
		PluginName string   `json:"pluginName"`
		Reducers   []string `json:"reducers"`
		Routes     []Route  `json:"routes"`
	}
)
