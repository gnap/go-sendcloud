/*
Sendcloud client in Go.
*/
package sendcloud

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	API_ENDPOINT = "https://sendcloud.sohu.com/webapi/"
	HTTP_TIMEOUT = 10 * time.Second
)

type Client struct {
	domains map[string]struct { // sending domains
		api_user string
		api_key  string
	}
	logger ErrorLogger
	httpClient *http.Client
}

func New() *Client {
	d := make(map[string]struct {
		api_user string
		api_key  string
	})
	l := FmtErrorLogger{}
	tr := &http.Transport{MaxIdleConnsPerHost: 10}
	return &Client{domains: d, logger: l, httpClient: &http.Client{Transport: tr, Timeout: time.Duration(HTTP_TIMEOUT)}}
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
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 30 * time.Second,
	}
	client := &http.Client{Transport: tr}
	rsp, err := client.PostForm(url, data)
	if err != nil {
		return
	}
	defer rsp.Body.Close()
	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}
	//err = fmt.Errorf("SendCloud error: %d %s", rsp.StatusCode, body)
	msg := string(body[:])
	go c.logger.ErrorLog("sendcloud.error", rsp.StatusCode, msg)
	return
}
