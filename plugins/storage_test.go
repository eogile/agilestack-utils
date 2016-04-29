package plugins

import (
	"encoding/json"
	"github.com/hashicorp/consul/api"
	"log"
	"testing"
)

//const expectedResourceJson = "{\"name\":\"accounts\",\"keys\":{\"lang\":\"accounts\",\"security\":\"rn:hydra:accounts\"},\"permissions\":[\"create\",\"get\",\"delete\",\"put:password\",\"put:data\"]}"
const expectedResourceJson = "{\"key\":\"accounts\",\"security_key\":\"rn:hydra:accounts\",\"permissions\":[\"create\",\"get\",\"delete\",\"put:password\",\"put:data\"]}"
const expectedResourceGetJson = "{\"key\":\"accountsGet\",\"security_key\":\"rn:hydra:accounts\",\"permissions\":[\"create\",\"get\",\"delete\",\"put:password\",\"put:data\"]}"

var resourceToStore = Resource{
	Key:         "accounts",
	SecurityKey: "rn:hydra:accounts",
	Permissions: []string{"create", "get", "delete", "put:password", "put:data"},
}

func TestConsulAfterInit(t *testing.T) {
	log.Println("Launching tests agilestack/utils/plugins/storage_test")

	// Get a handle to the KV API
	kv := consulTestClient.KV()

	// PUT a new KV pair
	p := &api.KVPair{Key: "foo", Value: []byte("test")}
	_, err := kv.Put(p, nil)
	if err != nil {
		t.Fatal("Cannot put key foo in consul store", err)
	}

	// Lookup the pair
	pair, _, err := kv.Get("foo", nil)
	if err != nil {
		t.Fatal("Cannot get key foo in consul store", err)
	}
	log.Printf("KV: %v", string(pair.Value))
	if pair == nil {
		t.Fatal("Got nil kv pair for index foo, expected one")
	}
	value := string(pair.Value)
	if "test" != value {
		t.Error("Expected test from get in consul store for foo, got ", value)
	}

	p.Value = []byte("modified")
	_, err = kv.Put(p, nil)
	if err != nil {
		t.Fatal("Cannot put key foo in consul store", err)
	}
	pair, _, err = kv.Get("foo", nil)
	if err != nil {
		t.Fatal("Cannot get key foo in consul store", err)
	}
	log.Printf("KV: %v", string(pair.Value))
	if pair == nil {
		t.Fatal("Got nil kv pair for index foo, expected one")
	}
	value = string(pair.Value)
	if "modified" != value {
		t.Error("Expected modified from get in consul store for foo, got ", value)
	}

	//remove
	_, errDel := kv.Delete("foo", nil)
	if errDel != nil {
		t.Fatal("Cannot delete foo from Consul Store")
	}

}

func TestStoreResource(t *testing.T) {

	err := consulTestStorageClient.StoreResource(resourceToStore)
	if err != nil {
		t.Fatal("CannotStore test resource", err)
	}

	// Get a handle to the KV API
	kv := consulTestClient.KV()

	// Lookup the pair
	pair, _, err := kv.Get("/agilestack/security/resources/accounts", nil)
	if err != nil {
		t.Fatal("Cannot get key /agilestack/security/resources/accounts in consul store", err)
	}
	if pair == nil {
		t.Fatal("Got nil kv pair for index /agilestack/security/resources/accounts, expected one")
	}
	log.Printf("KV: %v", string(pair.Value))
	value := string(pair.Value)
	if expectedResourceJson != value {
		t.Error("Expected test from get in consul store for foo, got ", value)
	}
}

func TestGetResource(t *testing.T) {

	resourceToStore.Key = "accountsGet"
	err := consulTestStorageClient.StoreResource(resourceToStore)
	if err != nil {
		t.Fatal("CannotStore test resource", err)
	}

	resourceFound, errGet := consulTestStorageClient.GetResource(resourceToStore.Key)
	if errGet != nil {
		t.Fatal("Unable to retreive resource ", resourceToStore.Key)
	}
	if resourceFound == nil {
		t.Error("Did not found any resource under ", resourceToStore.Key)
	} else {
		resourceFoundJson, errJson := json.Marshal(resourceFound)
		if errJson != nil {
			t.Fatal("unable to Marshal resource found", errJson)
		}
		if expectedResourceGetJson != string(resourceFoundJson) {
			t.Errorf("Expected %s, got %s", expectedResourceGetJson, string(resourceFoundJson))
		}
	}
	// Get a handle to the KV API
	//kv := consulTestClient.KV()
	//pairs, _, _ := kv.List("", nil)
	//for _, pair := range pairs {
	//	log.Printf("key:%s, value:%s", pair.Key, pair.Value)
	//}

}

func TestListResources(t *testing.T) {

	index1 := "resource1"
	index2 := "resource2"

	resources, errList := consulTestStorageClient.ListResources()
	if errList != nil {
		t.Fatal("Unable to list resources")
	}
	resourceToStore.Key = index1
	initialLength := len(resources)
	err := consulTestStorageClient.StoreResource(resourceToStore)
	if err != nil {
		t.Fatal("Got error when trying to store resource ", err)
	}
	resources, errList = consulTestStorageClient.ListResources()
	if errList != nil {
		t.Fatal("Unable to list resources")
	}
	if len(resources) != initialLength+1 {
		t.Errorf("Found %v resources, expected 1", len(resources))
	}
	resourceToStore.Key = index2
	err2 := consulTestStorageClient.StoreResource(resourceToStore)
	if err2 != nil {
		t.Fatal("Got error when trying to store resource ", err2)
	}
	resources, errList = consulTestStorageClient.ListResources()
	if errList != nil {
		t.Fatal("Unable to list resources")
	}
	if len(resources) != initialLength+2 {
		t.Errorf("Found %v resources, expected 2", len(resources))
	}
	//2nd store should not have any effect
	err = consulTestStorageClient.StoreResource(resourceToStore)
	if err != nil {
		t.Fatal("Got error when trying to store resource ", err)
	}
	resources, errList = consulTestStorageClient.ListResources()
	if errList != nil {
		t.Fatal("Unable to list resources")
	}
	if len(resources) != initialLength+2 {
		t.Errorf("Found %v resources, expected 2", len(resources))
	}

	err = consulTestStorageClient.DeleteResource(index1)
	if err != nil {
		t.Fatal("Got error when trying to delete resource ", err)
	}
	resources, errList = consulTestStorageClient.ListResources()
	if errList != nil {
		t.Fatal("Unable to list resources")
	}
	if len(resources) != initialLength+1 {
		t.Errorf("Found %v resources, expected 1", len(resources))
	}

	//kv := consulTestClient.KV()
	//pairs, _, _ := kv.List(resourcesPrefix, nil)
	//for _, pair := range pairs {
	//	log.Printf("in List --- key:%s, value:%s", pair.Key, pair.Value)
	//}
	//
	//for _, res := range resources {
	//	log.Printf("in List --- resource %v, permissions : %v", res.Keys.Lang, res.Permissions)
	//}
}

func TestDeleteResource(t *testing.T) {

	indexDel := "toBeDeleted"
	resourceToStore.Key = indexDel
	err := consulTestStorageClient.StoreResource(resourceToStore)
	if err != nil {
		t.Fatal("Got error when trying to store resource ", err)
	}
	err = consulTestStorageClient.DeleteResource(indexDel)
	if err != nil {
		t.Fatal("Got error when trying to delete resource ", err)
	}
	resourceFound, errGet := consulTestStorageClient.GetResource(indexDel)
	if errGet != nil {
		t.Fatal("Unable to retreive resource ", indexDel)
	}
	if resourceFound != nil {
		res, _ := json.Marshal(resourceFound)
		t.Error("Expected resource to be deleted but found", res)
	}
	//kv := consulTestClient.KV()
	//pairs, _, _ := kv.List("", nil)
	//for _, pair := range pairs {
	//	log.Printf("in Del --- key:%s, value:%s", pair.Key, pair.Value)
	//}
}
