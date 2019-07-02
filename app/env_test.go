package app

import (
	"os"
	"testing"
)

func TestOsEnv(t *testing.T) {
	envvalue := os.Getenv("testenvname")
	if envvalue != "" {
		t.Fatal(envvalue)
	}
	value := OsEnv.Getenv("testenvname")
	if value != "" {
		t.Fatal(value)
	}
	err := OsEnv.Setenv("testenvname", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		OsEnv.Setenv("testenvname", "")
	}()
	value = OsEnv.Getenv("testenvname")
	if value != "test" {
		t.Fatal(value)
	}
	envvalue = os.Getenv("testenvname")
	if envvalue != "test" {
		t.Fatal(envvalue)
	}
}

func TestPresetEnv(t *testing.T) {
	var TestEnv = PresetEnv(map[string]string{
		"test": "testvalue",
	})
	value := TestEnv.Getenv("test")
	if value != "testvalue" {
		t.Fatal(TestEnv)
	}
	err := TestEnv.Setenv("test", "testvalue2")
	if err != nil {
		t.Fatal(err)
	}
	value = TestEnv.Getenv("test")
	if value != "testvalue2" {
		t.Fatal(TestEnv)
	}
	app := NewApplication(Config)
	app.Env = TestEnv
	value = app.Getenv("test")
	if value != "testvalue2" {
		t.Fatal(TestEnv)
	}
	err = app.Setenv("test", "testvalue")
	if err != nil {
		t.Fatal(err)
	}
	value = app.Getenv("test")
	if value != "testvalue" {
		t.Fatal(TestEnv)
	}
}
