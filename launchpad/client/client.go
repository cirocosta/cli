package client

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/launchpad-project/cli/launchpad"
)

var (
	client             = &http.Client{}
	UserAgent          = "Launchpad/" + launchpad.Version + " (+https://launchpad.io)"
	UnexpectedResponse = errors.New("Unexpected response")
)

type Client struct {
	Id          int
	Time        time.Time
	RequestBody interface{}
	Request     *http.Request
	Response    *http.Response
}

func Url(uri string, body io.Reader) *Client {
	var time = time.Now()
	rand.Seed(time.UTC().UnixNano())

	var c = &Client{
		Id:   rand.Int(),
		Time: time,
	}

	req, err := http.NewRequest("GET", uri, body)

	c.Request = req

	if err != nil {
		panic(err)
	}

	c.Request.Header.Set("User-Agent", UserAgent)

	return c
}

func (c *Client) Auth(username, password string) *Client {
	c.Request.SetBasicAuth(username, password)
	return c
}

func (c *Client) ResponseJson(class interface{}) error {
	return json.NewDecoder(c.Response.Body).Decode(class)
}

func (c *Client) action(method string) *Client {
	c.Request.Method = method
	var resp, err = client.Do(c.Request)

	if err != nil {
		panic(err)
	}

	c.Response = resp

	if resp.StatusCode != http.StatusOK {
		panic(UnexpectedResponse.Error())
	}

	return c
}

func (c *Client) Get() *Client {
	return c.action("GET")
}

func (c *Client) Post() *Client {
	return c.action("POST")
}

func (c *Client) HEAD() *Client {
	return c.action("HEAD")
}

func (c *Client) Delete() *Client {
	return c.action("DELETE")
}

func (c *Client) Patch() *Client {
	return c.action("PATCH")
}
