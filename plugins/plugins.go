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
