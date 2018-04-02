package bitrescli

import (
	"net/http"
	"time"
	"net/url"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"errors"
)

type Client struct {
	BaseURI string

	Debug   bool
	MaxIdleConnections int
	RequestTimeout     int

	httpClient *http.Client
}

func (client *Client) Connect() {
	client.httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: client.MaxIdleConnections,
		},
		Timeout: time.Duration(client.RequestTimeout) * time.Second,
	}
}

func (client Client) EndpointURL(endpoint string) (string, error) {
	baseUri, err := url.Parse(client.BaseURI)
	if err != nil {
		return "", err
	}

	path, err := url.Parse(fmt.Sprintf("/rest/%v.json", endpoint))
	if err != nil {
		return "", err
	}

	final := baseUri.ResolveReference(path)

	return final.String(), nil
}

func (client Client) request(method string, endpoint string, s interface{}) error {
	fullUrl, err := client.EndpointURL(endpoint)
	if err != nil {
		return err
	}

	if client.Debug {
		fmt.Println(fullUrl)
	}

	req, err := http.NewRequest(method, fullUrl, nil)
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(resp.StatusCode)

	if client.Debug {
		fmt.Println("Body:", string(body))
	}

	if err := json.Unmarshal(body, &s); err != nil {
		return errors.New(string(body))
	}

	return nil
}
