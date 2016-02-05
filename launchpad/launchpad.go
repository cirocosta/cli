package launchpad

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"time"
)

var (
	defaultClient = &http.Client{}
	// ErrUnexpectedResponse is used when an unexpected response happens
	ErrUnexpectedResponse = errors.New("Unexpected response")
)

// Launchpad is the structure for a Launchpad query
type Launchpad struct {
	ID          int
	Time        time.Time
	RequestBody interface{}
	Request     *http.Request
	Response    *http.Response
	httpClient  *http.Client
}

// URL creates a new request object
func URL(uri string, body io.Reader) (*Launchpad, error) {
	var time = time.Now()
	rand.Seed(time.UTC().UnixNano())

	var l = &Launchpad{
		ID:   rand.Int(),
		Time: time,
	}

	req, err := http.NewRequest("GET", uri, body)

	l.Request = req
	l.httpClient = defaultClient

	if err == nil {
		l.Request.Header.Set("User-Agent", UserAgent)
	}

	return l, err
}

// Auth sets HTTP basic auth headers
func (l *Launchpad) Auth(username, password string) *Launchpad {
	l.Request.SetBasicAuth(username, password)
	return l
}

// DecodeJSON decodes a JSON response
func (l *Launchpad) DecodeJSON(class interface{}) error {
	return json.NewDecoder(l.Response.Body).Decode(class)
}

func (l *Launchpad) action(method string) (err error) {
	l.Request.Method = method
	l.Response, err = l.httpClient.Do(l.Request)

	if err == nil && l.Response.StatusCode >= 400 {
		err = ErrUnexpectedResponse
	}

	return err
}

// Delete method
func (l *Launchpad) Delete() error {
	return l.action("DELETE")
}

// Get method
func (l *Launchpad) Get() error {
	return l.action("GET")
}

// Head method
func (l *Launchpad) Head() error {
	return l.action("HEAD")
}

// Patch method
func (l *Launchpad) Patch() error {
	return l.action("PATCH")
}

// Post method
func (l *Launchpad) Post() error {
	return l.action("POST")
}

// Put method
func (l *Launchpad) Put() error {
	return l.action("PUT")
}
