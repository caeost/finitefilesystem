package main

import "testing"

func TestCaseA(t *testing.T) {

	Register("a", "X", []string{"true", "false"})
	Register("a", "Y", []string{"a", "b", "c"})

	if err := Generate("a"); err != nil {
		t.Error("Could not generate", err)
	}

	var hash string
	for i := 0; i < 1000; i++ {
		var err error
		hash, err = Store(generateInst())

		if err != nil {
			t.Error("could not Store", err)
		}
	}

	t.Log("hash: " + hash)

	str, err := Get(hash)

	if err != nil {
		t.Error("could not Get", err)
	}

	t.Log(str)
}

func generateInst() Instance {
	instance := make(Instance)

	instance["X"] = "true"
	instance["Y"] = "b"

	return instance
}
