//This package contains tools to help initi
package plugins

import (
	"fmt"
	"log"
	"net/http"

	"bytes"
	"github.com/nats-io/nats"
	"github.com/nats-io/nats/encoders/protobuf"
	"io/ioutil"
	"os"
	"github.com/eogile/agilestack-utils/plugins/registration"
	"github.com/eogile/agilestack-utils/plugins/menu"
	"github.com/eogile/agilestack-utils/plugins/components"
	"github.com/eogile/agilestack-utils/files"
)

/*
 * Establishes a connection to the NATS server.
 */
func EstablishConnection(natsServerURL string) (*nats.EncodedConn, error) {
	log.Println("Establishing connection to the NATS server :", natsServerURL)
	connection, err := nats.Connect(natsServerURL)
	if err != nil {
		return nil, err
	}

	return nats.NewEncodedConn(connection, protobuf.PROTOBUF_ENCODER)
}

func HandleHttpStatusUrl(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(writer, "Method %s not supported", request.Method)
	} else {
		fmt.Fprintln(writer, "OK")
	}
}

/*
 * Change the base Url in index.html to make it work behind a proxy if needed
 */
func ChangeBaseUrl(rootDir string) {
	baseUrl := os.Getenv("baseUrl")
	log.Println("baseUrl = ", baseUrl)
	if baseUrl != "" {

		indexPath := rootDir + "index.html"
		oldIndexContent, err := ioutil.ReadFile(indexPath)
		if err != nil {
			log.Fatalf("unable to find index.html :%v", err)
		}

		newIndexContent := bytes.Replace(oldIndexContent, []byte("window.baseUrl=\"/\""), []byte("window.baseUrl=\"/"+baseUrl+"\""), -1)
		err = ioutil.WriteFile(indexPath, newIndexContent, 0644)
	}

}

func Register(config FullRegistration) error {
	destination := "/shared/root-app-builder/web_modules/" + config.PluginName
	if err:= files.CopyDir(config.SourcesPath, destination); err!=nil {
		return err
	}
	if err := registration.StoreRoutesAndReducers(config.Config); err!= nil {
		return err
	}

	if err := menu.StoreMenu(config.Menu); err != nil {
		return err
	}

	if config.Components != nil {
		if err:= components.StoreComponents(config.Components); err != nil {
			return err
		}
	}

	return registration.LaunchApplicationBuild()
}

