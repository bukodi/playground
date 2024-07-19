package parsedocx

import (
	"testing"
)

import (
	"fmt"
	"os"

	"github.com/fumiama/go-docx"
)

func TestParseDocx(t *testing.T) {
	readFile, err := os.Open("testdata/kmip-spec-v2.1.docx")
	if err != nil {
		panic(err)
	}
	fileinfo, err := readFile.Stat()
	if err != nil {
		panic(err)
	}
	size := fileinfo.Size()
	doc, err := docx.Parse(readFile, size)
	if err != nil {
		panic(err)
	}
	//dumpItems(doc.Document.Body.Items)

	//dumpItems(extractItems(doc.Document.Body.Items, 2, "Key Wrap Type Enumeration"))
	root := &Section{
		Heading: "Root",
	}
	lastIdx := parseToSections(root, 0, doc.Document.Body.Items)
	if lastIdx != len(doc.Document.Body.Items) {
		panic(fmt.Errorf("unexpected end of items"))
	}
	root.DumpTOC(os.Stdout, "")

	root.FindChild("Enumerations").DumpTOC(os.Stdout, "")
	root.FindChild("Enumerations").FindChild("Cica").DumpTOC(os.Stdout, "")
	//dumpSection(root, "")
}
