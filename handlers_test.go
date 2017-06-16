package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
