package main

import "testing"

func TestCaseA(t *testing.T) {

	Register("a", "X", []string{"true", "false"})
	Register("a", "Y", []string{"a", "b", "c"})

	if err := Generate("a"); err != nil {
		t.Error("Could not generate")
	}

	instance := make(map[string]string)

	instance["X"] = "true"
	instance["Y"] = "b"
	hash, err := Store("a", instance)

	if err != nil {
		t.Error("could not Store")
	}

	str, err := Get(hash)

	if err != nil {
		t.Error("could not Get", err)
	}

	t.Log(str)
}
