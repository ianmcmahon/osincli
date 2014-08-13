package osincli

import (
	"fmt"
	"errors"
	"reflect"
	"strconv"
	"net/url"
)

type InfoRequest struct {
	client 				*Client
	token				string
	CustomParameters 	map[string]string
}

type InfoData struct {
	TokenType    string
	AccessToken string
	RefreshToken string
	Expiration *int32
	ResponseData ResponseData
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


func (c *InfoRequest) GetInfoData() (*InfoData, error) {
	iu := c.GetInfoUrl()

	ret := &InfoData{
		ResponseData: make(ResponseData),
	}

	// download data
	m := "POST"
	if c.client.config.UseGetAccessRequest {
		m = "GET"
	}
	err := downloadData(m, iu, nil, c.client.Transport, ret.ResponseData)
	if err != nil {
		return nil, err
	}

	// extract and convert received data
	token_type, ok := ret.ResponseData["token_type"]
	if !ok {
		return nil, errors.New("Invalid parameters received")
	}
	ret.TokenType = fmt.Sprintf("%v", token_type)

	access_token, ok := ret.ResponseData["access_token"]
	if !ok {
		return nil, errors.New("Invalid parameters received")
	}
	ret.AccessToken = fmt.Sprintf("%v", access_token)

	refresh_token, ok := ret.ResponseData["refresh_token"]
	if !ok {
		ret.RefreshToken = ""
	} else {
		ret.RefreshToken = fmt.Sprintf("%v", refresh_token)
	}

	expires_in_raw, ok := ret.ResponseData["expires_in"]
	if ok {
		rv := reflect.ValueOf(expires_in_raw)
		switch rv.Kind() {
		case reflect.Float64:
			// encoding/json always convert numbers fo float64
			ret.Expiration = new(int32)
			*ret.Expiration = int32(rv.Float())
		case reflect.String:
			// if string convert to integer
			ei, err := strconv.ParseInt(rv.String(), 10, 32)
			if err != nil {
				return nil, err
			}
			ret.Expiration = new(int32)
			*ret.Expiration = int32(ei)
		default:
			return nil, errors.New("Invalid parameter value")
		}
	}

	return ret, nil
}
