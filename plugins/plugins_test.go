package plugins

import (
	"github.com/hashicorp/consul/api"
	"log"
	"testing"
	"github.com/eogile/agilestack-utils/plugins/test"
	TestUtils "github.com/eogile/agilestack-utils/test"
)

var consulTestClient *api.Client
var consulTestStorageClient *ConsulStorageClient

func TestMain(m *testing.M) {
	log.Println("Launching tests agilestack/utils/plugins")

	consulTestClient, _ = TestUtils.NewConsulClient()
	consulTestStorageClient = &ConsulStorageClient{consulTestClient}

	test.DoTestMain(m)
}
