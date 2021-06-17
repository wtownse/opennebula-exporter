package api

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"github.com/bogue1979/xmlrpc"
)

// RPC represents an XML-RPC client.
type RPC struct {
	Client     *xmlrpc.Client
	URL        string
	AuthString string
}

// NewClient returns a prepared XML-RPC client.
func NewClient(url, user, password string, transport http.RoundTripper, timeout time.Duration) (*RPC, error) {
	client, err := xmlrpc.NewClient(url, transport, timeout)
	if err != nil {
		return nil, err
	}
	return newClient(url, user, password, client)
}

func newClient(url, user, password string, client *xmlrpc.Client) (*RPC, error) {
	return &RPC{
		Client:     client,
		URL:        url,
		AuthString: fmt.Sprintf("%s:%s", user, password),
	}, nil
}

// Call issues a request against the API endpoint.
func (c *RPC) Call(v interface{}, method string, args []interface{}) error {
	result := []interface{}{}

	if err := c.Client.Call(method, args, &result); err != nil {
		return err
	}

	apiCallSucceeded, ok := result[0].(bool)
	if !ok {
		return fmt.Errorf("malformed XMLRPC response")
	}
	if !apiCallSucceeded {
		switch e := result[1].(type) {
		case int64:
			return fmt.Errorf("API call against %s unsuccessful, error code %d", c.URL, e)
		case string:
			return fmt.Errorf("API call against %s unsuccessful, %s", c.URL, e)
		default:
			return fmt.Errorf("API call against %s unsuccessful", c.URL)
		}
	}

	switch r := result[1].(type) {
	case int64:
		return nil
	case string:
		return xml.Unmarshal([]byte(r), v)
	default:
		return fmt.Errorf("no known result type received from RPC call")
	}

}
