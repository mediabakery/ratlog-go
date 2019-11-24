package ratlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

type Data struct {
	Message string      `json:"message"`
	Fields  interface{} `json:"fields"`
	Tags    []string    `json:"tags"`
}

type Test struct {
	Log  string `json:"log"`
	Data Data   `json:"data"`
}

type Testsuite struct {
	Generic []Test `json:"generic"`
	Parsing []Test `json:"parsing"`
}

func getTestsuite() (Testsuite, error) {
	testsuiteFile, err := os.Open("ratlog.testsuite.json")
	if err != nil {
		return Testsuite{}, err
	}
	defer testsuiteFile.Close()
	testsuiteBytes, err := ioutil.ReadAll(testsuiteFile)
	if err != nil {
		return Testsuite{}, err
	}
	var testsuite Testsuite
	err = json.Unmarshal(testsuiteBytes, &testsuite)
	if err != nil {
		return Testsuite{}, err
	}

	return testsuite, nil
}

func getStringMap(in map[string]interface{}) map[string]string {
	out := make(map[string]string)
	for key, value := range in {
		if value != nil {
			strValue := fmt.Sprintf("%v", value)
			out[key] = strValue
		} else {
			out[key] = ""
		}
	}
	return out
}

func runTestsuite(t *testing.T, tests []Test, name string) {
	for index, test := range tests {
		var fields map[string]string
		if test.Data.Fields != nil {
			fields = getStringMap(test.Data.Fields.(map[string]interface{}))
		}
		var b bytes.Buffer
		ratlog := New(&b)
		ratlog.Log(Props{message: test.Data.Message, tags: test.Data.Tags, fields: fields})
		if test.Log != b.String() {
			t.Errorf("%s test # %d failed\nExpected: %#v\nReceived: %#v\n\n", name, index+1, test.Log, b.String())
		}
	}
}

func TestSuite(t *testing.T) {
	testsuite, err := getTestsuite()
	if err != nil {
		t.Error(err)
		return
	}
	runTestsuite(t, testsuite.Generic, "Generic")
	runTestsuite(t, testsuite.Parsing, "Parsing")
}

func TestTag(t *testing.T) {
	var b bytes.Buffer
	ratlog := New(&b)
	errorLag := ratlog.Tag("error", "debug")
	err := errorLag.Message("Hello World").Tag("bla", "baz").Fields("a", "2", "b", "3").Log()

	if err != nil {
		t.Error(err)
	}
	expected := "[error|debug|bla|baz] Hello World | a: 2 | b: 3\n"
	if b.String() != expected {
		t.Errorf("failed\nExpected: %#v\nReceived: %#v\n\n", expected, b.String())
	}
}

func TestBuilder(t *testing.T) {
	var b bytes.Buffer
	ratlog := New(&b)
	err := ratlog.Message("Hello World").Tag("bla", "baz").Fields("a", "2", "b", "3").Log()

	if err != nil {
		t.Error(err)
	}
	expected := "[bla|baz] Hello World | a: 2 | b: 3\n"
	if b.String() != expected {
		t.Errorf("failed\nExpected: %#v\nReceived: %#v\n\n", expected, b.String())
	}

}
