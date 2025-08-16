package jsonflatten

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFlattenTable_SimpleObject(t *testing.T) {
	input := map[string]interface{}{
		"a": 1,
		"b": "foo",
		"c": map[string]interface{}{
			"d": 2,
		},
	}
	want := []TableRow{
		{
			"a":   "1",
			"b":   "foo",
			"c/d": "2",
		},
	}
	got := FlattenTable(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SimpleObject: got %v, want %v", got, want)
	}
}

func TestFlattenTable_ArrayOfObjects(t *testing.T) {
	jsonStr := `[
		{"id": 1, "name": "a"},
		{"id": 2, "name": "b"}
	]`
	var input interface{}
	json.Unmarshal([]byte(jsonStr), &input)
	want := []TableRow{
		{"id": "1", "name": "a"},
		{"id": "2", "name": "b"},
	}
	got := FlattenTable(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ArrayOfObjects: got %v, want %v", got, want)
	}
}

func TestFlattenTable_NestedArray(t *testing.T) {
	jsonStr := `{
		"header": "X",
		"items": [
			{"a": 1, "b": {"c": 2}},
			{"a": 3, "b": {"c": 4}}
		]
	}`
	var input interface{}
	json.Unmarshal([]byte(jsonStr), &input)
	want := []TableRow{
		{"header": "X", "items/a": "1", "items/b/c": "2"},
		{"header": "X", "items/a": "3", "items/b/c": "4"},
	}
	got := FlattenTable(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NestedArray: got %v, want %v", got, want)
	}
}

func TestFlattenTable_MultiArray(t *testing.T) {
	jsonStr := `{
		"id": 1,
		"foo": [ {"a": 10}, {"a": 20} ],
		"bar": [ {"b": "x"}, {"b": "y"}, {"b": "z"} ]
	}`
	var input interface{}
	json.Unmarshal([]byte(jsonStr), &input)
	want := []TableRow{
		{"id": "1", "foo/a": "10", "bar/b": "x"},
		{"id": "1", "foo/a": "20", "bar/b": "y"},
		{"id": "1", "bar/b": "z"},
	}
	got := FlattenTable(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MultiArray: got %v, want %v", got, want)
	}
}

func TestFlattenTable_NonMapRoot(t *testing.T) {
	input := 42
	want := []TableRow{
		{"": "42"},
	}
	got := FlattenTable(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NonMapRoot: got %v, want %v", got, want)
	}
}

func TestFlattenTable_EmptyArray(t *testing.T) {
	jsonStr := `{"foo": []}`
	var input interface{}
	json.Unmarshal([]byte(jsonStr), &input)
	want := []TableRow{
		{"foo": ""},
	}
	got := FlattenTable(input)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("EmptyArray: got %v, want %v", got, want)
	}
}