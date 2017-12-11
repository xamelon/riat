package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"
)

type Input struct {
	K    int       `json:"K"`
	Sums []float32 `json:"Sums"`
	Muls []int     `json:"Muls"`
}

type Output struct {
	SumResult    float32   `json:"SumResult"`
	MulResult    float32   `json:"MulResult"`
	SortedInputs []float32 `json:"SortedInputs"`
}

type Config struct {
	Url string
}

func GetInputData(config Config) Input {
	url := fmt.Sprintf("%s/GetInputData", config.Url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var input Input
	err = json.Unmarshal(body, &input)
	if err != nil {
		log.Fatal(err)
	}
	return input
}

func WriteAnswer(config Config, output Output) {
	url := fmt.Sprintf("%s/WriteAnswer", config.Url)
	b, err := json.Marshal(output)
	if err != nil {
		log.Fatal(err)
	}
	body := bytes.NewReader(b)
	_, err = http.Post(url, "application/json", body)
	if err != nil {
		log.Fatal(err)
	}
}

type fl []float32

func (s fl) Len() int {
	return len(s)
}
func (s fl) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s fl) Less(i, j int) bool {
	return s[i] < s[j]
}

func Calculate(input Input) Output {
	var sortedInputs []float32
	sortedInputs = append(sortedInputs, input.Sums[:len(input.Sums)]...)
	var sumResults float32 = 0
	var mulResults float32 = 1.0
	for i := 0; i < len(input.Sums); i++ {
		sumResults += input.Sums[i] * float32(input.K)
	}

	for i := 0; i < len(input.Muls); i++ {
		mulResults *= float32(input.Muls[i])
		sortedInputs = append(sortedInputs, float32(input.Muls[i]))
	}
	sort.Sort(fl(sortedInputs))
	var output = Output{
		SumResult:    sumResults,
		MulResult:    mulResults,
		SortedInputs: sortedInputs,
	}

	return output
}

func Ping(config Config) int {
	url := fmt.Sprintf("%s/Ping", config.Url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	return resp.StatusCode
}

func main() {
	var port string
	_, err := fmt.Scanf("%s", &port)
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("htpp://127.0.0.1:%s/", port)
	config := Config{url}
	var status int
	for {
		if status != http.StatusOK {
			break
		}
		status = Ping(config)
		time.Sleep(1 * time.Second)
	}

	input := GetInputData(config)
	output := Calculate(input)
	WriteAnswer(config, output)

}
