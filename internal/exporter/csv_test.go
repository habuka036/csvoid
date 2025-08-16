package exporter

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
)

func TestExportCSV_Basic(t *testing.T) {
	rows := []TableRow{
		{"a": "1", "b": "x"},
		{"a": "2", "b": "y"},
	}
	var buf bytes.Buffer
	err := ExportCSV(rows, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := csv.NewReader(strings.NewReader(buf.String()))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("error reading csv: %v", err)
	}
	// Header should be sorted: a,b
	want := [][]string{
		{"a", "b"},
		{"1", "x"},
		{"2", "y"},
	}
	if !equal(records, want) {
		t.Errorf("got %v, want %v", records, want)
	}
}

func TestExportCSV_ColumnUnionAndEmpty(t *testing.T) {
	rows := []TableRow{
		{"a": "1", "b": "x"},
		{"b": "y", "c": "zzz"},
	}
	var buf bytes.Buffer
	err := ExportCSV(rows, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r := csv.NewReader(strings.NewReader(buf.String()))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("error reading csv: %v", err)
	}
	// Sorted: a,b,c
	want := [][]string{
		{"a", "b", "c"},
		{"1", "x", ""},
		{"", "y", "zzz"},
	}
	if !equal(records, want) {
		t.Errorf("got %v, want %v", records, want)
	}
}

func TestExportCSV_EmptyRows(t *testing.T) {
	rows := []TableRow{}
	var buf bytes.Buffer
	err := ExportCSV(rows, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty rows, got: %q", buf.String())
	}
}

/*
func equal(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
*/

func equal(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		maxLen := len(a[i])
		if len(b[i]) > maxLen {
			maxLen = len(b[i])
		}
		for j := 0; j < maxLen; j++ {
			var av, bv string
			if j < len(a[i]) {
				av = a[i][j]
			}
			if j < len(b[i]) {
				bv = b[i][j]
			}
			if av != bv {
				return false
			}
		}
	}
	return true
}