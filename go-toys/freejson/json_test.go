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

func TestUnmarshal(t *testing.T) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(testJson), &data); err != nil {
		t.Error(err)
	}
	// Output:
	out := new(strings.Builder)
	dumpMap(out, data, "")
	t.Log(out.String())
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
