package containers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/launchpad-project/cli/globalconfigmock"
	"github.com/launchpad-project/cli/servertest"
)

var bufOutStream bytes.Buffer

func TestMain(m *testing.M) {
	var defaultOutStream = outStream
	outStream = &bufOutStream

	ec := m.Run()

	outStream = defaultOutStream
	os.Exit(ec)
}

func TestList(t *testing.T) {
	defer servertest.Teardown()
	servertest.Setup()
	globalconfigmock.Setup()
	bufOutStream.Reset()

	var want = "Cloud Search 145911644072551330 (launchpad)\n"

	servertest.Mux.HandleFunc("/api/projects/123/containers", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,
			`{
	"456": {
		"template": "cloudsearch",
		"name": "Cloud Search",
		"image": "launchpad",
		"id": "145911644072551330",
		"basePath": "/search"
    }
}`)
	})

	List("123")

	if bufOutStream.String() != want {
		t.Errorf("Wanted %v, got %v instead", want, bufOutStream.String())
	}

	globalconfigmock.Teardown()
}
