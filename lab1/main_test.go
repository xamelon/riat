package main

import "testing"

func TestXMLSerialization(t *testing.T) {
	xml := XMLSerializer{}
	xmlString := "<Input><K>10</K><Sums><decimal>1.01</decimal><decimal>2.02</decimal></Sums><Muls><int>1</int><int>4</int></Muls></Input>"
	input := xml.Deserialize(xmlString)
	inputStr := xml.Serialize(input)
	if string(inputStr) != xmlString {
		t.Error("FAIL")
	}
}

func TestJSONSerialization(t *testing.T) {
	json := JSONSerializer{}
	jsonString := "{\"K\":10,\"Sums\":[1.01,2.02],\"Muls\":[1,4]}"
	input := json.Deserialize(jsonString)
	inputStr := json.Serialize(input)
	if string(inputStr) != jsonString {
		t.Error("FAIL")
	}
}
