package main

import (
	"os"
	"testing"
)

const MOCK_JSON_FILE = "mock.json"

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

	os.Remove(MOCK_JSON_FILE)
}
