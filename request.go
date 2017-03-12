package ulule

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) apiget(route string, res interface{}) error {
	req, err := http.NewRequest("GET", "https://api.ulule.com/v1"+route, nil)
	if err != nil {
		return err
	}
	c.authenticate(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("error %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(res)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) authenticate(req *http.Request) {
	if c.username != "" && c.password != "" {
		// TODO: basic auth
		// curl --basic --user "username:password" https://api.ulule.com/v1/...
	} else if c.username != "" && c.apikey != "" {
		req.Header.Add("Authorization", "Bearer "+c.username+":"+c.apikey)
	} else if c.accessToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.accessToken)
	}
}
