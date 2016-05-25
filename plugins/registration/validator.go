package registration

import (
	"fmt"
	"regexp"

	"github.com/go-errors/errors"
)

const (
	routeLinkPattern    = "^[a-z0-9\\-_/:]+$"
	jsIdentifierPattern = "^[a-zA-Z0-9]+$"
)

var (
	routesSliceNil = errors.New("The routes slice cannot be nil")
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
		return routesSliceNil
	}

	if err := validateRoutes(config.Routes); err != nil {
		return err
	}

	return nil
}

func validateReducers(reducers []string) error {
	for _, reducer := range reducers {
		if matched, err := regexp.MatchString(jsIdentifierPattern, reducer); !matched || err != nil {
			return fmt.Errorf("The name of a reducer does not match the pattern \"%s\": \"%s\"",
				jsIdentifierPattern, reducer)
		}
	}
	return nil
}

func validateRoutes(routes []Route) error {
	for _, route := range routes {
		if err := validateRoute(route); err != nil {
			return err
		}
	}
	return nil
}

func validateRoute(route Route) error {
	if matched, err := regexp.MatchString(jsIdentifierPattern, route.ComponentName); !matched || err != nil {
		return fmt.Errorf("The component name of a route does not match the pattern \"%s\": \"%s\"",
			jsIdentifierPattern, route.ComponentName)
	}

	if route.Routes == nil {
		return routesSliceNil
	}

	if route.Type != "content-route" && route.Type != "full-screen-route" {
		return fmt.Errorf("Invalid route type: \"%s\"", route.Type)
	}

	if route.IsIndex {
		if route.Href != "" {
			return errors.New("An index route cannot have a path")
		}
	} else {
		/*
	         * A non-index route without sub-routes is required to have a path.
	         */
		if len(route.Routes) == 0 || route.Href != "" {
			if matched, err := regexp.MatchString(routeLinkPattern, route.Href); !matched || err != nil {
				return fmt.Errorf("The link of a route does not match the pattern \"%s\": \"%s\"",
					routeLinkPattern, route.Href)
			}
		}
	}

	for _, subRoute := range route.Routes {
		if err := validateSubRoute(subRoute); err != nil {
			return err
		}
	}
	return nil
}

func validateSubRoute(subRoute SubRoute) error {
	if matched, err := regexp.MatchString(jsIdentifierPattern, subRoute.ComponentName); !matched || err != nil {
		return fmt.Errorf("The component name of a route does not match the pattern \"%s\": \"%s\"",
			jsIdentifierPattern, subRoute.ComponentName)
	}

	/*
	 * A route without sub-routes is required to have a path.
	 */
	if len(subRoute.Routes) == 0 || subRoute.Href != "" {
		if matched, err := regexp.MatchString(routeLinkPattern, subRoute.Href); !matched || err != nil {
			return fmt.Errorf("The link of a route does not match the pattern \"%s\": \"%s\"",
				routeLinkPattern, subRoute.Href)
		}
	}

	if subRoute.Routes == nil {
		return routesSliceNil
	}
	return nil
}
