package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandlers_Status_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(status))
	defer ts.Close()
	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if string(body) != `{"status": "ok"}` {
		t.Error(string(body))
	}
	if res.StatusCode != http.StatusOK {
		t.Error(res.StatusCode)
	}
}

func TestHandlers_Status_Post(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(status))
	defer ts.Close()
	res, err := http.Post(ts.URL, "", nil)
	if err != nil {
		t.Error(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if string(body) != http.StatusText(http.StatusNotFound)+"\n" {
		t.Errorf(string(body))
	}
	if res.StatusCode != http.StatusNotFound {
		t.Error(res.StatusCode)
	}
}

func TestHandlers_Message_Get(t *testing.T) {
	ts := httptest.NewServer(message("test"))
	defer ts.Close()
	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if string(body) != `test` {
		t.Error(string(body))
	}
	if res.StatusCode != http.StatusOK {
		t.Error(res.StatusCode)
	}
}

func TestHandlers_Message_Post(t *testing.T) {
	ts := httptest.NewServer(message("test"))
	defer ts.Close()
	res, err := http.Post(ts.URL, "", nil)
	if err != nil {
		t.Error(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if string(body) != http.StatusText(http.StatusNotFound)+"\n" {
		t.Errorf(string(body))
	}
	if res.StatusCode != http.StatusNotFound {
		t.Error(res.StatusCode)
	}
}

func TestHandlers_Message_Get_LogRequest(t *testing.T) {
	ts := httptest.NewServer(logRequest(message("test")))
	defer ts.Close()
	var buf bytes.Buffer

	log.SetOutput(&buf)
	log.SetFlags(0)
	res, err := http.Get(ts.URL)
	log.SetFlags(log.LstdFlags)
	log.SetOutput(os.Stderr)

	if err != nil {
		t.Error(err)
	}
	if buf.String() != "GET / 200\n" {
		t.Errorf("%#v", buf.String())
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if string(body) != `test` {
		t.Error(string(body))
	}
	if res.StatusCode != http.StatusOK {
		t.Error(res.StatusCode)
	}
}
