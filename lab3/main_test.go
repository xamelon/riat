package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPostInputData(t *testing.T) {
	go startServer("8080")
	inputObject := Input{
		10,
		[]float32{1.01, 2.02},
		[]int{1, 4},
	}
	outputString := "{\"SumResult\":30.300001,\"MulResult\":4,\"SortedInputs\":[1,1.01,2.02,4]}"
	b, err := json.Marshal(inputObject)
	if err != nil {
		t.Error(err)
	}
	buf := bytes.NewReader(b)
	resp, err := http.Post("http://127.0.0.1:8080/PostInputData", "application/json", buf)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("Status is not 200")
	}
	resp, err = http.Get("http://127.0.0.1:8080/GetAnswer")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Error("Status is not 200")
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if string(body) != outputString {
		t.Error("JSON strings are not equal")
	}
	_, err = http.Get("http://127.0.0.1:8080/Stop")
	if err != nil {
		t.Error(err)
	}
}
