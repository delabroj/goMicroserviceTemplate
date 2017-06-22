package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type endpointTest struct {
	Method       string
	Path         string
	RequestBody  io.Reader
	ResponseBody string
	Status       int
}

func (e endpointTest) String() string {
	result := fmt.Sprintf("Method: %#v\n", e.Method)
	result += fmt.Sprintf("Path: %#v\n", e.Path)
	result += fmt.Sprintf("RequestBody: %#v\n", e.RequestBody)
	result += fmt.Sprintf("ResponseBody: %#v\n", e.ResponseBody)
	result += fmt.Sprintf("Status: %#v\n", e.Status)
	return result
}

type endpointTests []endpointTest

func HandlerTest(t *testing.T, handler http.Handler, handlerName string, tests endpointTests) {
	server := httptest.NewServer(handler)
	defer server.Close()
	for _, test := range tests {
		context := fmt.Sprintf("\nhandler: %v\n\nendpointTest: %v", handlerName, test)
		r, err := http.NewRequest(test.Method, server.URL+test.Path, test.RequestBody)
		if err != nil {
			t.Errorf(context+"\nerror: %v", err)
		}
		response, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Errorf(context+"\nerror: %v", err)
		}
		actualBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Errorf(context+"\nerror: %v", err)
		}
		actualBodyString := string(actualBody)
		if test.ResponseBody != actualBodyString {
			t.Errorf(context+"\ntest.ResponseBody != actualBodyString: %#v != %#v", test.ResponseBody, actualBodyString)
		}
		if test.Status != response.StatusCode {
			t.Errorf(context+"\ntest.Status != response.StatusCode: %v != %v", test.Status, response.StatusCode)
		}
	}
}

var statusTests = endpointTests{{
	Method:       "GET",
	Path:         "/status",
	ResponseBody: `{"status": "ok"}`,
	Status:       http.StatusOK,
}, {
	Method:       "POST",
	Path:         "/status",
	RequestBody:  strings.NewReader(`request body`),
	ResponseBody: http.StatusText(http.StatusNotFound) + "\n",
	Status:       http.StatusNotFound,
}}

func TestStatus(t *testing.T) {
	HandlerTest(t, http.HandlerFunc(status), "http.HandlerFunc(status)", statusTests)
}

var messageTests = endpointTests{{
	Method:       "GET",
	Path:         "/time",
	ResponseBody: `test`,
	Status:       http.StatusOK,
}, {
	Method:       "POST",
	Path:         "/time",
	RequestBody:  strings.NewReader(`request body`),
	ResponseBody: http.StatusText(http.StatusNotFound) + "\n",
	Status:       http.StatusNotFound,
}}

func TestMessage(t *testing.T) {
	HandlerTest(t, message("test"), `message("test")`, messageTests)
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
