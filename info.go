package osincli

import (
	"fmt"
	"net/url"
)

type InfoRequest struct {
	client 				*Client
	token				string
	CustomParameters 	map[string]string
}

// Creates a new info request
func (c *Client) NewInfoRequest(token string) *InfoRequest {
	return &InfoRequest{
		client:           c,
		token:			  token,
		CustomParameters: make(map[string]string),
	}
}

// returns the info url
func (c *InfoRequest) GetInfoUrl() *url.URL {
	return c.GetInfoUrlWithParams("")
}

// returns the info url
func (c *InfoRequest) GetInfoUrlWithParams(state string) *url.URL {
	u := *c.client.configcache.infoUrl
	uq := u.Query()
	uq.Add("code", c.token)

	if c.client.config.Scope != "" {
		uq.Add("scope", c.client.config.Scope)
	}

	if state != "" {
		uq.Add("state", state)
	}

	if c.CustomParameters != nil {
		for pn, pv := range c.CustomParameters {
			uq.Add(pn, pv)
		}
	}

	u.RawQuery = uq.Encode()
	return &u
}


func (c *InfoRequest) GetInfoData() (*ResponseData, error) {
	iu := c.GetInfoUrl()

	ret := make(ResponseData)

	// download data
	m := "POST"
	if c.client.config.UseGetAccessRequest {
		m = "GET"
	}
	err := downloadData(m, iu, nil, c.client.Transport, ret)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Received info response: %v\n", ret)

	return &ret, nil
}
