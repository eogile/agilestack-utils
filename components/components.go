package components

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats"

	"errors"

	composerPb "bitbucket.org/eogile/agilestackreloaded/plugins/composer2-api/proto"
)

type Component struct {
	Name string
	Path string
}

/*
 * Sends a message on the "composer" topic
 */
func RegisterComponents(connection *nats.EncodedConn, pluginName string, components []Component) error {
	pbComponents := []*composerPb.Component{}
	for _, component := range components {
		pbComponents = append(pbComponents, &composerPb.Component{
			Name: component.Name,
			Path: component.Path,
		})
	}

	request := &composerPb.RegisterComponentRequest{
		PluginName: pluginName,
		Components: pbComponents,
	}
	var result = composerPb.Response{}
	err := connection.Request(composerPb.RegisterComponentTopic, request, &result, 2*time.Minute)
	if err != nil {
		return err
	}
	if result.Status == composerPb.Status_KO {
		return errors.New(result.Msg)
	}

	initDeregisterComponentsHook(connection, pluginName)

	return nil
}

/*
 * Sends a message on the "composer" topic
 */
func deregisterComponents(connection *nats.EncodedConn, pluginName string) {
	request := &composerPb.RegisterComponentRequest{
		PluginName: pluginName,
	}
	var result = composerPb.Response{}
	err := connection.Request(composerPb.DeregisterComponentTopic, request, &result, 2*time.Minute)
	if err != nil {
		log.Fatalf("Error during composition : %v", err)
	}
	if result.Status == composerPb.Status_KO {
		log.Fatal("Error during composition :", result.Msg)
	}
}

func initDeregisterComponentsHook(connection *nats.EncodedConn, pluginName string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	log.Println("Initializing component deregister hook")
	go func() {
		for _ = range c {
			log.Println("Intercepting \"Interrupt\" signal")

			/*
			 * Unregistering to the proxy
			 */
			deregisterComponents(connection, pluginName)

			process, err := os.FindProcess(os.Getpid())
			if err != nil {
				log.Printf("Error while finding process : %v", err)
				os.Exit(1)
			} else {
				/*
				 * Sending a SIGINT signal.
				 *
				 * Sending a SIGTERM signal does not work.
				 */
				close(c)
				signal.Stop(c)
				process.Signal(os.Interrupt)
			}
		}
	}()
}
