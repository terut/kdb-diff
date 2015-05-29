package main

import (
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	masterDB, _ := readXML("master_sample.xml")
	conflictDB, _ := readXML("conflict_sample.xml")
	masterEntries := make(map[string]Entry)
	conflictEntries := make(map[string]Entry)
	filterEntries(masterDB.Groups, &masterEntries)
	filterEntries(conflictDB.Groups, &conflictEntries)
	result := diff(masterEntries, conflictEntries)

	actual := result.MasterOnly
	expected := []string{"Sakura"}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
	actual = result.ConflictOnly
	expected = []string{"Hinata"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
	actual = result.Diff
	expected = []string{"Sasuke"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
