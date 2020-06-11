package helper

import (
	"io/ioutil"
	"testing"
)

func TestGenerate(t *testing.T) {
	avatar, err := ioutil.ReadFile("../test/avatar.png")
	if err != nil {
		t.Fatal(err)
	}
	logo, err := ioutil.ReadFile("../test/logo.png")
	if err != nil {
		t.Fatal(err)
	}
	code, err := ioutil.ReadFile("../test/code.png")
	if err != nil {
		t.Fatal(err)
	}
	if result, err := Generate(code, avatar, logo, map[string]string{}); err != nil {
		t.Fatal(err)
	} else {
		ioutil.WriteFile("../test/result.png", result, 0777)
	}
}
