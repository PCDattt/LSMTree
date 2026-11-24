package main

import (
	"os"
	"testing"
)

func TestEndToEnd(t *testing.T) {
	testCurrentFileName := "testCurrentFile.txt"
	testCompactedFileName := "testCompactedFile.txt"	
	s, err := NewStorage(testCurrentFileName, testCompactedFileName)
	if err != nil {
		t.Errorf("New Storage error: %v", err)
	}
	// Add data
	s.Add(5, "A")
	s.Add(10, "B")
	s.Add(1, "C")
	s.Add(7, "D")
	s.Add(3, "E")

	// Find from memory, before flush
	str, err := s.Find(5)
	if err != nil {
		t.Errorf("Find(5) error: %v", err)
	}
	n, err:= NodeUnmarshal(str)
	if err != nil {
		t.Errorf("NodeUnmarshal error: %v", err)
	}
	if n.Value != "A" {
		t.Errorf("s.Find(5) want A, get %s", n.Value)
	}

	str, err = s.Find(1)
	if err != nil {
		t.Errorf("Find(1) error: %v", err)
	}
	n, err= NodeUnmarshal(str)
	if err != nil {
		t.Errorf("NodeUnmarshal error: %v", err)
	}
	if n.Value != "C" {
		t.Errorf("s.Find(1) want C, get %s", n.Value)
	}

	str, err = s.Find(7)
	if err != nil {
		t.Errorf("Find(7) error: %v", err)
	}
	n, err= NodeUnmarshal(str)
	if err != nil {
		t.Errorf("NodeUnmarshal error: %v", err)
	}
	if n.Value != "D" {
		t.Errorf("s.Find(7) want D, get %s", n.Value)
	}

	err = s.Flush()
	if err != nil {
		t.Errorf("Flush 1 error: %v", err)
	}

	// Find from disk and memory, after flush 1
	str, err = s.Find(10)
	if err != nil {
		t.Errorf("Find(10) error: %v", err)
	}
	n, err = NodeUnmarshal(str)
	if err != nil {
		t.Errorf("NodeUnmarshal error: %v", err)
	}
	if n.Value != "B" {
		t.Errorf("s.Find(10) want B, get %s", n.Value)
	}
	
	str, err = s.Find(3)
	if err != nil {
		t.Errorf("Find(3) error: %v", err)
	}
	n, err= NodeUnmarshal(str)
	if err != nil {
		t.Errorf("NodeUnmarshal error: %v", err)
	}
	if n.Value != "E" {
		t.Errorf("s.Find(3) want E, get %s", n.Value)
	}

	// Add new data then find
	s.Add(5, "U")
	s.Add(3, "P")

	str, err = s.Find(5)
	if err != nil {
		t.Errorf("Find(5) after flush 1 error: %v", err)
	}
	n, err= NodeUnmarshal(str)
	if err != nil {
		t.Errorf("NodeUnmarshal error: %v", err)
	}
	if n.Value != "U" {
		t.Errorf("s.Find(5) after flush 1 want U, get %s", n.Value)
	}

	err = s.Compact()
	if err != nil {
		t.Errorf("Compact 1 error: %v", err)
	}

	err = s.Flush()
	if err != nil {
		t.Errorf("Flush 2 error: %v", err)
	}

	// Find from disk and memory, after compact and flush 2
	str, err = s.Find(3)
	if err != nil {
		t.Errorf("Find(3) after flush 2 error: %v", err)
	}
	n, err= NodeUnmarshal(str)
	if err != nil {
		t.Errorf("NodeUnmarshal error: %v", err)
	}
	if n.Value != "P" {
		t.Errorf("s.Find(3) after flush 2 want P, get %s", n.Value)
	}

	// Test construct index from file
	newIndex := NewIndex()
	newIndex.ConstructFromFile(s.compactedFile)
	newIndex.ConstructFromFile(s.currentFile)
	for k, m := range s.Index.MapIndex {
		newIndexMetadata, ok := newIndex.MapIndex[k]
		if !ok {
			t.Error("ConstructFromFile error")
		}
		if newIndexMetadata.FileName != m.FileName || newIndexMetadata.Offset != m.Offset {
			t.Errorf("Metadata want name: %s, offset: %d get name: %s, offset: %d", m.FileName, m.Offset, newIndexMetadata.FileName, newIndexMetadata.Offset)
		}
	}

	err = s.Compact()
	if err != nil {
		t.Errorf("Compact 2 error: %v", err)
	}

	str, err = s.Find(7)
	if err != nil {
		t.Errorf("Find(7) after compact 2 error: %v", err)
	}
	n, err= NodeUnmarshal(str)
	if err != nil {
		t.Errorf("NodeUnmarshal error: %v", err)
	}
	if n.Value != "D" {
		t.Errorf("s.Find(7) after compact 2 want D, get %s", n.Value)
	}

	// Clear test file
	s.Close()
	os.Remove(testCurrentFileName)
	os.Remove(testCompactedFileName)
}
