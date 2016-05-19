package menu

type (
	/*
	A menu object represents all the menu  entries provided by a plugin.
	The application menu should be made of all the top-level menu entries
	of all deployed plugins.
	 */
	Menu struct {
		// The name of the plugin
		PluginName string `json:"pluginName"`

		// The menu entries
		Entries    []MenuEntry `json:"entries"`
	}

	MenuEntry struct {

		// The name of the menu entry.
		Name   string `json:"name"`

		// The front-end route.
		Route  string `json:"route"`

		// The menu entry weight (used to determine the rank of the menu entry).
		Weight int `json:"weight"`

		// The sub menu entries
		Entries    []MenuEntry `json:"entries"`
	}
)