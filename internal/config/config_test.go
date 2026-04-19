package config

import (
	"testing"
)

func TestGetConfigFilePath(t *testing.T) {
	path, err := getConfigFilePath()
	if err != nil {
		t.Errorf("Expected no errors")
		t.Fail()
	}

	expected := "/home/todd"
	if path != expected {
		t.Errorf("expected %s != actual %s", expected, path)
		t.Fail()
	}
}

func TestRead(t *testing.T) {
	expected := Config{"postgres://example", "todd"}

	actual := Read()

	if actual != expected {
		t.Errorf("expected %v != actual %v", expected, actual)
		t.Fail()
	}
}
