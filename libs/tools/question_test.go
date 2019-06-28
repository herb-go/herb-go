package tools

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func TestQuestion(t *testing.T) {
	question := NewQuestion()
	question.
		SetDescription("test description").
		AddAnswer("0", "select0", "result0").
		AddAnswer("1", "select1", "result1").
		SetDefaultKey("0")
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
