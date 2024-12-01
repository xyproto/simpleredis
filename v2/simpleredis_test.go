package simpleredis

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/xyproto/pinterface"
)

var pool *ConnectionPool

func TestLocalConnection(t *testing.T) {
	if err := TestConnection(); err != nil {
		if strings.HasSuffix(err.Error(), "i/o timeout") {
			log.Println("Try the 'latency doctor' command in the redis-cli if I/O timeouts happens often.")
		}
		t.Errorf(err.Error())
	}
}

func TestRemoteConnection(t *testing.T) {
	if err := TestConnectionHost("foobared@ :6379"); err != nil {
		t.Errorf(err.Error())
	}
}

func TestConnectionPool(t *testing.T) {
	pool = NewConnectionPool()
}

func TestConnectionPoolHost(t *testing.T) {
	pool = NewConnectionPoolHost("localhost:6379")
}

// Tests with password "foobared" if the previous connection test
// did not result in a connection that responds to PING.
func TestConnectionPoolHostPassword(t *testing.T) {
	if pool.Ping() != nil {
		// Try connecting with the default password
		pool = NewConnectionPoolHost("foobared@localhost:6379")
	}
}

func TestList(t *testing.T) {
	const (
		listname = "abc123_test_test_test_123abc"
		testdata = "123abc"
	)
	list := NewList(pool, listname)

	// Check that the list qualifies for the IList interface
	var _ pinterface.IList = list

	list.SelectDatabase(1)
	if err := list.Add(testdata); err != nil {
		t.Errorf("Error, could not add item to list! %s", err.Error())
	}
	items, err := list.All()
	if err != nil {
		t.Errorf("Error, could not retrieve list! %s", err.Error())
	}
	if len(items) != 1 {
		t.Errorf("Error, wrong list length! %v", len(items))
	}
	if (len(items) > 0) && (items[0] != testdata) {
		t.Errorf("Error, wrong list contents! %v", items)
	}
	err = list.Remove()
	if err != nil {
		t.Errorf("Error, could not remove list! %s", err.Error())
	}
}

func TestRemove(t *testing.T) {
	const (
		kvname    = "abc123_test_test_test_123abc"
		testkey   = "sdsdf234234"
		testvalue = "asdfasdf1234"
	)
	kv := NewKeyValue(pool, kvname)

	// TODO: Also do this check for ISet and IHashMap
	// Check that the key/value qualifies for the IKeyValue interface
	var _ pinterface.IKeyValue = kv

	kv.SelectDatabase(1)
	if err := kv.Set(testkey, testvalue); err != nil {
		t.Errorf("Error, could not set key and value! %s", err.Error())
	}
	if val, err := kv.Get(testkey); err != nil {
		t.Errorf("Error, could not get key! %s", err.Error())
	} else if val != testvalue {
		t.Errorf("Error, wrong value! %s != %s", val, testvalue)
	}
	kv.Remove()
	if _, err := kv.Get(testkey); err == nil {
		t.Errorf("Error, could get key! %s", err.Error())
	}
}

func TestInc(t *testing.T) {
	const (
		kvname     = "kv_234_test_test_test"
		testkey    = "key_234_test_test_test"
		testvalue0 = "9"
		testvalue1 = "10"
		testvalue2 = "1"
	)
	kv := NewKeyValue(pool, kvname)
	kv.SelectDatabase(1)
	if err := kv.Set(testkey, testvalue0); err != nil {
		t.Errorf("Error, could not set key and value! %s", err.Error())
	}
	if val, err := kv.Get(testkey); err != nil {
		t.Errorf("Error, could not get key! %s", err.Error())
	} else if val != testvalue0 {
		t.Errorf("Error, wrong value! %s != %s", val, testvalue0)
	}
	incval, err := kv.Inc(testkey)
	if err != nil {
		t.Errorf("Error, could not INCR key! %s", err.Error())
	}
	if val, err := kv.Get(testkey); err != nil {
		t.Errorf("Error, could not get key! %s", err.Error())
	} else if val != testvalue1 {
		t.Errorf("Error, wrong value! %s != %s", val, testvalue1)
	} else if incval != testvalue1 {
		t.Errorf("Error, wrong inc value! %s != %s", incval, testvalue1)
	}
	kv.Remove()
	if _, err := kv.Get(testkey); err == nil {
		t.Errorf("Error, could get key! %s", err.Error())
	}
	// Creates "0" and increases the value with 1
	kv.Inc(testkey)
	if val, err := kv.Get(testkey); err != nil {
		t.Errorf("Error, could not get key! %s", err.Error())
	} else if val != testvalue2 {
		t.Errorf("Error, wrong value! %s != %s", val, testvalue2)
	}
	kv.Remove()
	if _, err := kv.Get(testkey); err == nil {
		t.Errorf("Error, could get key! %s", err.Error())
	}
}

func TestTwoFields(t *testing.T) {
	test, test23, ok := twoFields("test1@test2@test3", "@")
	if ok && ((test != "test1") || (test23 != "test2@test3")) {
		t.Error("Error in twoFields functions")
	}
}

func TestICreator(t *testing.T) {
	// Check if the struct comforms to ICreator
	var _ pinterface.ICreator = NewCreator(pool, 1)
}

func TestKeyValue(t *testing.T) {
	const (
		kvname  = "kv_abc123_test_test_test_123abc"
		testkey = "token"
		testval = "123abc"
		fakekey = "hurdygurdy32"
	)
	kv := NewKeyValue(pool, kvname)

	// Check that the list qualifies for the IList interface
	var _ pinterface.IKeyValue = kv

	kv.SelectDatabase(1)

	if err := kv.Set(testkey, testval); err != nil {
		t.Errorf("Error, could not set key and value! %s", err.Error())
	}
	retval, err := kv.Get(testkey)
	if err != nil {
		t.Errorf("Error, could not get value! %s", err.Error())
	} else if retval != testval {
		t.Errorf("Error, got the wrong return value! %s", retval)
	}
	if err := kv.Del(testkey); err != nil {
		t.Errorf("Error, could not remove key! %s", err.Error())
	}
	_, err = kv.Get(testkey)
	if err == nil {
		t.Errorf("Error, key should be gone #1! %s", err.Error())
	}
	_, err = kv.Get(fakekey)
	if err == nil {
		t.Errorf("Error, key should be gone #2! %s", err.Error())
	}
	err = kv.Remove()
	if err != nil {
		t.Errorf("Error, could not remove KeyValue! %s", err.Error())
	}
}

func TestExpire(t *testing.T) {
	const (
		kvname  = "kv_abc123_test_test_test_123abc_exp"
		testkey = "token"
		testval = "123abc"
	)
	kv := NewKeyValue(pool, kvname)

	// Check that the list qualifies for the IList interface
	var _ pinterface.IKeyValue = kv

	kv.SelectDatabase(1)

	if err := kv.SetExpire(testkey, testval, time.Second*1); err != nil {
		t.Errorf("Error, could not set key and value! %s", err.Error())
	}
	retval, err := kv.Get(testkey)
	if err != nil {
		t.Errorf("Error, could not get value! %s", err.Error())
	} else if retval != testval {
		t.Errorf("Error, got the wrong return value! %s", retval)
	}
	ttl, err := kv.TimeToLive(testkey)
	if err != nil {
		t.Errorf("Error, retrieving time to live: %s", err.Error())
	} else if ttl.String() != "1s" {
		t.Errorf("Error, there should only be 1 second left, but there is: %s!", ttl.String())
	}
	// Wait a bit extra, testing on external hosts may take some time
	time.Sleep(3 * time.Second)

	_, err2 := kv.Get(testkey)
	if err2 == nil {
		t.Errorf("Error, key should be gone! %s", testkey)
	}
	err = kv.Remove()
	if err != nil {
		t.Errorf("Error, could not remove KeyValue! %s", err.Error())
	}
}

func TestExpireHashMapKey(t *testing.T) {
	const (
		hname    = "hk_abc123_test_test_test_123abc_exp"
		testkey  = "token"
		testval  = "123abc"
		username = "bob"
	)
	hm := NewHashMap(pool, hname)

	// Check that the list qualifies for the IList interface
	var _ pinterface.IHashMap = hm

	hm.SelectDatabase(1)

	if err := hm.SetExpire(username, testkey, testval, time.Second*1); err != nil {
		t.Errorf("Error, could not set key and value! %s", err.Error())
	}
	retval, err := hm.Get(username, testkey)
	if err != nil {
		t.Errorf("Error, could not get value! %s", err.Error())
	} else if retval != testval {
		t.Errorf("Error, got the wrong return value! %s", retval)
	}
	// Wait a bit more than just 1 second. Testing on Travis can take some time.
	time.Sleep(3 * time.Second)

	_, err2 := hm.Get(username, testkey)
	if err2 == nil {
		t.Errorf("Error, key should be gone! %s", testkey)
	}
	err = hm.Remove()
	if err != nil {
		t.Errorf("Error, could not remove Hash! %s", err.Error())
	}
}

func TestHashMap(t *testing.T) {
	const (
		hashname  = "abc123_test_test_test_123abc_123"
		testid    = "bob"
		testidInv = "b:ob"
		testkey   = "password"
		testvalue = "hunter1"
	)

	hash := NewHashMap(pool, hashname)

	// Check that the list qualifies for the IList interface
	var _ pinterface.IHashMap = hash

	hash.SelectDatabase(1)
	hash.Clear()

	//if err := hash.Set(testidInv, testkey, testvalue); err == nil {
	//	t.Error("Should not be allowed to use an element id with \":\"")
	//}
	if err := hash.Set(testid, testkey, testvalue); err != nil {
		t.Errorf("Error, could not add item to hash map! %s", err.Error())
	}
	value2, err := hash.Get(testid, testkey)
	if err != nil {
		t.Error(err)
	}
	if value2 != testvalue {
		t.Errorf("Got a different value in return! %s != %s", value2, testvalue)
	}
	items, err := hash.All()
	if err != nil {
		t.Error(err)
	}
	if len(items) != 1 {
		t.Errorf("Error, wrong hash map length! %v", len(items))
	}
	if (len(items) > 0) && (items[0] != testid) {
		t.Errorf("Error, wrong hash map id! %v", items)
	}
	props, err := hash.Keys(testid)
	if err != nil {
		t.Error(err)
	}
	// only "password"
	if len(props) != 1 {
		t.Errorf("Error, wrong properties: %v\n", props)
	}
	if props[0] != "password" {
		t.Errorf("Error, wrong properties: %v\n", props)
	}

	err = hash.Remove()
	if err != nil {
		t.Errorf("Error, could not remove hash map! %s", err.Error())
	}
}

func TestHashMapFindByValue(t *testing.T) {
	const (
		hashname      = "abc123_test_test_test_123abc_123"
		elementID     = "bob"
		keyEmail      = "email"
		valueEmail    = "bob@zombo.com"
		keyUsername   = "username"
		valueUsername = "bob"
	)

	// Create a new connection pool
	pool := NewConnectionPool()
	defer pool.Close()

	// Create a new HashMap instance
	hash := NewHashMap(pool, hashname)
	hash.SelectDatabase(1)

	// Ensure the hash map is clean before the test
	if err := hash.Remove(); err != nil {
		t.Errorf("Error removing hash map: %v", err)
	}

	// Set key-value pairs for the elementID
	if err := hash.Set(elementID, keyEmail, valueEmail); err != nil {
		t.Errorf("Error setting email: %v", err)
	}
	if err := hash.Set(elementID, keyUsername, valueUsername); err != nil {
		t.Errorf("Error setting username: %v", err)
	}

	// Use FindKeyByValue to find the key associated with the email value
	foundKey, err := hash.FindKeyByValue(elementID, valueEmail)
	if err != nil {
		t.Errorf("Error finding key by value: %v", err)
	} else if foundKey != keyEmail {
		t.Errorf("Expected to find key '%s', but found '%s'", keyEmail, foundKey)
	} else {
		t.Logf("Successfully found key '%s' for value '%s'", foundKey, valueEmail)
	}

	// Retrieve the username using the found key (if needed)
	retrievedUsername, err := hash.Get(elementID, keyUsername)
	if err != nil {
		t.Errorf("Error getting username: %v", err)
	} else if retrievedUsername != valueUsername {
		t.Errorf("Expected username '%s', but got '%s'", valueUsername, retrievedUsername)
	} else {
		t.Logf("Successfully retrieved username '%s'", retrievedUsername)
	}

	// Clean up after the test
	if err := hash.Remove(); err != nil {
		t.Errorf("Error removing hash map: %v", err)
	}
}
