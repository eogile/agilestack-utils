package components

import (
	"errors"
	"regexp"
	"fmt"
)

const jsIdentifierPattern = "^[a-zA-Z0-9]+$"

func Validate(components *Components) error {
	if components == nil {
		return errors.New("Components cannot not be nil")
	}

	if matched, err := regexp.MatchString("^\\s*$", components.PluginName); matched || err != nil {
		return errors.New("Plugin name must not be blank")
	}

	if matched, err := regexp.MatchString(jsIdentifierPattern, components.AppComponent); !matched || err != nil {
		return fmt.Errorf("The App component name does not match the pattern \"%s\": \"%s\"",
			jsIdentifierPattern, components.AppComponent)
	}

	if matched, err := regexp.MatchString(jsIdentifierPattern, components.MainComponent); !matched || err != nil {
		return fmt.Errorf("The Main component name does not match the pattern \"%s\": \"%s\"",
			jsIdentifierPattern, components.MainComponent)
	}
	return nil
}
