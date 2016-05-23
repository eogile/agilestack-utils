package menu

var consulAddress = "consul.agilestacknet:8500"

func SetConsulAddress(consultAddress string) {
	consulAddress = consultAddress
}

// Stores the menu and returns the result error.
func StoreMenu(menu *Menu) error {
	store, err := newMenuStore()
	if err != nil {
		return err
	}
	return store.StoreMenu(menu)
}

// Lists all the existing menus.
func  ListMenus() ([]Menu, error) {
	store, err := newMenuStore()
	if err != nil {
		return nil, err
	}
	return store.ListMenus()
}

// Deletes the menu matching the given plugin name.
func DeleteMenu(pluginName string) error {
	store, err := newMenuStore()
	if err != nil {
		return err
	}
	return store.DeleteMenu(pluginName)
}