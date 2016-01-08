package main_test

import (
	"testing"

	"golang.org/x/tools/container/intsets"
)

func TestSetInsert(t *testing.T) {
	var X intsets.Sparse
	r := X.Insert(1)

	if !r  {
		t.Errorf("Insert: got %s, want %s", false, true)
	}
	r = X.Insert(1)
	if r  {
		t.Errorf("Insert: got %s, want %s", true, false)
	}

	r = X.Insert(8)
	if !r  {
		t.Errorf("Insert: got %s, want %s", false, true)
	}
	r = X.Insert(8)
	if r  {
		t.Errorf("Insert: got %s, want %s", true, false)
	}
}
