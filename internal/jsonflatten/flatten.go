package jsonflatten

import (
	"fmt"
	"sort"
)

// TableRow represents a single row in the output table (for CSV/Excel export)
type TableRow map[string]string

// FlattenTable flattens a nested JSON-like structure into a slice of TableRow for tabular output.
// - root: input JSON data (map[string]interface{} or []interface{})
// - Arrays are expanded as row groups, nested keys are joined by "/".
// - If multiple arrays exist at the same level, they're expanded in parallel by row index.
func FlattenTable(root interface{}) []TableRow {
	switch v := root.(type) {
	case []interface{}:
		var rows []TableRow
		for _, elem := range v {
			rows = append(rows, flattenRow(elem, ""))
		}
		return rows
	case map[string]interface{}:
		return flattenObjectTable(v, "")
	default:
		return []TableRow{{"": fmt.Sprintf("%v", v)}}
	}
}

// flattenObjectTable: expands arrays in a map as row groups, flattens the rest.
func flattenObjectTable(obj map[string]interface{}, prefix string) []TableRow {
	// Find array fields at this level
	arrayKeys := []string{}
	for k, v := range obj {
		if _, ok := v.([]interface{}); ok {
			arrayKeys = append(arrayKeys, k)
		}
	}
	sort.Strings(arrayKeys) // for deterministic output

	if len(arrayKeys) == 0 {
		// No arrays: single row output
		return []TableRow{flattenRow(obj, prefix)}
	}

	// There are arrays: expand by max array length
	maxLen := 0
	arrays := map[string][]interface{}{}
	for _, k := range arrayKeys {
		arr, _ := obj[k].([]interface{})
		arrays[k] = arr
		if len(arr) > maxLen {
			maxLen = len(arr)
		}
	}
	// flatten non-array part as base row
	baseRow := flattenRow(obj, prefix)
	for _, k := range arrayKeys {
		delete(baseRow, joinKey(prefix, k))
	}
	if maxLen == 0 {
		// すべての配列が空の場合、空カラムの1行を出す
		row := make(TableRow)
		for k, v := range baseRow {
			row[k] = v
		}
		for _, k := range arrayKeys {
			row[joinKey(prefix, k)] = ""
		}
		return []TableRow{row}
	}
	var rows []TableRow
	for i := 0; i < maxLen; i++ {
		row := make(TableRow)
		// copy base fields
		for k, v := range baseRow {
			row[k] = v
		}
		// flatten array elements if present
		for _, k := range arrayKeys {
			arr := arrays[k]
			if i < len(arr) {
				flat := flattenRow(arr[i], joinKey(prefix, k))
				for kk, vv := range flat {
					row[kk] = vv
				}
			}
		}
		rows = append(rows, row)
	}
	return rows
}

// flattenRow: flatten a value (non-array) as a single row
func flattenRow(data interface{}, prefix string) TableRow {
	row := make(TableRow)
	switch v := data.(type) {
	case map[string]interface{}:
		for k, vv := range v {
			key := joinKey(prefix, k)
			switch vv.(type) {
			case []interface{}:
				// handled by flattenObjectTable
				continue
			default:
				for kk, vv2 := range flattenRow(vv, key) {
					row[kk] = vv2
				}
			}
		}
	case []interface{}:
		// arrays are handled in flattenObjectTable; skip here
	default:
		row[prefix] = fmt.Sprintf("%v", v)
	}
	return row
}

func joinKey(a, b string) string {
	if a == "" {
		return b
	}
	return a + "/" + b
}