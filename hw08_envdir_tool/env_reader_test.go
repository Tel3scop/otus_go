package main

import (
	"os"
	"testing"
)

func TestReadDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	err = os.WriteFile(dir+"/FOO", []byte("123\n"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(dir+"/BAR", []byte("value"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	env, err := ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	if env["FOO"].Value != "123" {
		t.Errorf("Expected FOO=123, got %s", env["FOO"].Value)
	}

	if env["BAR"].Value != "value" {
		t.Errorf("Expected BAR=value, got %s", env["BAR"].Value)
	}
}
