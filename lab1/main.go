package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"sort"
)

type TypeSerialization int

const (
	XML  TypeSerialization = 1 + iota
	JSON TypeSerialization = 2 + iota
)

type Input struct {
	K1    int       `json:"K" xml:"K"`
	Sums1 []float32 `json:"Sums" xml:"Sums>decimal"`
	Muls1 []int     `json:"Muls" xml:"Muls>int"`
}

type Output struct {
	SumResult    float32   `json:"SumResult" xml:"SumResult`
	MulResult    float32   `json:"MulResult" xml:"MulResult"`
	SortedInputs []float32 `json:"SortedInputs" xml:"SortedInputs>decimal"`
}

type Serializer interface {
	CanSerialize(_type TypeSerialization) bool
	Serialize(v interface{}) []byte
	Deserialize(jsonStr string) Input
}

type JSONSerializer struct{}

func (ser JSONSerializer) CanSerialize(_type TypeSerialization) bool {
	if _type == JSON {
		return true
	} else {
		return false
	}
}

func (ser JSONSerializer) Serialize(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

func (ser JSONSerializer) Deserialize(jsonStr string) Input {
	var obe Input
	err := json.Unmarshal([]byte(jsonStr), &obe)
	if err != nil {
		fmt.Println(err)
	}
	return obe
}

type XMLSerializer struct{}

func (ser XMLSerializer) CanSerialize(_type TypeSerialization) bool {
	if _type == XML {
		return true
	} else {
		return false
	}
}

func (ser XMLSerializer) Serialize(v interface{}) []byte {
	data, err := xml.Marshal(v)
	if err != nil {
		fmt.Println(err)
	}
	return data
}

func (ser XMLSerializer) Deserialize(xmlString string) Input {
	var obe Input
	err := xml.Unmarshal([]byte(xmlString), &obe)
	if err != nil {
		fmt.Println(err)
	}
	return obe
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
func main() {
	var SerializationType TypeSerialization

	var inputString string

	_, err := fmt.Scanf("%d", &SerializationType)
	if err != nil {
		fmt.Print(err)
	}
	var serializer Serializer
	if SerializationType == XML {
		serializer = XMLSerializer{}
		inputString = "<Input><K>10</K><Sums><decimal>1.01</decimal><decimal>2.02</decimal></Sums><Muls><int>1</int><int>4</int></Muls></Input>"
	} else {
		inputString = "{\"K\":10,\"Sums\":[1.01,2.02],\"Muls\":[1,4]}"
		serializer = JSONSerializer{}
	}
	input := serializer.Deserialize(inputString)
	var sortedInputs []float32
	sortedInputs = append(sortedInputs, input.Sums1[:len(input.Sums1)]...)
	var sumResults float32 = 0
	var mulResults float32 = 1.0
	for i := 0; i < len(input.Sums1); i++ {
		sumResults += input.Sums1[i] * float32(input.K1)
	}

	for i := 0; i < len(input.Muls1); i++ {
		mulResults *= float32(input.Muls1[i])
		sortedInputs = append(sortedInputs, float32(input.Muls1[i]))
	}
	sort.Sort(fl(sortedInputs))
	var output = Output{
		SumResult:    sumResults,
		MulResult:    mulResults,
		SortedInputs: sortedInputs,
	}

	fmt.Println(string(serializer.Serialize(output)))

}
