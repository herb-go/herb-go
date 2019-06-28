package tools

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/herb-go/herb-go/app"
)

func TestQuestion(t *testing.T) {
	App := app.NewApplication(app.Config)
	inputr, inputw, _ := os.Pipe()
	_, outputw, _ := os.Pipe()
	App.Stdin = inputr
	App.Stdout = outputw
	question := NewQuestion()
	question.
		SetDescription("test description").
		AddAnswer("0", "select0", "result0").
		AddAnswer("1", "select1", "result1").
		SetDefaultKey("0")
	err := question.Exec(App, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	var result string = ""
	go func() {
		_, err = inputw.Write([]byte("1\n"))
		if err != nil {
			t.Fatal(err)
		}
	}()

	err = question.Exec(App, true, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result != "result1" {
		t.Fatal(err)
	}

	go func() {
		_, err = inputw.Write([]byte(" \n"))
		if err != nil {
			t.Fatal(err)
		}
	}()

	err = question.Exec(App, true, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result != "result0" {
		t.Fatal(err)
	}
}

func TestAnswer(t *testing.T) {
	answer := NewAnswer()
	answer.Key = "0"
	answer.Label = "select0"
	answer.Value = "result0"
	buf := bytes.NewBuffer([]byte{})
	answer.Println(buf, "0")
	data, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "*0:select0\r\n" {
		t.Fatal(string(data))
	}
	answer.Println(buf, "1")
	data, err = ioutil.ReadAll(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "0:select0\r\n" {
		t.Fatal(string(data))
	}
}
