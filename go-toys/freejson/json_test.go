package freejson

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
)

const testJson = `{
	"string": "string",
	"int": 1,
	"subtree" : {
		"subkey": "value"
	},
	"array": [1,2,3,4],
	"float": 1.1
}`

func TestSingleString(t *testing.T) {
	var data any
	if err := json.Unmarshal([]byte(`"hello"`), &data); err != nil {
		t.Error(err)
	}
	// Output:
	out := new(strings.Builder)
	dumpNode(out, data, "")
	t.Log(out.String())
}

func TestUnmarshal(t *testing.T) {
	var data any
	if err := json.Unmarshal([]byte(testJson), &data); err != nil {
		t.Error(err)
	}
	// Output:
	out := new(strings.Builder)
	dumpNode(out, data, "")
	t.Log(out.String())
}

func dumpNode(out io.StringWriter, d any, indent string) {
	if d == nil {
		out.WriteString("null\n")
	} else if s, ok := d.(string); ok {
		out.WriteString(fmt.Sprintf("%s%q\n", indent, s))
	} else if m, ok := d.(map[string]any); ok {
		dumpMap(out, m, indent)
	}
}

func dumpMap(out io.StringWriter, m map[string]interface{}, indent string) {
	for k, v := range m {
		switch v := v.(type) {
		case string:
			out.WriteString(fmt.Sprintf("%s%q=%q\n", indent, k, v))
		case int:
			out.WriteString(fmt.Sprintf("%s%q=%d\n", indent, k, v))
		case float64:
			out.WriteString(fmt.Sprintf("%s%q=%f\n", indent, k, v))
		case []interface{}:
			out.WriteString(fmt.Sprintf("%s%q array:\n", indent, k))
			for _, v := range v {
				out.WriteString(fmt.Sprintf("%s%v\n", indent+"  ", v))
			}
		case map[string]interface{}:
			out.WriteString(fmt.Sprintf("%s%q map:\n", indent, k))
			dumpMap(out, v, indent+"  ")
		}
	}
}
