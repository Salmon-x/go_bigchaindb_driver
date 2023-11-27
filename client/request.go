package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/pkg/errors"
)

func StructToJSON(obj interface{}) (*bytes.Buffer, error) {
	// Convert struct to map
	result := make(map[string]interface{})
	value := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)
	for i := 0; i < value.NumField(); i++ {
		result[typ.Field(i).Name] = value.Field(i).Interface()
	}

	// Convert map to JSON string
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonBytes), nil
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {

	rel := &url.URL{Path: path}
	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		if method != http.MethodGet {
			buf = new(bytes.Buffer)
			err := json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}
		} else {
			q := u.Query()
			for k, v := range body.(map[string]string) {
				q.Set(k, v)
			}
			u.RawQuery = q.Encode()
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	req.Header = c.baseHeader

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > http.StatusIMUsed {
		return errors.New(resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&v)
	return err
}
