package exporter

import (
	"io"
	"sort"

	"github.com/xuri/excelize/v2"
)

// ExportExcel writes the rows as an Excel (xlsx) file to the given io.Writer.
// - rows: slice of TableRow (must not be nil; may be empty)
// - w: output stream (e.g. *os.File, *bytes.Buffer)
// - The first row is the header row (column names, sorted alphabetically)
func ExportExcel(rows []TableRow, w io.Writer) error {
	if len(rows) == 0 {
		// No data: write empty Excel with no header
		f := excelize.NewFile()
		err := f.Write(w)
		return err
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

	f := excelize.NewFile()
	sheet := f.GetSheetName(0)
	// Write header
	for i, col := range cols {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, col)
	}
	// Write rows
	for rowIdx, row := range rows {
		for colIdx, col := range cols {
			val := row[col]
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheet, cell, val)
		}
	}
	// Write to io.Writer
	return f.Write(w)
}
