package main

import (
	"bytes"
	"testing"
)

func TestGetKubeNATTableLine(t *testing.T) {
	got, err := getKubeNATTableLines(exampleTables)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, wantTables) {
		t.Error("got wrong table")
		t.Error(string(wantTables))
		t.Error(string(got))
	}
}
