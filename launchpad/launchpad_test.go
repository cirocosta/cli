package launchpad

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var mux *http.ServeMux
var server *httptest.Server

func resetServer() {
	if server != nil {
		server.Close()
	}

	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	defaultClient = &http.Client{Transport: transport}
}

func assertStatusCode(t *testing.T, want int, got int) {
	if got != want {
		t.Errorf("Expected status code %d, got %d", want, got)
	}
}

func assertMethod(t *testing.T, want string, got string) {
	if got != want {
		t.Errorf("%s method expected, found %s instead", want, got)
	}
}

func assertURI(t *testing.T, want string, got string) {
	if got != want {
		t.Errorf("Expected URL %s, got %s", want, got)
	}
}

func assertBodyResponse(t *testing.T, want string, got io.ReadCloser) {
	body, err := ioutil.ReadAll(got)

	if err != nil {
		t.Error(err)
	}

	var bString = string(body)

	if bString != want {
		t.Errorf("Expected %s response from the server, got %s instead", want, bString)
	}
}

func TestURL(t *testing.T) {
	var r, err = URL("https://example.com/foo/bah", nil)

	if err != nil {
		t.Error(err)
	}

	assertURI(t, "https://example.com/foo/bah", r.Request.URL.String())
}

func TestURLErrorDueToInvalidURI(t *testing.T) {
	_, err := URL("://example.com/foo/bah", nil)

	var kind = reflect.TypeOf(err).String()

	if kind != "*url.Error" {
		t.Errorf("Expected error %s isn't *URL.Error", kind)
	}
}

func TestUserAgent(t *testing.T) {
	var r, err = URL("http://localhost/foo", nil)

	if err != nil {
		t.Error(err)
	}

	var actual = r.Request.Header.Get("User-Agent")
	var expected = "Launchpad/master (+https://launchpad.io)"

	if actual != expected {
		t.Errorf("Expected User-Agent %s doesn't match with %s", actual, expected)
	}
}

func TestBasicAuth(t *testing.T) {
	var r, err = URL("http://localhost/", nil)

	if err != nil {
		t.Error(err)
	}

	r.Auth("admin", "safe")

	var username, password, ok = r.Request.BasicAuth()

	if username != "admin" || password != "safe" || ok != true {
		t.Errorf("Wrong user credentials")
	}
}

func TestDeleteRequest(t *testing.T) {
	resetServer()

	mux.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `"body"`)
	})

	var wantURL = "http://example.com/url"
	var wantMethod = "DELETE"

	var req, err = URL(wantURL, nil)

	if err != nil {
		t.Error(err)
	}

	if err := req.Delete(); err != nil {
		t.Error(err)
	}

	assertURI(t, wantURL, req.Request.URL.String())
	assertMethod(t, wantMethod, req.Request.Method)
	assertStatusCode(t, 200, req.Response.StatusCode)
	assertBodyResponse(t, `"body"`, req.Response.Body)
}

func TestGetRequest(t *testing.T) {
	resetServer()

	mux.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `"body"`)
	})

	var wantURL = "http://example.com/url"
	var wantMethod = "GET"

	var req, err = URL(wantURL, nil)

	if err != nil {
		t.Error(err)
	}

	if err := req.Get(); err != nil {
		t.Error(err)
	}

	assertURI(t, wantURL, req.Request.URL.String())
	assertMethod(t, wantMethod, req.Request.Method)
	assertStatusCode(t, 200, req.Response.StatusCode)
	assertBodyResponse(t, `"body"`, req.Response.Body)
}

func TestHeadRequest(t *testing.T) {
	resetServer()

	mux.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
	})

	var wantURL = "http://example.com/url"
	var wantMethod = "HEAD"

	var req, err = URL(wantURL, nil)

	if err != nil {
		t.Error(err)
	}

	if err := req.Head(); err != nil {
		t.Error(err)
	}

	assertURI(t, wantURL, req.Request.URL.String())
	assertMethod(t, wantMethod, req.Request.Method)
	assertStatusCode(t, 200, req.Response.StatusCode)
}

func TestPatchRequest(t *testing.T) {
	resetServer()

	mux.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `"body"`)
	})

	var wantURL = "http://example.com/url"
	var wantMethod = "PATCH"

	var req, err = URL(wantURL, nil)

	if err != nil {
		t.Error(err)
	}

	if err := req.Patch(); err != nil {
		t.Error(err)
	}

	assertURI(t, wantURL, req.Request.URL.String())
	assertMethod(t, wantMethod, req.Request.Method)
	assertStatusCode(t, 200, req.Response.StatusCode)
	assertBodyResponse(t, `"body"`, req.Response.Body)
}

func TestPostRequest(t *testing.T) {
	resetServer()

	mux.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `"body"`)
	})

	var wantURL = "http://example.com/url"
	var wantMethod = "POST"

	var req, err = URL(wantURL, nil)

	if err != nil {
		t.Error(err)
	}

	if err := req.Post(); err != nil {
		t.Error(err)
	}

	assertURI(t, wantURL, req.Request.URL.String())
	assertMethod(t, wantMethod, req.Request.Method)
	assertStatusCode(t, 200, req.Response.StatusCode)
	assertBodyResponse(t, `"body"`, req.Response.Body)
}

func TestPutRequest(t *testing.T) {
	resetServer()

	mux.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `"body"`)
	})

	var wantURL = "http://example.com/url"
	var wantMethod = "PUT"

	var req, err = URL(wantURL, nil)

	if err != nil {
		t.Error(err)
	}

	if err := req.Put(); err != nil {
		t.Error(err)
	}

	assertURI(t, wantURL, req.Request.URL.String())
	assertMethod(t, wantMethod, req.Request.Method)
	assertStatusCode(t, 200, req.Response.StatusCode)
	assertBodyResponse(t, `"body"`, req.Response.Body)
}

func TestPostFormRequest(t *testing.T) {
	resetServer()

	var wantURL = "http://example.com/url"
	var wantMethod = "POST"
	var wantContentType = "application/x-www-form-urlencoded"
	var wantTitle = "foo"

	mux.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `"body"`)

		var gotContentType = r.Header.Get("Content-Type")
		var gotTitle = r.PostFormValue("title")

		if gotContentType != wantContentType {
			t.Errorf("Expected content type %s, got %s instead", wantContentType, gotContentType)
		}

		if gotTitle != wantTitle {
			t.Errorf("Expected title %s, got %s instead", wantTitle, gotTitle)
		}
	})

	var form = url.Values{}
	form.Add("title", wantTitle)

	var content = strings.NewReader(form.Encode())

	req, err := URL(wantURL, content)

	if err != nil {
		t.Error(err)
	}

	req.Request.Header.Add("Content-Type", wantContentType)

	if err := req.Post(); err != nil {
		t.Error(err)
	}

	assertURI(t, wantURL, req.Request.URL.String())
	assertMethod(t, wantMethod, req.Request.Method)
	assertStatusCode(t, 200, req.Response.StatusCode)
	assertBodyResponse(t, `"body"`, req.Response.Body)
}

func TestDecodeJSON(t *testing.T) {
	resetServer()

	var wantURL = "http://example.com/url"
	var wantMethod = "GET"
	var wantTitle = "body"

	mux.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"title": "body"}`)
	})

	req, err := URL(wantURL, nil)

	if err != nil {
		t.Error(err)
	}

	if err := req.Get(); err != nil {
		t.Error(err)
	}

	assertURI(t, wantURL, req.Request.URL.String())
	assertMethod(t, wantMethod, req.Request.Method)
	assertStatusCode(t, 200, req.Response.StatusCode)

	var content struct {
		Title string `json:"title"`
	}

	err = req.DecodeJSON(&content)

	if err != nil {
		t.Error(err)
	}

	if content.Title != wantTitle {
		t.Errorf("Expected title %s, got %s instead", wantTitle, content.Title)
	}
}

func TestErrorStatusCode404(t *testing.T) {
	resetServer()

	mux.HandleFunc("/url", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})

	var wantURL = "http://example.com/url"

	req, err := URL(wantURL, nil)

	if err != nil {
		t.Error(err)
	}

	if err := req.Get(); err != ErrUnexpectedResponse {
		t.Errorf("Missing error %s", ErrUnexpectedResponse)
	}

	assertURI(t, wantURL, req.Request.URL.String())
	assertMethod(t, "GET", req.Request.Method)
	assertStatusCode(t, 404, req.Response.StatusCode)
}
