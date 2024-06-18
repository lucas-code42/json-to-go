package main

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
)

const (
	MOCK_JSON_FILE string = "mock.json"
	MOCK_JSON_DATA string = `{
		"name": "json-to-struct",
		"decimal": 2.0,
		"boolean": false,
		"data": [
			{
			"foo": "bar"
			}
		],
		"number": 10
	}`
)

var MOCK_EXPECTED_STRUCT []string = []string{
	"packageauto_generated",
	"typeAutoGeneratedstruct{",
	"Namestringjson:name,omitempty",
	"Decimalfloat64json:decimal,omitempty",
	"Booleanbooljson:boolean,omitempty",
	"Data[]interface{}json:data,omitempty",
	"Numberfloat64json:number,omitempty",
}

func TestValidExtesion(t *testing.T) {
	result := validJsonExtension("test.js")
	if result {
		t.Errorf("expected false, but got %v", result)
		return
	}

	result = validJsonExtension("test.json")
	if !result {
		t.Errorf("expected true, but got %v", result)
		return
	}
}

func TestReadJsonFile(t *testing.T) {
	_, err := readJsonFile("test.js")
	if err == nil {
		t.Errorf("expected error, but got %e", err)
	}

	mockFile := func() {
		_, err := os.Create(MOCK_JSON_FILE)
		if err != nil {
			t.Errorf("could not create mock file")
			return
		}
	}
	mockFile()

	_, err = readJsonFile(MOCK_JSON_FILE)
	if err != nil {
		t.Errorf("expected nil, but got %e", err)
	}

	if err := os.Remove(MOCK_JSON_FILE); err != nil {
		t.Error("could not remove MOCK_JSON_FILE")
		return
	}
}

func TestValidStructContent(t *testing.T) {
	err := os.WriteFile(MOCK_JSON_FILE, []byte(MOCK_JSON_DATA), os.ModePerm)
	if err != nil {
		t.Error("could not create file 'test.json'")
		return
	}

	jsonBytes, err := readJsonFile(MOCK_JSON_FILE)
	if err != nil {
		t.Error("could not read MOCK_JSON_FILE")
		return
	}

	if err := os.Remove(MOCK_JSON_FILE); err != nil {
		t.Error("could not remove MOCK_JSON_FILE")
		return
	}

	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &jsonMap)
	if err != nil {
		t.Error("could not unmarshal, probably is not a valid json data for test")
		return
	}

	if err := writeGolangStruct(jsonMap); err != nil {
		t.Errorf("Expected nil but got %v", err)
		return
	}

	if _, err := os.Stat(FILE_NAME); errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected auto_generated.go but got %v", err)
		return
	}

	result, err := readJsonFile(FILE_NAME)
	if err != nil {
		t.Error("could not read FILE_NAME")
		return
	}

	if err := os.Remove(FILE_NAME); err != nil {
		t.Error("could not remove FILE_NAME")
		return
	}

	resultBuffer := string(result)
	resultBuffer = strings.ReplaceAll(resultBuffer, "`", "")
	resultBuffer = strings.ReplaceAll(resultBuffer, "\"", "")
	resultBuffer = strings.ReplaceAll(resultBuffer, "\n", "")
	resultBuffer = strings.ReplaceAll(resultBuffer, "\t", "")
	resultBuffer = strings.ReplaceAll(resultBuffer, " ", "")

	for i := range MOCK_EXPECTED_STRUCT {
		result := strings.Contains(resultBuffer, MOCK_EXPECTED_STRUCT[i])
		if !result {
			t.Error("Expected struct content doesn't match")
		}
	}
}

func TestMain(t *testing.T) {
	err := os.WriteFile(MOCK_JSON_FILE, []byte(MOCK_JSON_DATA), os.ModePerm)
	if err != nil {
		t.Error("could not create file 'test.json'")
		return
	}

	if err := os.Remove(MOCK_JSON_FILE); err != nil {
		t.Error("could not remove MOCK_JSON_FILE")
		return
	}

	os.Args = []string{"", MOCK_JSON_FILE}
	main()

	if err := os.Remove(FILE_NAME); err != nil {
		t.Error("could not remove FILE_NAME")
		return
	}
}
