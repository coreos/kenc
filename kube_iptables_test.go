package main

import (
	"reflect"
	"testing"
)

func TestPurgeNonKubeLinesForIPtables(t *testing.T) {
	got, err := purgeNonKubeLinesForIPtables(exampleTables)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, wantTables) {
		t.Error("got wrong table after purging")
		t.Error(string(got))
	}
}
