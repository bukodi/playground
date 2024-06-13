package parsedocx

import (
	"strconv"
	"strings"
	"testing"
)

import (
	"fmt"
	"os"

	"github.com/fumiama/go-docx"
)

type Section struct {
	Heading  string
	Level    int
	Index    int
	Parent   *Section
	Children []*Section
	Items    []any
}

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
	dumpSection(root, "")
}

func parseToSections(parent *Section, startIdx int, items []any) (nextIdx int) {
	for idx := startIdx; idx < len(items); {
		it := items[idx]
		if p, ok := it.(*docx.Paragraph); ok {
			hl := extractHeadingLevel(p)
			text := p.String()
			if strings.HasPrefix(text, "Appendix") {
				fmt.Printf(" %s\n", text)
			}
			if hl == 0 {
				// Skip non-heading paragraph
				parent.Items = append(parent.Items, p)
				idx++
				continue
			}
			if hl > parent.Level {
				newChild := &Section{
					Heading: text,
					Level:   hl,
					Parent:  parent,
				}
				parent.Children = append(parent.Children, newChild)
				idx = parseToSections(newChild, idx+1, items)
				continue
			} else if hl <= parent.Level {
				fmt.Printf("New Section: %s (%d/%d)\n", parent.Heading, idx, len(items))
				return idx
			}
		} else if t, ok := it.(*docx.Table); ok {
			parent.Items = append(parent.Items, t)
			idx++
			continue
		} else {
			idx++
			continue
		}
	}
	return len(items)
}

func dumpSection(parent *Section, tab string) {
	fmt.Printf("%s%s\n", tab, parent.Heading)
	for _, child := range parent.Children {
		dumpSection(child, tab+"  ")
	}
}

func dumpItems(items []any) {
	for _, it := range items {
		if p, ok := it.(*docx.Paragraph); ok {
			fmt.Println(p)
		} else if t, ok := it.(*docx.Table); ok {
			fmt.Println(t)
		}
	}
}

func extractItems(items []any, headingLevel int, headingPrefix string) []any {
	var result []any
	for _, it := range items {
		if p, ok := it.(*docx.Paragraph); ok {
			hl := extractHeadingLevel(p)
			text := p.String()
			if result == nil && hl == headingLevel && strings.HasPrefix(text, headingPrefix) {
				// Start collection items
				result = make([]any, 0)
				result = append(result, it)
			} else if result != nil && hl > headingLevel {
				// Add to result
				result = append(result, it)
			} else if result != nil && hl <= headingLevel {
				// End collection items
				return result
			}
		} else if t, ok := it.(*docx.Table); ok {
			if result != nil {
				result = append(result, t)
			}
		}
	}
	panic("unexpected end of items")
}

func styleName(p *docx.Paragraph) string {
	if p.Properties != nil {
		if p.Properties.Style != nil {
			return p.Properties.Style.Val
		}
	}
	return ""
}

func extractHeadingLevel(p *docx.Paragraph) int {
	if p.Properties != nil {
		if p.Properties.Style != nil {
			name := p.Properties.Style.Val
			if strings.HasPrefix(name, "Heading") {
				idxtxt := name[len("Heading"):]
				idx, err := strconv.Atoi(idxtxt)
				if err != nil {
					panic(fmt.Errorf("cannot parse heading level from %s", name))
				}
				return idx
			}
		}
	}
	return 0
}
