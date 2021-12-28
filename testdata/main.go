// This program is a test program used to facilitate unit testing with Tarmac.
package main

import (
	"encoding/base64"
	"fmt"
	"github.com/valyala/fastjson"
	wapc "github.com/wapc/wapc-guest-tinygo"
)

func main() {
	// Tarmac uses waPC to facilitate WASM module execution. Modules must register their custom handlers under the
	// appropriate method as shown below.
	wapc.RegisterFunctions(wapc.Functions{
		// Register a GET request handler
		"http:GET": NoHandler,
		// Register a POST request handler
		"http:POST": Handler,
		// Register a PUT request handler
		"http:PUT": Handler,
		// Register a DELETE request handler
		"http:DELETE": NoHandler,
		// Register a handler for scheduled tasks
		"scheduler:RUN": Handler,
	})
}

// NoHandler is a custom Tarmac Handler function that will return a tarmac.ServerResponse JSON that denies
// the client request.
func NoHandler(payload []byte) ([]byte, error) {
	return []byte(`{"status":{"code":503,"status":"Not Implemented"}}`), nil
}

// Handler is the custom Tarmac Handler function that will receive a tarmac.ServerRequest JSON payload and
// must return a tarmac.ServerResponse JSON payload along with a nil error.
func Handler(payload []byte) ([]byte, error) {
	// Parse the JSON request
	rq, err := fastjson.ParseBytes(payload)
	if err != nil {
		return []byte(fmt.Sprintf(`{"status":{"code":500,"status":"Failed to call parse json - %s"}}`, err)), nil
	}

	// Decode the payload
	s, err := base64.StdEncoding.DecodeString(string(rq.GetStringBytes("payload")))
	if err != nil {
		return []byte(fmt.Sprintf(`{"status":{"code":500,"status":"Failed to perform base64 decode - %s"}}`, err)), nil
	}
	b := []byte(s)

	// Logger
	_, err = wapc.HostCall("tarmac", "logger", "debug", payload)
	if err != nil {
		return []byte(fmt.Sprintf(`{"status":{"code":500,"status":"Failed to call host callback - %s"}}`, err)), nil
	}

	/*
		// HTTP Client
		r, err := wapc.HostCall("tarmac", "httpclient", "call", []byte(`{"method":"GET","insecure":true,"url":"http://localhost:9000/health"}`))
		if err != nil {
			return []byte(fmt.Sprintf(`{"status":{"code":500,"status":"Failed to call host callback"}}`)), nil
		}

		hr, err := fastjson.ParseBytes(r)
		if err != nil {
			return []byte(fmt.Sprintf(`{"status":{"code":500,"status":"Failed to parse host callback json"}}`)), nil
		}

		if hr.GetInt("code") == 200 {
			return []byte(fmt.Sprintf(`{"status":{"code":500,"status":"Host callback results failure"}}`)), nil
		}
	*/

	// KVStore
	_, err = wapc.HostCall("tarmac", "kvstore", "set", []byte(`{"key":"test-data","data":"aSBhbSBhIGxpdHRsZSB0ZWFwb3Q="}`))
	if err != nil {
		return []byte(fmt.Sprintf(`{"status":{"code":500,"status":"Failed to call host callback - %s"}}`, err)), nil
	}

	// Return the payload via a ServerResponse JSON
	return []byte(fmt.Sprintf(`{"payload":"%s","status":{"code":200,"status":"Success"}}`, base64.StdEncoding.EncodeToString(b))), nil
}
