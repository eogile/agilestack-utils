package registration

import (
	"github.com/go-errors/errors"
	"regexp"
	"fmt"
)

func Validate(config *PluginConfiguration) error {
	if config == nil {
		return errors.New("The configuration cannot not be nil")
	}

	if matched, err := regexp.MatchString("^\\s*$", config.PluginName); matched || err != nil {
		return errors.New("Plugin name must not be blank")
	}

	if config.Reducers == nil {
		return errors.New("The reducers slice cannot be nil")
	}

	if err := validateReducers(config.Reducers); err != nil {
		return err
	}

	if config.Routes == nil {
		return errors.New("The routes slice cannot be nil")
	}

	if err := validateRoutes(config.Routes); err != nil {
		return err
	}

	return nil
}

func validateReducers(reducers []string) error {
	for _, reducer := range reducers {
		if matched, err := regexp.MatchString("^[a-zA-Z0-9]+$", reducer); !matched || err != nil {
			return fmt.Errorf("The name of a reducer does not match the pattern \"^[a-zA-Z0-9]+$\": \"%s\"",
				reducer)
		}
	}
	return nil
}

func validateRoutes(routes []Route) error {
	for _, route := range routes {
		if matched, err := regexp.MatchString("^[a-zA-Z0-9]+$", route.ComponentName); !matched || err != nil {
			return fmt.Errorf("The component name of a route does not match the pattern \"^[a-zA-Z0-9]+$\": \"%s\"",
				route.ComponentName)
		}

		if matched, err := regexp.MatchString("^/[a-z0-9\\-_/]+$", route.Href); !matched || err != nil {
			return fmt.Errorf("The link of a route does not match the pattern \"^/[a-z0-9\\-_/]+$\": \"%s\"",
				route.Href)
		}
	}
	return nil
}