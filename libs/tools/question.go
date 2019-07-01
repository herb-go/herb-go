package tools

import (
	"bufio"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/herb-go/herb-go/app"
)

type Answer struct {
	Key   string
	Label string
	Value interface{}
}

func (a *Answer) Println(w io.Writer, defaultKey string) {
	if a.Key == defaultKey && defaultKey != "" {
		fmt.Fprintf(w, "*%s:%s\r\n", a.Key, a.Label)
		return
	}
	fmt.Fprintf(w, "%s:%s\r\n", a.Key, a.Label)
}
func NewAnswer() *Answer {
	return &Answer{}
}

type Question struct {
	Description string
	Answers     []*Answer
	DefaultKey  string
}

func (q *Question) SetDescription(d string) *Question {
	q.Description = d
	return q
}

func (q *Question) SetDefaultKey(k string) *Question {
	q.DefaultKey = k
	return q
}

func (q *Question) AddAnswer(key string, label string, value interface{}) *Question {
	a := NewAnswer()
	a.Key = key
	a.Label = label
	a.Value = value
	q.Answers = append(q.Answers, a)

	return q
}

func (q *Question) ExecIf(a *app.Application, conditon bool, result interface{}) error {
	if conditon == false {
		return nil
	}
	a.Println(q.Description)
	if len(q.Answers) == 0 {
		return nil
	}
	a.Printf("Please choose.\r\n")
	if q.DefaultKey != "" {
		a.Printf("Default choice is %s .\r\n", q.DefaultKey)
	}
	for _, v := range q.Answers {
		v.Println(a.Stdout, q.DefaultKey)
	}
	var s string
	scanner := bufio.NewScanner(a.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		return err
	}
	s = scanner.Text()
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	var key = s
	if s == "" && q.DefaultKey != "" {
		key = q.DefaultKey
	}

	for _, v := range q.Answers {
		if strings.TrimSpace(strings.ToLower(v.Key)) == key {
			reflect.Indirect(reflect.ValueOf(result)).Set(reflect.Indirect(reflect.ValueOf(v.Value)))
			return nil
		}
	}

	return fmt.Errorf("Choice \"%s\" is not  avaliable", s)
}

func NewQuestion() *Question {
	return &Question{
		Answers: []*Answer{},
	}
}
