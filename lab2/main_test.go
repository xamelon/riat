package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetInputData(t *testing.T) {
	jsonString := "{\"K\":10,\"Sums\":[1.01,2.02],\"Muls\":[1,4]}"
	inputObject := Input{
		10,
		[]float32{1.01, 2.02},
		[]int{1, 4},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(jsonString))
	}))
	defer ts.Close()
	config := Config{ts.URL}
	input := GetInputData(config)

	if inputObject.K != input.K {
		t.Error("K are not Equal")
	}

	if len(inputObject.Muls) != len(input.Muls) {
		t.Error("Muls are not Equal")
	}
	for i := 0; i < len(input.Muls); i++ {
		if inputObject.Muls[i] != input.Muls[i] {
			t.Error("Muls are not Equal")
		}
	}

	if len(inputObject.Sums) != len(input.Sums) {
		t.Error("Sums are not Equal")
	}
	for i := 0; i < len(input.Sums); i++ {
		if inputObject.Sums[i] != input.Sums[i] {
			t.Error("Sums are not Equal")
		}
	}
}

func TestWriteAnswer(t *testing.T) {
	output := Output{
		30.300001,
		4,
		[]float32{1, 1.01, 2.02, 4},
	}
	outputJson, err := json.Marshal(output)
	if err != nil {
		t.Error(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			t.Error(err)
		}

		if string(resp) != string(outputJson) {
			t.Error("JSON strings are not equal")
		}
	}))
	defer ts.Close()
	config := Config{ts.URL}
	WriteAnswer(config, output)
}

func TestCalculate(t *testing.T) {
	inputObject := Input{
		10,
		[]float32{1.01, 2.02},
		[]int{1, 4},
	}
	output := Output{
		30.300001,
		4,
		[]float32{1, 1.01, 2.02, 4},
	}
	outputObject := Calculate(inputObject)
	if output.SumResult != outputObject.SumResult {
		t.Error("SumResult are not equal")
	}

	if len(output.SortedInputs) != len(outputObject.SortedInputs) {
		t.Error("SortedInputs are not equal")
	}

	for i := 0; i < len(output.SortedInputs); i++ {
		if output.SortedInputs[i] != outputObject.SortedInputs[i] {
			t.Error("SortedInputs are not equal")
		}
	}

	if output.MulResult != outputObject.MulResult {
		t.Error("MulResult are not equals")
	}
}

func TestPing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	config := Config{ts.URL}
	status := Ping(config)
	if status != http.StatusOK {
		t.Error("Status is not 200")
	}
}
