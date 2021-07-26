package app

import (
	"encoding/base64"
	"fmt"
	"github.com/madflojo/tarmac"
	"github.com/pquerna/ffjson/ffjson"
	"time"
)

// kvStore provides access to Host Callbacks that interact with the key:value datastores within Tarmac. The callbacks
// within kvStore provided all of the logic and error handlings of accessing and interacting with a key:value
// database. Users will send the specified JSON request and receive an appropriate JSON response.
type kvStore struct{}

// Get will fetch the stored data using the key specified within the incoming JSON. Logging, error handling, and
// base64 encoding of data are all handled via this function. Note, this function expects the KVStoreGetRequest
// JSON type as input and will return a KVStoreGetResponse JSON.
func (k *kvStore) Get(b []byte) ([]byte, error) {
	now := time.Now()
	// Start Response Message assuming everything is good
	r := tarmac.KVStoreGetResponse{}
	r.Status.Code = 200
	r.Status.Status = "OK"

	// Parse incoming Request
	var rq tarmac.KVStoreGet
	err := ffjson.Unmarshal(b, &rq)
	if err != nil {
		r.Status.Code = 400
		r.Status.Status = "Error Parsing Input"
	}

	// Fetch data from KVStore if we do not have any other errors
	if r.Status.Code == 200 {
		data, err := kv.Get(rq.Key)
		if err != nil {
			r.Status.Code = 404
			r.Status.Status = fmt.Sprintf("Unable to fetch key %s - %s", rq.Key, err)
		}

		// Encode Fetched Data to store within JSON
		r.Data = base64.StdEncoding.EncodeToString(data)

	}

	// Marshal a resposne JSON to return to caller
	rsp, err := ffjson.Marshal(r)
	if err != nil {
		log.Errorf("Unable to marshal kvstore:get response - %s", err)
		stats.kvstore.WithLabelValues("get").Observe(time.Since(now).Seconds())
		return []byte(""), fmt.Errorf("unable to marshal kvstore:get response")
	}

	if r.Status.Code == 200 {
		stats.kvstore.WithLabelValues("get").Observe(time.Since(now).Seconds())
		return rsp, nil
	}
	stats.kvstore.WithLabelValues("get").Observe(time.Since(now).Seconds())
	return rsp, fmt.Errorf("%s", r.Status.Status)
}

// Set will store data within the key:value datastore using the key specified within the incoming JSON. Logging, error
// handling, and base64 decoding of data are all handled via this function. Note, this function expects the
// KVStoreSetRequest JSON type as input and will return a KVStoreSetResponse JSON.
func (k *kvStore) Set(b []byte) ([]byte, error) {
	now := time.Now()
	// Start Response Message assuming everything is good
	r := tarmac.KVStoreSetResponse{}
	r.Status.Code = 200
	r.Status.Status = "OK"

	// Parse incoming Request
	var rq tarmac.KVStoreSet
	err := ffjson.Unmarshal(b, &rq)
	if err != nil {
		r.Status.Code = 400
		r.Status.Status = "Error Parsing Input"
	}

	// Decode data to store
	data, err := base64.StdEncoding.DecodeString(rq.Data)
	if err != nil {
		r.Status.Code = 400
		r.Status.Status = fmt.Sprintf("Unable to decode data - %s", err)
	}

	// Store data in KVStore if we do not have any other errors
	if r.Status.Code == 200 {
		err = kv.Set(rq.Key, data)
		if err != nil {
			r.Status.Code = 500
			r.Status.Status = fmt.Sprintf("Unable to store data using key %s - %s", rq.Key, err)
		}
	}

	// Marshal a resposne JSON to return to caller
	rsp, err := ffjson.Marshal(r)
	if err != nil {
		log.Errorf("Unable to marshal kvstore:set response - %s", err)
		stats.kvstore.WithLabelValues("set").Observe(time.Since(now).Seconds())
		return []byte(""), fmt.Errorf("unable to marshal kvstore:set response")
	}

	if r.Status.Code == 200 {
		stats.kvstore.WithLabelValues("set").Observe(time.Since(now).Seconds())
		return rsp, nil
	}
	stats.kvstore.WithLabelValues("set").Observe(time.Since(now).Seconds())
	return rsp, fmt.Errorf("%s", r.Status.Status)
}

// Delete will remove the key and data stored within the key:value datastore using the key specified within the incoming
// JSON. Logging and error handling are all handled via this callback function. Note, this function expects the
// KVStoreDeleteRequest JSON type as input and will return a KVStoreDeleteResponse JSON.
func (k *kvStore) Delete(b []byte) ([]byte, error) {
	now := time.Now()
	// Start Response Message assuming everything is good
	r := tarmac.KVStoreDeleteResponse{}
	r.Status.Code = 200
	r.Status.Status = "OK"

	// Parse incoming Request
	var rq tarmac.KVStoreDelete
	err := ffjson.Unmarshal(b, &rq)
	if err != nil {
		r.Status.Code = 400
		r.Status.Status = "Error Parsing Input"
	}

	// Delete data in KVStore if we do not have any other errors
	if r.Status.Code == 200 {
		err = kv.Delete(rq.Key)
		if err != nil {
			r.Status.Code = 404
			r.Status.Status = fmt.Sprintf("Unable to delete key %s - %s", rq.Key, err)
		}
	}

	// Marshal a response JSON to return to caller
	rsp, err := ffjson.Marshal(r)
	if err != nil {
		log.Errorf("Unable to marshal kvstore:delete response - %s", err)
		stats.kvstore.WithLabelValues("delete").Observe(time.Since(now).Seconds())
		return []byte(""), fmt.Errorf("unable to marshal kvstore:delete response")
	}

	if r.Status.Code == 200 {
		stats.kvstore.WithLabelValues("delete").Observe(time.Since(now).Seconds())
		return rsp, nil
	}
	stats.kvstore.WithLabelValues("delete").Observe(time.Since(now).Seconds())
	return rsp, fmt.Errorf("%s", r.Status.Status)
}

// Keys will return a list of all keys stored within the key:value datastore. Logging and error handling are
// all handled via this callback function. Note, this function will return a KVStoreKeysResponse JSON.
func (k *kvStore) Keys(b []byte) ([]byte, error) {
	now := time.Now()
	// Start Response Message assuming everything is good
	r := tarmac.KVStoreKeysResponse{}
	r.Status.Code = 200
	r.Status.Status = "OK"

	// Fetch keys from datastore
	var err error
	r.Keys, err = kv.Keys()
	if err != nil {
		r.Status.Code = 500
		r.Status.Status = fmt.Sprintf("Unable to fetch keys - %s", err)
	}

	// Marshal a response JSON to return to caller
	rsp, err := ffjson.Marshal(r)
	if err != nil {
		log.Errorf("Unable to marshal kvstore:delete response - %s", err)
		stats.kvstore.WithLabelValues("keys").Observe(time.Since(now).Seconds())
		return []byte(""), fmt.Errorf("unable to marshal kvstore:delete response")
	}

	if r.Status.Code == 200 {
		stats.kvstore.WithLabelValues("keys").Observe(time.Since(now).Seconds())
		return rsp, nil
	}
	stats.kvstore.WithLabelValues("keys").Observe(time.Since(now).Seconds())
	return rsp, fmt.Errorf("%s", r.Status.Status)
}
