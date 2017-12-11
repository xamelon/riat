package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

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

type Input struct {
	K    int       `json:"K"`
	Sums []float32 `json:"Sums"`
	Muls []int     `json:"Muls"`
}

type Output struct {
	SumResult    float32   `json:"SumResult" xml:"SumResult`
	MulResult    float32   `json:"MulResult" xml:"MulResult"`
	SortedInputs []float32 `json:"SortedInputs" xml:"SortedInputs>decimal"`
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

var input Input
var server *http.Server

func Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func PostInputData(w http.ResponseWriter, r *http.Request) {
	resp, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(resp, &input)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetAnswer(w http.ResponseWriter, r *http.Request) {
	output := Calculate(input)
	jsonStr, err := json.Marshal(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(jsonStr))
	w.WriteHeader(http.StatusOK)
}

func Stop(w http.ResponseWriter, r *http.Request) {
	err := server.Shutdown(nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func startServer(port string) *http.Server {
	port = fmt.Sprintf(":%s", port)
	srv := &http.Server{Addr: port}
	http.HandleFunc("/Ping", Ping)
	http.HandleFunc("/PostInputData", PostInputData)
	http.HandleFunc("/GetAnswer", GetAnswer)
	log.Fatal(srv.ListenAndServe())
	return srv
}

func main() {
	port := "8080"
	server = startServer(port)
}
