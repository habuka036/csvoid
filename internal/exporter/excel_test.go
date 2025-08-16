package exporter

import (
	"bytes"
	"testing"

	"github.com/xuri/excelize/v2"
)

func readExcelFile(buf *bytes.Buffer) ([][]string, error) {
	f, err := excelize.OpenReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, err
	}
	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func TestExportExcel_Basic(t *testing.T) {
	rows := []TableRow{
		{"a": "1", "b": "x"},
		{"a": "2", "b": "y"},
	}
	var buf bytes.Buffer
	err := ExportExcel(rows, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := readExcelFile(&buf)
	if err != nil {
		t.Fatalf("failed to read excel: %v", err)
	}
	want := [][]string{
		{"a", "b"},
		{"1", "x"},
		{"2", "y"},
	}
	if !equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestExportExcel_ColumnUnionAndEmpty(t *testing.T) {
	rows := []TableRow{
		{"a": "1", "b": "x"},
		{"b": "y", "c": "zzz"},
	}
	var buf bytes.Buffer
	err := ExportExcel(rows, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := readExcelFile(&buf)
	if err != nil {
		t.Fatalf("failed to read excel: %v", err)
	}
	want := [][]string{
		{"a", "b", "c"},
		{"1", "x", ""},
		{"", "y", "zzz"},
	}
	if !equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestExportExcel_EmptyRows(t *testing.T) {
	rows := []TableRow{}
	var buf bytes.Buffer
	err := ExportExcel(rows, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := readExcelFile(&buf)
	if err != nil {
		t.Fatalf("failed to read excel: %v", err)
	}
	// excelize returns at least one sheet, and for empty, no rows
	if len(got) > 0 && len(got[0]) > 0 {
		t.Errorf("expected no output for empty rows, got: %v", got)
	}
}

// reuse the same equal function as csv_test.go
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
