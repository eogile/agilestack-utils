package menu

import (
	"errors"
	"regexp"
	"fmt"
)

func ValidateMenu(menu *Menu) error {
	if menu == nil {
		return errors.New("Menu must not be nil")
	}

	if matched, err := regexp.MatchString("^\\s*$", menu.PluginName); matched || err != nil {
		return errors.New("Plugin name must not be blank")
	}

	if menu.Entries == nil {
		return errors.New("The menu entries slice must not be nil")
	}

	for _, menuEntry := range menu.Entries {
		if err := validateMenuEntry(menuEntry); err != nil {
			return err
		}
	}
	return nil
}

func validateMenuEntry(menuEntry MenuEntry) error {
	if matched, err := regexp.MatchString("^\\s*$", menuEntry.Name); matched || err != nil {
		return errors.New("Menu entry name must be not blank")
	}
	if matched, err := regexp.MatchString("^[a-z0-9\\-_/]+$", menuEntry.Route); !matched || err != nil {
		return fmt.Errorf("Menu entry route does not match the pattern \"^[a-z0-9\\-_/]+$\": \"%s\"",
			menuEntry.Route)
	}
	if menuEntry.Entries == nil {
		return errors.New("The menu entries slice must not be nil")
	}

	for _, menuEntry := range menuEntry.Entries {
		if err := validateMenuEntry(menuEntry); err != nil {
			return err
		}
	}
	return nil
}