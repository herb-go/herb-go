package tools

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/herb-go/herb-go/app"
)

func TestQuestion(t *testing.T) {
	var result string
	App := app.NewApplication(app.Config)
	inputr, inputw, _ := os.Pipe()
	outputw := bytes.NewBuffer([]byte{})
	App.Stdin = inputr
	App.Stdout = outputw
	question := NewQuestion()
	err := question.ExecIf(App, true, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result != "" {
		t.Fatal(result)
	}
	question = NewQuestion()
	question.
		SetDescription("test description").
		AddAnswer("0", "select0", "result0").
		AddAnswer("1", "select1", "result1").
		SetDefaultKey("0")
	err = question.ExecIf(App, false, nil)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		_, err = inputw.Write([]byte("1\n"))
		if err != nil {
			t.Fatal(err)
		}
	}()

	err = question.ExecIf(App, true, &result)
	outputs, err := ioutil.ReadAll(outputw)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(outputs), "test description") {
		t.Fatal(string(outputs))
	}
	if !strings.Contains(string(outputs), "Default choice is") {
		t.Fatal(string(outputs))
	}
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

	err = question.ExecIf(App, true, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result != "result0" {
		t.Fatal(err)
	}

	go func() {
		_, err = inputw.Write([]byte("a\n"))
		if err != nil {
			t.Fatal(err)
		}
	}()

	err = question.ExecIf(App, true, &result)
	if err == nil {
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
