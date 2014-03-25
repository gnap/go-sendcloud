/*
Sendcloud client in Go.
*/
package sendcloud

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	API_ENDPOINT = "https://sendcloud.sohu.com/webapi/"
)

type Client struct {
	domains map[string]struct { // sending domains
		api_user string
		api_key  string
	}
    logger ErrorLogger
}

func New() *Client {
	d := make(map[string]struct {
		api_user string
		api_key  string
	})
	l := FmtErrorLogger{}
	return &Client{domains: d, logger: l}
}

// add a sending domain with its authentication info
func (c *Client) AddDomain(domain, api_user, api_key string) {
	c.domains[domain] = struct {
		api_user string
		api_key  string
	}{api_user, api_key}
}

func (c *Client) SetLogger(l ErrorLogger) {
	c.logger = l
}

// invoke the remote API
func (c *Client) do(target, domain string, data url.Values) (body []byte, err error) {
	url := fmt.Sprintf("%s%s.json", API_ENDPOINT, target)
	s, ok := c.domains[domain]
	if !ok {
		return nil, fmt.Errorf("unknown domain: %s", domain)
	}
	data.Add("api_user", s.api_user)
	data.Add("api_key", s.api_key)
	rsp, err := http.PostForm(url, data)
	if err != nil {
		return
	}
	defer rsp.Body.Close()
	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}
	if rsp.StatusCode != 200 {
		//err = fmt.Errorf("SendCloud error: %d %s", rsp.StatusCode, body)
        msg := string(body[:])
		err = c.logger.ErrorLog("SendCloud SMTP", rsp.StatusCode, msg)
	}
	return
}
