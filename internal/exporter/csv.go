package exporter

import (
	"encoding/csv"
	"io"
	"sort"
)

// TableRow is a map representing a single row (flattened)
type TableRow map[string]string

// ExportCSV writes the rows as CSV to the given io.Writer.
// - rows: slice of TableRow (must not be nil; may be empty)
// - w: output stream (e.g. *os.File, *bytes.Buffer)
// - The first row is the header row (column names, sorted alphabetically)
func ExportCSV(rows []TableRow, w io.Writer) error {
	if len(rows) == 0 {
		// No data: write nothing
		return nil
	}
	// Collect all column names (union)
	colSet := map[string]struct{}{}
	for _, row := range rows {
		for col := range row {
			colSet[col] = struct{}{}
		}
	}
	cols := make([]string, 0, len(colSet))
	for col := range colSet {
		cols = append(cols, col)
	}
	sort.Strings(cols)

	cw := csv.NewWriter(w)
	defer cw.Flush()

	// Write header
	if err := cw.Write(cols); err != nil {
		return err
	}
	// Write data rows
	for _, row := range rows {
		record := make([]string, len(cols))
		for i, col := range cols {
			record[i] = row[col]
		}
		if err := cw.Write(record); err != nil {
			return err
		}
	}
	return cw.Error()
}