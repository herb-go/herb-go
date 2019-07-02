package tools

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestReplaceLine(t *testing.T) {
	data := "test\r\n\t    test\r\n      test2"
	data1 := "    test1\r\n    test1\r\n      test2"
	data2 := "    test1\r\n    test1\r\ntest2!"
	file, err := ioutil.TempFile("", "herbtest*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	file.Close()
	err = ioutil.WriteFile(file.Name(), []byte(data), 600)
	if err != nil {
		t.Fatal(err)
	}
	found, err := ReplaceLine("fileshoudnotexists.notexists", "test", "    test1")
	if err != nil {
		t.Fatal(err)
	}
	if found == true {
		t.Fatal(found)
	}
	found, err = ReplaceLine(file.Name(), "notfound", "    test1")
	if err != nil {
		t.Fatal(err)
	}
	if found == true {
		t.Fatal(found)
	}
	found, err = ReplaceLine(file.Name(), "test", "    test1")
	if err != nil {
		t.Fatal(err)
	}
	if found == false {
		t.Fatal(found)
	}
	bs, err := ioutil.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != data1 {
		t.Fatal(string(bs))
	}

	found, err = ReplaceLine(file.Name(), "test2", "test2!")
	if err != nil {
		t.Fatal(err)
	}
	if found == false {
		t.Fatal(found)
	}
	bs, err = ioutil.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}
	if string(bs) != data2 {
		t.Fatal(string(bs))
	}

}
